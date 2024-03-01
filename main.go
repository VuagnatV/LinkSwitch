package main

import (
	"linkSwitch/database"
	"linkSwitch/services"
	"log"
	"net/http"
)

func main() {
	db, err := database.InitMongo("mongodb://localhost:27017")
	address := ":8000"
	router := services.RunServer(db)

	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err)
	}

}
