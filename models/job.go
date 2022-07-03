package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// delete auth token setelah 12 jam
func DeleteTokenAndCheckout() {
	// nanti tambah buat set checkin jadi false
	go func() {

		for {
			DoDeleteAndCheckout()
			time.Sleep(250 * time.Second)
		}
	}()
}

func DoDeleteAndCheckout() {
	conn := GetDB()
	defer conn.Close()

	acc := Account{}

	err := conn.Where("updated_at + interval'12 hour' < now()").Find(&acc).Updates(map[string]interface{}{"auth_token": ""}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		fmt.Println(err.Error())
	}

	err2 := conn.Where("date_part('day',updated_at) < date_part('day', now())").Find(&acc).Updates(map[string]interface{}{"check_in": false}).Error
	if err2 != nil {
		if err == gorm.ErrRecordNotFound {
			return
		}
		fmt.Println(err.Error())
	}
}
