package models

import (
	"strings"
	u "tes-mitramas/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

type JsonLogin struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type ApiKeyAccount struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(data *JsonLogin) map[string]interface{} {
	if data.Username == "" || data.Password == "" {
		return u.Message(false, "username atau password tidak boleh kosong")
	}

	if len(strings.TrimSpace(data.Username)) == 0 || len(strings.TrimSpace(data.Password)) == 0 {
		return u.Message(false, "username atau password tidak boleh kosong")
	}

	conn := GetDB()
	defer conn.Close()

	acc := &Account{}
	data.Username = strings.TrimSpace(data.Username)
	data.Username = strings.ToLower(data.Username)
	err := conn.Where("username = ?", data.Username).First(acc).RecordNotFound()
	if err {
		return u.Message(false, "Username belum terdaftar")
	}

	if acc.AuthToken != "" {
		return u.Message(true, "restored session")
	}

	key := "IN1KEYD3CRYP7YYY"

	encrypted := u.Encrypt(key, data.Password)

	if encrypted != acc.Password {
		return u.Message(false, "username atau password salah")
	}

	ttl := 28800 * time.Second
	claims := ApiKeyAccount{
		uint64(acc.ID),
		acc.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(ttl).Unix(),
			Issuer:    "M1TR4MASNET",
			Id:        "1NI7ESM1TR4MAS",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errtoken := token.SignedString([]byte("M1TR4MAS"))
	if errtoken != nil {
		return u.Message(false, "gagal buat token")
	}

	aktivitas := &Aktivitas{
		NamaAktivitas: "login",
		AccountId:     acc.ID,
	}

	respcreate := CreateAct(aktivitas)

	if !respcreate["status"].(bool) {
		return u.Message(false, "gagal buat aktivitas")
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return u.Message(false, "gagal open transaksi")
	}

	acc.SaveToken(tx, tokenString)

	tx.Commit()

	resp := u.Message(true, "Logged In")
	resp["id_account"] = acc.ID
	return resp
}

func Logout(acc *Account) map[string]interface{} {
	if acc.AuthToken == "" {
		return u.Message(true, "anda telah logout")
	}

	conn := GetDB()
	defer conn.Close()

	aktivitas := &Aktivitas{
		NamaAktivitas: "logout",
		AccountId:     acc.ID,
	}

	resp := CreateAct(aktivitas)

	if !resp["status"].(bool) {
		return u.Message(false, "gagal buat aktivitas")
	}

	tx := conn.Begin()
	if tx.Error != nil {
		return u.Message(false, "gagal open transaksi")
	}

	errdel := acc.DeleteToken(tx)
	if errdel != nil {
		tx.Rollback()
		return u.Message(false, errdel.Error())
	}

	tx.Commit()
	return u.Message(false, "berhasil logout")
}
