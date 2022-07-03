package controllers

import (
	"encoding/json"
	"net/http"
	"tes-mitramas/models"
	u "tes-mitramas/utils"
)

type Valcontex struct {
	Username string `json:"username"`
	Id       int64  `json:"id"`
}

var LoginCont = func(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusMethodNotAllowed)
		return
	}

	login := &models.JsonLogin{}

	err := json.NewDecoder(r.Body).Decode(login)
	if err != nil {
		u.Respond(w, u.Message(false, "Request tidak valid"))
		return
	}

	resp := models.Login(login)
	u.Respond(w, resp)
}

var LogoutCont = func(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusMethodNotAllowed)
		return
	}

	accid := r.Context().Value("values").(Valcontex).Id
	username := r.Context().Value("values").(Valcontex).Username

	acc := models.GetAccount(username, accid)

	if acc == nil {
		u.RespondError(w, u.Message(false, "akun tidak ditemukan"), http.StatusBadRequest)
		return
	}

	resp := models.Logout(acc)

	if !resp["status"].(bool) {
		u.RespondError(w, resp, http.StatusNotAcceptable)
		return
	}
	u.Respond(w, resp)
}

var RegisterCont = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusMethodNotAllowed)
		return
	}

	jsonregist := &models.RegisterAccountJson{}

	err := json.NewDecoder(r.Body).Decode(jsonregist)
	if err != nil {
		u.Respond(w, u.Message(false, "Request tidak valid"))
		return
	}

	resp := jsonregist.CreateAccount()

	if !resp["status"].(bool) {
		u.RespondError(w, resp, http.StatusNotAcceptable)
		return
	}

	u.Respond(w, resp)
}

var GetTokenCont = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusMethodNotAllowed)
		return
	}

	accid := r.Context().Value("values").(Valcontex).Id
	username := r.Context().Value("values").(Valcontex).Username

	acc := models.GetAccount(username, accid)

	if acc == nil {
		u.RespondError(w, u.Message(false, "akun tidak ditemukan"), http.StatusBadRequest)
		return
	}

	resp := models.GetToken(username, accid)

	if !resp["status"].(bool) {
		u.RespondError(w, resp, http.StatusNotAcceptable)
		return
	}

	u.Respond(w, resp)
}
