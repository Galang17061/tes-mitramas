package models

import (
	"errors"
	"fmt"
	"strings"
	u "tes-mitramas/utils"
	"time"

	"github.com/jinzhu/gorm"
)

type Account struct {
	ID        int64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255" json:"name"`
	Username  string    `gorm:"size:255" json:"username"`
	Password  string    `gorm:"size:255" json:"password"`
	CheckIn   bool      `gorm:"default:false" json:"checkin"`
	AuthToken string    `gorm:"size:255" json:"auth_token"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type RegisterAccountJson struct {
	ID        int64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255" json:"name"`
	Username  string    `gorm:"size:255" json:"username"`
	Password  string    `gorm:"size:255" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Account) TableName() string {
	return "account"
}

func (acc *RegisterAccountJson) CreateAccount() map[string]interface{} {
	if acc.Username == "" && acc.Password == "" || acc.Name == "" {
		return u.Message(false, "data tidak valid")
	}

	if len(strings.TrimSpace(acc.Username)) == 0 || len(strings.TrimSpace(acc.Password)) == 0 || len(strings.TrimSpace(acc.Name)) == 0 {
		return u.Message(false, "data tidak valid")
	}

	notexist := GetAccountByUsername(acc.Username)
	if notexist {
		key := "IN1KEYD3CRYP7YYY"
		acc.Username = strings.TrimSpace(acc.Username)
		acc.Username = strings.ToLower(acc.Username)

		newacc := Account{
			Name:     acc.Name,
			Username: acc.Username,
			Password: u.Encrypt(key, acc.Password),
		}

		conn := GetDB()
		defer conn.Close()

		tx := conn.Begin()
		if tx.Error != nil {
			return u.Message(false, "gagal open transaksi")
		}

		err := tx.Create(&newacc).Error
		if err != nil {
			tx.Rollback()
			return u.Message(false, "gagal buat akun baru")
		}

		tx.Commit()
		return u.Message(true, "sukses buat akun baru")
	}

	return u.Message(false, "username sudah digunakan")
}

func GetAccount(username string, accid int64) *Account {
	conn := GetDB()
	defer conn.Close()
	acc := &Account{}

	err := conn.Where("id = ? and username = ?", accid, username).Find(&acc).RecordNotFound()

	if err {
		return nil
	}

	return acc
}

func GetAccountByUsername(username string) bool {
	conn := GetDB()
	defer conn.Close()
	acc := &Account{}

	isexist := conn.Where("username = ?", username).Find(&acc).RecordNotFound()

	return isexist
}

func GetToken(username string, accid int64) map[string]interface{} {
	conn := GetDB()
	defer conn.Close()
	acc := &Account{}

	conn.Where("username = ? and account_id = ?", username, accid).Find(&acc)
	resp := map[string]interface{}{}
	resp["status"] = true
	resp["data"] = acc.AuthToken
	return resp
}

func (acc *Account) SaveToken(conn *gorm.DB, token string) error {
	conn.Where("username = ?", acc.Username).Find(&acc)

	err := conn.Model(&Account{}).Where("id=?", acc.ID).Updates(map[string]interface{}{"auth_token": token}).Error
	if err != nil {
		return errors.New("gagal simpan token")
	}
	return nil
}

func (acc *Account) DeleteToken(conn *gorm.DB) error {
	err := conn.Where("username = ?", acc.Username).Find(&acc).Updates(map[string]interface{}{"auth_token": ""}).Error
	if err != nil {
		return errors.New("gagal hapus token")
	}
	return nil
}

func (acc *Account) DoCheckin() error {
	if acc.AuthToken == "" {
		return errors.New("anda telah logout, harap login kembali")
	}

	conn := GetDB()
	defer conn.Close()

	// untuk cek apakah checkin ini adalah yang pertama kali atau tidak
	// jika ya maka tambah record ke absensi
	absensi := &Absensi{}

	filter1 := "date_part('day', created_at) = date_part('day', now())"
	filter2 := "date_part('month', created_at) = date_part('month', now())"

	errfind := conn.Where(fmt.Sprintf("%s and %s", filter1, filter2)).Find(&absensi).RecordNotFound()

	aktivitas := &Aktivitas{
		NamaAktivitas: "checkin",
		AccountId:     acc.ID,
	}

	resp := CreateAct(aktivitas)

	if !resp["status"].(bool) {
		return errors.New("gagal buat aktivitas")
	}

	tx := conn.Begin()

	if tx.Error != nil {
		return errors.New("gagal open transaksi")
	}

	err := tx.Model(&Account{}).Where("id=?", acc.ID).Updates(map[string]interface{}{"check_in": true}).Error

	if err != nil {
		return errors.New("gagal checkin")
	}

	if errfind {
		newabsen := &Absensi{
			AccountId: acc.ID,
		}
		err := tx.Create(&newabsen).Error
		if err != nil {
			tx.Rollback()
			return errors.New("gagal buat absen baru")
		}
	}

	tx.Commit()
	return nil
}

func (acc *Account) DoCheckout() error {
	if acc.AuthToken == "" {
		return errors.New("anda telah logout, harap login kembali")
	}

	conn := GetDB()
	defer conn.Close()

	aktivitas := &Aktivitas{
		NamaAktivitas: "checkout",
		AccountId:     acc.ID,
	}

	resp := CreateAct(aktivitas)

	if !resp["status"].(bool) {
		return errors.New("gagal buat aktivitas")
	}

	tx := conn.Begin()

	if tx.Error != nil {
		return errors.New("gagal open transaksi")
	}

	err := tx.Model(&Account{}).Where("id=?", acc.ID).Updates(map[string]interface{}{"check_in": false}).Error

	if err != nil {
		return errors.New("gagal checkout")
	}

	tx.Commit()
	return nil
}
