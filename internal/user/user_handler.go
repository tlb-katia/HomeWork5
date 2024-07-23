package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	Service
	*slog.Logger
}

func NewHandler(log *slog.Logger, s Service) *Handler {
	return &Handler{s, log}
}

func (h *Handler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
	logResponseStatusError(h.Logger, message, statusCode)
}

func (h *Handler) sendSuccessResponse(w http.ResponseWriter, ur *UserRes, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ur)
	logResponseSuccess(h.Logger, message, statusCode)
}

func logResponseStatusError(log *slog.Logger, message string, statusCode int) {
	log.Error("Request error", "status", statusCode, "error", message)
}

func logResponseSuccess(log *slog.Logger, message string, statusCode int) {
	log.Info("Request success", slog.Int("status", statusCode), slog.String("message", message))
}

// CreateUser godoc
// @Summary      create a user
// @Description  Create a new user with username, email, and password
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      UserReq  true  "User request body"
// @Success      200   {object}  User
// @Failure      400   {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /signup [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u UserReq

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.sendErrorResponse(w, "Invalid request to create a user", http.StatusBadRequest)
		return
	}

	userRes, err := h.Service.CreateUser(r.Context(), &u)
	if err != nil {
		h.sendErrorResponse(w, "Couldn't create a user", http.StatusInternalServerError)

		return
	}

	h.sendSuccessResponse(w, userRes, "User created successfully", http.StatusCreated)
}

// LoginUser godoc
// @Summary      log in a user
// @Description  Log in a user with username, email, and password
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      UserReq  true  "User request body"
// @Success      200   {object}  User
// @Failure      400   {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /login [post]
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var u UserReq

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.sendErrorResponse(w, "Invalid request to create a user", http.StatusBadRequest)
		return
	}

	loginUser, err := h.Service.Login(r.Context(), &u)
	if err != nil {
		h.sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	tokenString := loginUser.Token

	cookie := &http.Cookie{
		Name:   "token",
		Value:  tokenString,
		MaxAge: 24 * 60 * 60,
		Path:   "/",
	}
	http.SetCookie(w, cookie)

	h.sendSuccessResponse(w, &UserRes{Message: "user was successfully logged in"}, "User logged in successfully", http.StatusOK)
}

// LogoutUser godoc
// @Summary      Log out user
// @Description  Logs out the user by clearing all cookies and redirecting to the login page.
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  UserRes  "Successfully logged out"
// @Router       /logout [get]
func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(w, &c)
	}
	h.sendSuccessResponse(w, &UserRes{Message: "user was successfully logged out"}, "Logout successful", http.StatusOK)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
