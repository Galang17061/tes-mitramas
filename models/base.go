package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func BaseInit() {
	conn := GetDB()
	defer conn.Close()

	conn.AutoMigrate(&Account{}, &Absensi{}, &Aktivitas{})
	insertDefaultAccount(conn)
	insertDefaultAbsensi(conn)
	insertDefaultAktivitas(conn)
}

func GetDB() *gorm.DB {
	username := "root2"
	password := "root2"
	dbName := "tesmitramasdb"
	dbHost := "localhost"

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}
	return conn
}

func insertDefaultAccount(conn *gorm.DB) {
	defaultAccount := []string{
		"login", "logout", "checkout",
	}
	defaultUsername := []string{
		"kikep08", "rudiyanto", "ahmad",
	}

	var acc Account
	var count int64

	conn.Find(&acc).Count(&count)
	if count == 0 {
		for i, val := range defaultAccount {
			insertAccount := &Account{
				Name:     val,
				Username: defaultUsername[i],
				Password: "12345678",
			}
			conn.Create(&insertAccount)
		}
	}
}

func insertDefaultAbsensi(conn *gorm.DB) {
	var acc Absensi
	var count int64

	conn.Find(&acc).Count(&count)
	if count == 0 {
		for i := 0; i < 4; i++ {
			absenList := Absensi{
				AccountId: int64(i),
			}
			conn.Create(&absenList)
		}
	}
}

func insertDefaultAktivitas(conn *gorm.DB) {
	defaultAktivitas := []string{
		"login", "logout", "checkout",
	}

	var act Aktivitas
	var count int64

	conn.Find(&act).Count(&count)
	if count == 0 {
		for i, val := range defaultAktivitas {
			insertAktivitas := &Aktivitas{
				NamaAktivitas: val,
				AccountId:     int64(i + 1),
			}
			conn.Create(&insertAktivitas)
		}
	}
}
