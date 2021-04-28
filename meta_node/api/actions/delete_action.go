package actions

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/GDVFox/dflow/meta_node/api/common"
	"github.com/GDVFox/dflow/meta_node/external"
	"github.com/GDVFox/dflow/util/httplib"
	"github.com/GDVFox/dflow/util/storage"
)

// DeleteAction удаляем действие, если оно существует.
func DeleteAction(r *http.Request) (*httplib.Response, error) {
	vars := mux.Vars(r)
	actionName := vars["action_name"]
	if actionName == "" {
		return httplib.NewBadRequestResponse(httplib.NewErrorBody(common.BadNameErrorCode, "action_name must be not empty")), nil
	}

	if err := external.ETCD.DeleteAction(r.Context(), actionName); err != nil {
		if errors.Cause(err) == storage.ErrNotFound {
			return httplib.NewNotFoundResponse(httplib.NewErrorBody(common.NameNotFoundErrorCode, err.Error())), nil
		}
		return httplib.NewInternalErrorResponse(httplib.NewErrorBody(common.ETCDErrorCode, err.Error())), nil
	}

	return httplib.NewOKResponse(nil, false), nil
}
