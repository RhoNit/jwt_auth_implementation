package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/RhoNit/jwt_auth_implementation/database"
	"github.com/RhoNit/jwt_auth_implementation/models"
	"github.com/RhoNit/jwt_auth_implementation/utils"
	"github.com/jackc/pgx/v5"
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
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to connect to database.. /signup"})
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

func Login(c echo.Context) error {
	// check the binding of request body to golang-structure
	var req models.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "incorrect request body"})
	}

	// authenticate the login credentials
	conn, err := database.ConnectToDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to connect to database.. /login"})
	}

	var userID int
	var hashedPasswd string

	query := "SELECT id, password FROM users WHERE email=$1"
	err = conn.QueryRow(context.Background(), query, req.Email).Scan(&userID, &hashedPasswd)
	if err == pgx.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user associated with this email not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database error: " + err.Error()})
	}

	// compare the password with the hashed one
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid username or password"})
	}

	// generate access token
	signedTokenString, err := utils.GenerateToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate a token"})
	}

	// send the token in reponse body
	return c.JSON(http.StatusOK, echo.Map{
		"message":      "access token has been generated successfully",
		"access_token": signedTokenString,
	})
}
