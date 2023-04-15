package main

import (
	"final-project-mygram/database"
	"final-project-mygram/router"
)

func main() {
	database.StartDB()
	// var PORT = os.Getenv("PORT")
	router.StartApp().Run(":8080")
}
