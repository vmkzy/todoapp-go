package users_transport_http

import (
	"net/http"

	core_logger "github.com/vmkzy/todoapp-go/internal/core/logger"
	core_http_response "github.com/vmkzy/todoapp-go/internal/core/transport/http/response"
	core_http_utils "github.com/vmkzy/todoapp-go/internal/core/transport/http/utils"
)

func (h *UserHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userId path value",
		)
	}
	if err := h.userService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}
	responseHandler.NoContentResponse()
}
