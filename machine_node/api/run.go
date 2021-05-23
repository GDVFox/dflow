package api

import (
	"encoding/json"
	"net/http"

	"github.com/GDVFox/dflow/machine_node/config"
	"github.com/GDVFox/dflow/machine_node/external"
	"github.com/GDVFox/dflow/machine_node/watcher"
	"github.com/GDVFox/dflow/util"
	"github.com/GDVFox/dflow/util/httplib"
	"github.com/GDVFox/dflow/util/message"
	"github.com/GDVFox/dflow/util/storage"
	"github.com/pkg/errors"
)

// RunAction запускает action.
func RunAction(r *http.Request) (*httplib.Response, error) {
	logger := r.Context().Value(httplib.RequestLogger).(*util.Logger)

	req := &message.RunActionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httplib.NewBadRequestResponse(httplib.NewErrorBody(BadUnmarshalRequestErrorCode, err.Error())), nil
	}

	actionBytes, err := external.ETCD.LoadAction(r.Context(), req.Action)
	if err != nil {
		if errors.Cause(err) == storage.ErrNotFound {
			return httplib.NewNotFoundResponse(httplib.NewErrorBody(NoActionErrorCode, err.Error())), nil
		}
		return httplib.NewInternalErrorResponse(httplib.NewErrorBody(ETCDErrorCode, err.Error())), nil
	}
	logger.Debugf("binary action '%s' received", req.Action)

	opt := &watcher.RuntimeOptions{
		Port:             req.Port,
		In:               req.In,
		Out:              req.Out,
		RuntimePath:      config.Conf.RuntimePath,
		RuntimeLogsDir:   config.Conf.RuntimeLogsDir,
		RuntimeLogsLevel: config.Conf.RuntimeLogsLevel,
		ActionStartRetry: config.Conf.ActionStartRetry,
		ActionOptions: &watcher.ActionOptions{
			Args: req.Args,
			Env:  req.Env,
		},
	}
	runtime := watcher.NewRuntime(req.SchemeName, req.ActionName, actionBytes, logger, opt)

	if err := watcher.RuntimeWatcher.StartRuntime(r.Context(), runtime); err != nil {
		return httplib.NewInternalErrorResponse(httplib.NewErrorBody(InternalError, err.Error())), nil
	}
	logger.Debugf("runtime '%s' started", runtime.Name())

	logger.Infof("started action '%s' from scheme '%s'", req.ActionName, req.SchemeName)
	return httplib.NewOKResponse(nil, httplib.ContentTypeRaw), nil
}
