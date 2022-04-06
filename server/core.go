package server

import (
	"log"

	"autumn-2021-intern-assignment/database"
	"autumn-2021-intern-assignment/utils/config"
	"github.com/valyala/fasthttp"
)

func Start(db *database.DB, host config.Host) {
	log.Println("Server started at 127.0.0.1:" + host.Port)

	err := fasthttp.ListenAndServe(":"+host.Port, Handler(db))
	log.Fatalln("ListenAndServe: ", err)
}
