package root

import "streakr-backend/pkg/http"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Init() {
	db := DbConnect()
	http.InitHttp(db)
}