package main

import apps "autumn-2021-intern-assignment/app"

func main() {
	app := apps.New()

	sync := app.Init()
	defer sync()

	closeDB, ok := app.Database()
	if !ok {
		return
	}
	defer closeDB()

	app.Start()
}
