package models

import (
	"fmt"
	u "tes-mitramas/utils"
	"time"

	"github.com/jinzhu/gorm"
)

type Aktivitas struct {
	ID            int64     `gorm:"primary_key;auto_increment" json:"id"`
	NamaAktivitas string    `gorm:"size:255" json:"nama_aktivitas"`
	AccountId     int64     `gorm:"not null"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type EditAktivitasJson struct {
	IdAktivitas   int64  `gorm:"not null" json:"id"`
	NamaAktivitas string `gorm:"size:255" json:"nama_aktivitas"`
}

func (Aktivitas) TableName() string {
	return "aktivitas"
}

func CreateAct(aktivitas *Aktivitas) map[string]interface{} {
	conn := GetDB()
	defer conn.Close()

	tx := conn.Begin()

	if tx.Error != nil {
		return u.Message(false, tx.Error.Error())
	}

	err := aktivitas.DoCreateAct(conn)

	if err != nil {
		tx.Rollback()
		return u.Message(false, err.Error())
	}

	tx.Commit()

	return u.Message(true, "berhasil buat aktivitas")
}

func (aktivitas *Aktivitas) DoCreateAct(conn *gorm.DB) error {
	err := conn.Create(&aktivitas).Error
	if err != nil {
		return err
	}
	return nil
}

func EditAct(aktivitas *EditAktivitasJson, accid int64) map[string]interface{} {
	if aktivitas.NamaAktivitas == "" {
		return u.Message(false, "nama aktivitas tidak boleh kosong")
	}
	conn := GetDB()
	defer conn.Close()

	act := &Aktivitas{}
	errnotfound := conn.Where("id = ? and account_id = ?", aktivitas.IdAktivitas, accid).Find(&act).RecordNotFound()

	if errnotfound {
		return u.Message(false, "aktivitas tidak ditemukan")
	}

	if act.NamaAktivitas == "login" || act.NamaAktivitas == "logout" || act.NamaAktivitas == "checkin" || act.NamaAktivitas == "checkout" {
		resperr := fmt.Sprintf("tidak bisa edit aktivitas %s", act.NamaAktivitas)
		return u.Message(false, resperr)
	}

	tx := conn.Begin()

	if tx.Error != nil {
		return u.Message(false, "gagal open transaksi")
	}

	err := tx.Model(&Aktivitas{}).Where("id = ? and account_id = ?", aktivitas.IdAktivitas, accid).Updates(map[string]interface{}{"nama_aktivitas": aktivitas.NamaAktivitas}).Error

	if err != nil {
		tx.Rollback()
		return u.Message(false, "gagal rubah aktivitas")
	}

	tx.Commit()

	return u.Message(true, "berhasil rubah aktivitas")
}

func DeleteAct(aktivitas *EditAktivitasJson, accid int64) map[string]interface{} {
	conn := GetDB()
	defer conn.Close()

	act := &Aktivitas{}
	errnotfound := conn.Where("id = ? and account_id = ?", aktivitas.IdAktivitas, accid).Find(&act).RecordNotFound()

	if errnotfound {
		return u.Message(false, "aktivitas tidak ditemukan")
	}

	if act.NamaAktivitas == "login" || act.NamaAktivitas == "logout" || act.NamaAktivitas == "checkin" || act.NamaAktivitas == "checkout" {
		resperr := fmt.Sprintf("tidak bisa hapus aktivitas %s", act.NamaAktivitas)
		return u.Message(false, resperr)
	}

	tx := conn.Begin()

	if tx.Error != nil {
		return u.Message(false, "gagal open transaksi")
	}

	err := tx.Delete(act).Error

	if err != nil {
		tx.Rollback()
		return u.Message(false, "gagal hapus aktivitas")
	}

	tx.Commit()

	return u.Message(true, "berhasil hapus aktivitas")
}
