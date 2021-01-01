package utils

import "log"

// FailIfErr logs and terminates the app if an error occurs
func FailIfErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
