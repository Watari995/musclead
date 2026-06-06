package userhandler

import (
	"net/http"

	"time"

	"github.com/Watari995/musclead/internal/myerror"
	shareddto "github.com/Watari995/musclead/internal/shared/dto"
	"github.com/Watari995/musclead/internal/shared/httpx"
	userdto "github.com/Watari995/musclead/internal/user/dto"
	userusecase "github.com/Watari995/musclead/internal/user/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
)

type UserHandler struct {
	me         *userusecase.Me
	register   *userusecase.RegisterUser
	find       *userusecase.FindUser
	updateUser *userusecase.UpdateUser
	delete     *userusecase.DeleteUser
}

func NewPublic(register *userusecase.RegisterUser) http.Handler {
	h := &UserHandler{
		register: register,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.Register)
	return mux
}

func NewAuthenticated(me *userusecase.Me, find *userusecase.FindUser, updateUser *userusecase.UpdateUser, delete *userusecase.DeleteUser) http.Handler {
	// ServeHTTP interfaceを満たしている必要がある
	h := &UserHandler{
		me:         me,
		find:       find,
		updateUser: updateUser,
		delete:     delete,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/me", h.Me)
	mux.HandleFunc("GET /users/{id}", h.Find)
	mux.HandleFunc("PATCH /users/me", h.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", h.Delete)
	return mux
}

// Me godoc
//
// @Summary ユーザー情報取得
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userdto.UserDTO
// @Failure 401 {object} httpx.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	params := userusecase.MeInput{
		UserID: userID,
	}
	output, err := h.me.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := userdto.NewUserDTO(output.User.ID(), output.User.Name(), output.User.Email(), output.User.Birthday(), output.User.CreatedAt(), output.User.UpdatedAt())
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// Register godoc
//
// @Summary ユーザー登録
// @Description 新規ユーザーを作成する(認証不要)
// @Tags users
// @Accept json
// @Produce json
// @Param request body userdto.RegisterRequest true "ユーザー登録情報"
// @Success 201 {object} userdto.RegisterResponse "ユーザー登録成功"
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Router /users [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req userdto.RegisterRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}

	name, err := valueobject.NewString50(req.Name)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
		return
	}
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid email"))
		return
	}
	var birthday *time.Time
	if req.Birthday != nil {
		t, err := time.Parse("2006-01-02", *req.Birthday)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid birthday"))
			return
		}
		birthday = &t
	}

	params := userusecase.RegisterUserInput{
		Name:     *name,
		Email:    *email,
		Birthday: birthday,
		Password: req.Password,
	}

	output, err := h.register.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	resp := userdto.RegisterResponse{
		UserID: output.UserID.Value(),
	}
	httpx.WriteJSON(w, http.StatusCreated, resp)
}

// Find godoc
//
// @Summary ユーザー取得
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "対象 UserID"
// @Success 200 {object} userdto.UserDTO
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) Find(w http.ResponseWriter, r *http.Request) {
	// path parameterからuserIDを取得
	userID, err := valueobject.NewPrimaryIDFromString[valueobject.UserID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid userID"))
		return
	}
	params := userusecase.FindUserInput{
		UserID: *userID,
	}
	output, err := h.find.Execute(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	resp := userdto.NewUserDTO(output.UserID, output.Name, output.Email, output.Birthday, output.CreatedAt, output.UpdatedAt)
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req userdto.UpdateUserRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid request body"))
		return
	}
	if req.Name.Set && req.Name.Null {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("name is required"))
		return
	}
	var namePatch shareddto.Patch[valueobject.String50]
	// decodeされて値がセットされている場合に限り更新する
	if req.Name.Set {
		namePatch.Set = true
		nameVO, err := valueobject.NewString50(req.Name.Value)
		if err != nil {
			httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid name"))
			return
		}
		namePatch.Value = *nameVO
	}
	var birthdayPatch shareddto.Patch[time.Time]
	// decodeされて値がセットされている場合に限り更新する
	if req.Birthday.Set {
		birthdayPatch.Set = true
		if req.Birthday.Null {
			birthdayPatch.Null = true
		} else {
			birthdayTime, err := time.Parse("2006-01-02", req.Birthday.Value)
			if err != nil {
				httpx.WriteError(w, err)
				return
			}
			birthdayPatch.Value = birthdayTime
		}
	}
	output, err := h.updateUser.Execute(r.Context(), userusecase.UpdateUserInput{
		UserID:   userID,
		Name:     namePatch,
		Birthday: birthdayPatch,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, userdto.UpdateUserResponse{UserID: output.UserID.Value()})
}

// Delete godoc
//
// @Summary ユーザー削除
// @Tags users
// @Security BearerAuth
// @Param id path string true "対象 UserID"
// @Success 204
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := valueobject.NewPrimaryIDFromString[valueobject.UserID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid userID"))
		return
	}
	if err := h.delete.Execute(r.Context(), userusecase.DeleteUserInput{UserID: *userID}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteNoContent(w)
}
