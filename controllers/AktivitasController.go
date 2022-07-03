package controllers

import (
	"encoding/json"
	"net/http"
	"tes-mitramas/models"
	u "tes-mitramas/utils"
)

var CheckInCont = func(w http.ResponseWriter, r *http.Request) {
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

	if acc.CheckIn {
		u.Respond(w, u.Message(false, "anda telah checkin hari ini"))
		return
	}

	err := acc.DoCheckin()
	if err != nil {
		u.RespondError(w, u.Message(false, err.Error()), http.StatusBadRequest)
		return
	}

	resp := u.Message(true, "berhasil checkin")
	u.Respond(w, resp)
}

var CheckOutCont = func(w http.ResponseWriter, r *http.Request) {
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

	if !acc.CheckIn {
		u.Respond(w, u.Message(false, "anda telah checkout"))
		return
	}

	err := acc.DoCheckout()
	if err != nil {
		u.RespondError(w, u.Message(false, err.Error()), http.StatusBadRequest)
		return
	}

	resp := u.Message(true, "berhasil checkout")
	u.Respond(w, resp)
}

var CreateAktivitas = func(w http.ResponseWriter, r *http.Request) {
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

	aktivitas := &models.Aktivitas{}

	err := json.NewDecoder(r.Body).Decode(aktivitas)

	if err != nil {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusBadRequest)
		return
	}

	if aktivitas.AccountId != accid {
		u.RespondError(w, u.Message(false, "akun id tidak sama"), http.StatusBadRequest)
		return
	}

	resp := models.CreateAct(aktivitas)

	if !resp["status"].(bool) {
		u.RespondError(w, resp, http.StatusBadRequest)
		return
	}

	u.Respond(w, resp)
}

var EditAktivitas = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusMethodNotAllowed)
		return
	}

	jsonedit := &models.EditAktivitasJson{}

	err := json.NewDecoder(r.Body).Decode(&jsonedit)
	if err != nil {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusBadRequest)
		return
	}

	accid := r.Context().Value("values").(Valcontex).Id
	username := r.Context().Value("values").(Valcontex).Username

	acc := models.GetAccount(username, accid)

	if acc == nil {
		u.RespondError(w, u.Message(false, "akun tidak ditemukan"), http.StatusBadRequest)
		return
	}

	resp := models.EditAct(jsonedit, accid)
	if !resp["status"].(bool) {
		u.RespondError(w, resp, http.StatusBadRequest)
		return
	}

	u.Respond(w, resp)
}

var DeleteAktivitas = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusMethodNotAllowed)
		return
	}

	jsonedit := &models.EditAktivitasJson{}

	err := json.NewDecoder(r.Body).Decode(&jsonedit)
	if err != nil {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusBadRequest)
		return
	}

	accid := r.Context().Value("values").(Valcontex).Id
	username := r.Context().Value("values").(Valcontex).Username

	acc := models.GetAccount(username, accid)

	if acc == nil {
		u.RespondError(w, u.Message(false, "akun tidak ditemukan"), http.StatusBadRequest)
		return
	}

	resp := models.DeleteAct(jsonedit, accid)
	if !resp["status"].(bool) {
		u.RespondError(w, resp, http.StatusBadRequest)
		return
	}

	u.Respond(w, resp)
}
