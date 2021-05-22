package main

import (
	"encoding/binary"
	"flag"
	"strconv"

	"github.com/GDVFox/dflow/lib/go-actionlib"
)

var (
	mod int
)

func init() {
	flag.IntVar(&mod, "mod", 1, "passes messages that are multiples of mod")
}

func main() {
	flag.Parse()

	for {
		data, err := actionlib.ReadMessage()
		if err != nil {
			actionlib.WriteError(err)
		}

		number := binary.BigEndian.Uint32(data)
		if int(number)%mod != 0 {
			actionlib.AckMessage()
			continue
		}

		outputData := []byte(strconv.Itoa(int(number)))
		if err := actionlib.WriteMessage(outputData); err != nil {
			actionlib.WriteError(err)
		}
	}

}
