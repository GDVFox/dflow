package watcher

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/GDVFox/dflow/util"
)

const (
	// PingCommand команда для проверки состояния runtime.
	PingCommand uint8 = 0x1
)

const (
	// OKResponse ответ, предполагающий успешное выполнение действия.
	OKResponse uint8 = 0x0
	// FailResponse ответ, предполагающий ошибочное выполнение действия.
	FailResponse uint8 = 0x1
)

// Возможные ошибки
var (
	ErrPingFailed = errors.New("ping returned not OK response")
)

// ActionOptions опции для запуска действия
type ActionOptions struct {
	Args []string          `json:"args"`
	Env  map[string]string `json:"env"`
}

// RuntimeOptions набор параметров при запуске действия.
type RuntimeOptions struct {
	Port          int
	Replicas      int
	In            []string
	Out           []string
	ActionOptions *ActionOptions

	RuntimePath      string
	RuntimeLogsDir   string
	ActionStartRetry *util.RetryConfig
}

// Runtime структура, представляющая собой запущенное действие
type Runtime struct {
	name string
	bin  []byte
	opt  *RuntimeOptions

	binPath         string
	cmd             *exec.Cmd
	stderr          io.ReadCloser
	serviceSockPath string
	serviceConn     net.Conn

	logger *util.Logger
}

// NewRuntime создает новое действие.
func NewRuntime(schemeName, actionName string, bin []byte, l *util.Logger, opt *RuntimeOptions) *Runtime {
	return &Runtime{
		name: buildRuntimeName(schemeName, actionName),
		bin:  bin,
		opt:  opt,

		logger: l,
	}
}

// Name возвращает имя запущенного действия
func (r *Runtime) Name() string {
	return r.name
}

// Start запускает действие и выходит, в случае успешного запуска.
func (r *Runtime) Start(ctx context.Context) error {
	var err error
	r.binPath, err = r.createTmpBinary(r.name, r.bin)
	if err != nil {
		return fmt.Errorf("can not create binary: %w", err)
	}
	r.logger.Debugf("created tmp binary with path %s", r.binPath)

	actionOptions, err := json.Marshal(r.opt.ActionOptions)
	if err != nil {
		return fmt.Errorf("invalid action options: %w", err)
	}

	r.serviceSockPath = filepath.Join("/", "tmp", r.name+strconv.Itoa(r.opt.Port)+".sock")
	logFileAddr := filepath.Join("/", r.opt.RuntimeLogsDir, r.name+strconv.Itoa(r.opt.Port)+".log")
	r.cmd = exec.Command(
		r.opt.RuntimePath,
		"--action="+r.binPath,
		"--replicas="+strconv.Itoa(r.opt.Replicas),
		"--port="+strconv.Itoa(r.opt.Port),
		"--service-sock="+r.serviceSockPath,
		"--log-file="+logFileAddr,
		"--log-level=debug",
		"--in="+strings.Join(r.opt.In, ","),
		"--out="+strings.Join(r.opt.Out, ","),
		"--action-opt="+string(actionOptions),
	)

	r.stderr, err = r.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("can not get stderr connect: %w", err)
	}

	if err := r.cmd.Start(); err != nil {
		return fmt.Errorf("can not start action: %w", err)
	}
	r.logger.Debugf("runtime started with command: %s", r.cmd.String())

	if err := r.connect(ctx); err != nil {
		return fmt.Errorf("can not connect to action: %w", err)
	}

	return nil
}

func (r *Runtime) createTmpBinary(name string, bin []byte) (string, error) {
	actionFile, err := ioutil.TempFile("", name)
	if err != nil {
		return "", fmt.Errorf("can not create tmp file for bin: %w", err)
	}
	defer actionFile.Close()

	if _, err := actionFile.Write(bin); err != nil {
		return "", fmt.Errorf("can not writer bin: %w", err)
	}

	if err := os.Chmod(actionFile.Name(), 0700); err != nil {
		return "", fmt.Errorf("can not change file mod: %w", err)
	}

	return actionFile.Name(), nil
}

func (r *Runtime) connect(ctx context.Context) error {
	dialErr := util.Retry(ctx, r.opt.ActionStartRetry, func() error {
		var err error
		r.serviceConn, err = net.Dial("unix", r.serviceSockPath)
		if err != nil {
			return err
		}

		if err := r.Ping(); err != nil {
			return err
		}

		r.logger.Debug("start confirm received")
		return nil
	})
	if dialErr != nil {
		if r.cmd.Process == nil {
			return dialErr
		}

		if err := r.Stop(); err != nil {
			return err
		}
		return dialErr
	}
	return nil
}

// Ping проверяет работоспособность действия с помощью отправки ping.
func (r *Runtime) Ping() error {
	if err := binary.Write(r.serviceConn, binary.BigEndian, PingCommand); err != nil {
		return err
	}

	var resp uint8
	if err := binary.Read(r.serviceConn, binary.BigEndian, &resp); err != nil {
		return err
	}

	if resp != OKResponse {
		return ErrPingFailed
	}
	return nil
}

// Stop завершает работу действия, возвращает ошибку из stderr.
func (r *Runtime) Stop() error {
	defer os.Remove(r.binPath)
	defer os.Remove(r.serviceSockPath)
	defer func() {
		if r.stderr != nil {
			r.stderr.Close()
		}
	}()
	defer func() {
		if r.serviceConn != nil {
			r.serviceConn.Close()
		}
	}()

	if err := r.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("can not send SIGTERM to runtime: %w", err)
	}

	state, err := r.cmd.Process.Wait()
	if err != nil {
		return fmt.Errorf("can not wait and of runtime process: %w", err)
	}

	// успешное завершение процесса
	if state.Success() {
		return nil
	}

	b := &bytes.Buffer{}
	if _, err := io.Copy(b, r.stderr); err != nil {
		return fmt.Errorf("can not copy stderr: %w", err)
	}

	return errors.New(b.String())
}

func buildRuntimeName(schemeName, actionName string) string {
	return schemeName + "_" + actionName
}