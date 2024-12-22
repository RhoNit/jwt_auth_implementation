package main

import (
	"log"
	"os"

	"github.com/RhoNit/jwt_auth_implementation/config"
	"github.com/RhoNit/jwt_auth_implementation/controllers"
	"github.com/labstack/echo/v4"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error while loading the ENV variables: %q", err)
	}

	// initiaize the database instance
	// _, err := database.ConnectToDB()
	// if err != nil {
	// 	log.Fatalf("Error whiile connecting to the database instance: %q", err)
	// }
	// log.Println("Connected to the Database")

	// initialize the echo engine
	e := echo.New()

	// routes
	e.POST("/signup", controllers.Signup)

	// run the server
	//
	// addr := string(os.Getenv("SERVER_HOST")) + ":" + string(os.Getenv("SERVER_PORT"))
	log.Printf("Server is running on port: %s", os.Getenv("SERVER_PORT"))
	e.Logger.Fatal(e.Start(":8080"))
}
