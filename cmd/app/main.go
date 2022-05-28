package main

import apps "autumn-2021-intern-assignment/app"

// @title           Balance API
// @version         1.0
// @description     Test task for the position of trainee golang backend developer
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/illiafox
// @contact.email  illiadimura@gmail.com

// @license.name  Boost Software License - Version 1.0
// @license.url   https://opensource.org/licenses/BSL-1.0

// @host      localhost:8080
// @BasePath  /api
func main() {
	app, sync := apps.Init()
	defer sync()

	closeDB, ok := app.Database()
	if !ok {
		return
	}
	defer closeDB()

	app.Start()
}
