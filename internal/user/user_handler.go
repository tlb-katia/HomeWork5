package user

import (
	"encoding/json"
	"io"
	"net/http"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u UserReq

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRes, err := h.Service.CreateUser(r.Context(), &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(userRes)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var u UserReq

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginUser, err := h.Service.Login(r.Context(), &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// we need to delete jwt cookies
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(w, &c)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
