package util

import "log"

func CheckErr(err error) (ok bool) {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
