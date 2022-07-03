package controllers

import (
	"net/http"
	"tes-mitramas/models"
	u "tes-mitramas/utils"
	"time"
)

var RiwayatAbsensi = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		u.RespondError(w, u.Message(false, "Request tidak valid"), http.StatusMethodNotAllowed)
		return
	}

	accid := r.Context().Value("values").(Valcontex).Id
	username := r.Context().Value("values").(Valcontex).Username

	absenlists, err := models.RiwayatAbsensi(username, accid)
	if err != nil {
		u.RespondError(w, u.Message(false, err.Error()), http.StatusBadRequest)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = absenlists
	u.Respond(w, resp)
}

var RiwayatAbsensiByTanggal = func(w http.ResponseWriter, r *http.Request) {
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

	startdate := r.URL.Query().Get("start_date")
	enddate := r.URL.Query().Get("end_date")

	starttime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", startdate+" 00:00:00 +0700 WIB")
	if err != nil {
		u.Respond(w, u.Message(false, "Request filter date tidak valid"))
		return
	}

	endtime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", enddate+" 00:00:00 +0700 WIB")
	if err != nil {
		u.Respond(w, u.Message(false, "Request filter date tidak valid"))
		return
	}

	absenlists, err := models.RiwayatAbsensiPerTanggal(username, accid, starttime, endtime)
	if err != nil {
		u.RespondError(w, u.Message(false, err.Error()), http.StatusBadRequest)
		return
	}
	resp := u.Message(true, "success")
	resp["data"] = absenlists
	u.Respond(w, resp)
}
