package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vmkzy/todoapp-go/internal/core/domain"
	core_logger "github.com/vmkzy/todoapp-go/internal/core/logger"
	core_http_request "github.com/vmkzy/todoapp-go/internal/core/transport/http/request"
	core_http_response "github.com/vmkzy/todoapp-go/internal/core/transport/http/response"
	core_http_types "github.com/vmkzy/todoapp-go/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Petya"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+375777777777"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}
		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}
	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be beetween 10 and 15 symbols")
			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber mist starts with '+' symbol")
			}
		}

	}
	return nil
}

type PatchUserResponse UserDTORepsonse

// PatchUser 	godoc
// @Summary 	Изменение пользователя
// @Description Изменение существующего пользователя по ID
// @Description ### Логика обновление полей (Three-state logic):
// @Description 1. **Поле не передано**: `phone_bumber` игнорируется, значение в бд не меняется
// @Description 2. **Явно передано значение**: `phone_number:"+375..."` устанавливается новое значение в бд
// @Description 3. **Передано null**: `phone_number:null` очищается поле в бд
// @Description Ограничения: `full_name` не может быть выставлен как null
// @Tags 		users
// @Accept		json
// @Produce 	json
// @Param 		id path int true "ID пользователя"
// @Param		request body PatchUserRequest true "PatchUser тело запроса"
// @Success 	200 {object} PatchUserResponse "Пользователь успешно изменен"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 	409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users/{id} [patch]
func (h *UserHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID patch value",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validete HTTP request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)
	userDomain, err := h.userService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}
	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}
func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
