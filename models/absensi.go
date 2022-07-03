package models

import (
	"errors"
	"time"
)

type Absensi struct {
	ID        int64     `gorm:"primary_key;auto_increment" json:"id"`
	AccountId int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Absensi) TableName() string {
	return "absensi"
}

func RiwayatAbsensi(username string, accid int64) ([]*Absensi, error) {
	conn := GetDB()
	defer conn.Close()

	acc := GetAccount(username, accid)

	if acc == nil {
		return nil, errors.New("akun tidak ditemukan")
	}

	if acc.AuthToken == "" {
		return nil, errors.New("anda telah logout, harap login kembali")
	}

	absenlists := make([]*Absensi, 0)

	err := conn.Where("account_id = ? and date_part('day',created_at) = date_part('day',now())", accid).Find(&absenlists).RecordNotFound()

	if err {
		return nil, errors.New("anda belum checkin hari ini")
	}

	return absenlists, nil
}

func RiwayatAbsensiPerTanggal(username string, accid int64, starttime time.Time, endtime time.Time) ([]*Absensi, error) {
	conn := GetDB()
	defer conn.Close()

	acc := GetAccount(username, accid)

	if acc == nil {
		return nil, errors.New("akun tidak ditemukan")
	}

	if acc.AuthToken == "" {
		return nil, errors.New("anda telah logout, harap login kembali")
	}

	absenlists := make([]*Absensi, 0)

	err := conn.Where("account_id = ? and created_at >= ? and created_at <= ?", accid, starttime, endtime).Find(&absenlists).RecordNotFound()

	if err {
		return nil, errors.New("anda belum checkin hari ini")
	}

	return absenlists, nil
}
