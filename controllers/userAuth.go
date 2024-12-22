package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/RhoNit/jwt_auth_implementation/database"
	"github.com/RhoNit/jwt_auth_implementation/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {
	// de-serialize the request body into native type
	var userRequest models.User
	if err := c.Bind(&userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "incorrect request body"})
	}

	// convert human-readable password into hashed form
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error while hashing the password"})
	}
	userRequest.Password = string(hashedPassword)

	// save the de-serialized data to the database
	query := "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id"

	conn, err := database.ConnectToDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to connect to database"})
	}
	log.Println("Connected to Database")

	err = conn.QueryRow(context.Background(), query, userRequest.FirstName, userRequest.LastName, userRequest.Email, userRequest.Password).Scan(&userRequest.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to save user into the database"})
	}

	// form the response type
	createdUserResponse := models.UserCreationResponse{
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
	}

	log.Println("Saving the user to the database")
	return c.JSON(http.StatusCreated, echo.Map{"message": "user created", "user": createdUserResponse})
}
