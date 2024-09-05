package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	error2 "todo/error"
	"todo/interfaces"
	"todo/logger"
	"todo/model"
)

type TodoHandler struct {
	service interfaces.TodoServiceInterface
}

type AuthHandler struct {
	service interfaces.UserServiceInterface
}

func (handler *TodoHandler) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	cError := CheckToken(authHeader)
	if cError != nil {
		logger.Error(cError.Message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewBadRequestError("Unauthorized"))
		return
	}
	var requestBody *model.DTO
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		if err.Error() == "EOF" {
			logger.Error("Request Body EOF(empty)")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(error2.NewBadRequestError("Body is empty").Code)
			json.NewEncoder(w).Encode(error2.NewBadRequestError("Body is empty"))
			//fmt.Fprint(w, error2.NewBadRequestError("Body is empty").Message)
		} else {
			logger.Error("Request Body JSON Error")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
			json.NewEncoder(w).Encode(error2.NewBadRequestError("Bad Request"))
		}
		return
	}
	responseBody, cError := handler.service.CreateTodo(requestBody)
	if cError != nil {
		logger.Error(cError.Message)
		w.WriteHeader(cError.Code)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseBody)
	logger.Info("Created Todo")
}

// Get all to-do items
func (handler *TodoHandler) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	cError := CheckToken(authHeader)
	if cError != nil {
		logger.Error(cError.Message) //changes here
		w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewBadRequestError("Unauthorized"))
		return
	}
	responses, error := handler.service.GetTodos()
	if error != nil {
		logger.Error(error.Message)
		w.WriteHeader(error.Code)
		fmt.Fprint(w, error.Message)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
	logger.Info("Retrieved Todos")
}

// Update a to-do item
func (handler *TodoHandler) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	cError := CheckToken(authHeader)
	if cError != nil {
		logger.Error(cError.Message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewBadRequestError("Unauthorized"))
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		logger.Error("Request Parse Error")
		w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewBadRequestError("Invalid To-do ID"))
		return
	}

	var updatedTodo *model.DTO
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		if err.Error() == "EOF" {
			logger.Error("Request Body EOF(empty)")
			w.WriteHeader(error2.NewBadRequestError("Body is empty").Code)
			json.NewEncoder(w).Encode(error2.NewBadRequestError("Body is empty"))
		} else {
			logger.Error("Request Body JSON Error")
			w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
			json.NewEncoder(w).Encode(error2.NewBadRequestError("Bad Request"))
		}
		return
	}
	cError = handler.service.UpdateTodo(id, updatedTodo)
	if cError != nil {
		logger.Error(cError.Message)
		w.WriteHeader(cError.Code)
		json.NewEncoder(w).Encode(cError)
		return
	}
	updatedTodo, cError = handler.service.GetTodoById(id)
	if cError != nil {
		logger.Error(cError.Message)
		w.WriteHeader(cError.Code)
		json.NewEncoder(w).Encode(cError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTodo)
	logger.Info("Updated Todo")
}

// Delete a to-do item
func (handler *TodoHandler) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	cError := CheckToken(authHeader)
	if cError != nil {
		logger.Error(cError.Message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewBadRequestError("Unauthorized"))
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.Error("Request Parse Error")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewBadRequestError("Bad Request"))
		return
	}
	cError = handler.service.DeleteTodo(id)
	if cError != nil {
		logger.Error(cError.Message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(cError.Code)
		json.NewEncoder(w).Encode(cError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("Deleted successfully!")
	if err != nil {
		logger.Log.Error("error in encoding:" + err.Error())
	}
	logger.Info("Deleted Todo")
}

func (handler *TodoHandler) updateTodoStatus(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	cError := CheckToken(authHeader)
	if cError != nil {
		logger.Error(cError.Message)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewBadRequestError("Unauthorized"))
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		logger.Error("Request Parse Error")
		w.WriteHeader(error2.NewBadRequestError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewBadRequestError("Bad Request"))
		return
	}
	existingDto, cError := handler.service.GetTodoById(id)
	if cError != nil {
		logger.Error(cError.Message)
		w.WriteHeader(cError.Code)
		json.NewEncoder(w).Encode(cError)
		return
	}

	// Decode the incoming partial update
	var updatedFields map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updatedFields)
	if err != nil {
		if err.Error() == "EOF" {
			logger.Error("Request Body EOF(empty)")
			w.WriteHeader(error2.NewBadRequestError("Body is empty").Code)
			json.NewEncoder(w).Encode(error2.NewBadRequestError("Body is empty"))
		}
		logger.Error("Request Body JSON Error")
		w.WriteHeader(error2.NewInternalServerError("Bad Request").Code)
		json.NewEncoder(w).Encode(error2.NewInternalServerError("Bad Request"))
		return
	}
	if title, ok := updatedFields["title"].(string); ok {
		existingDto.Title = title
	}
	if status, ok := updatedFields["status"].(string); ok {
		existingDto.Status = status
	}

	cError = handler.service.UpdateTodo(id, existingDto)
	if cError != nil {
		logger.Error(cError.Message)
		w.WriteHeader(cError.Code)
		json.NewEncoder(w).Encode(cError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingDto)
	logger.Info("Updated Todo")
}

func (handler *AuthHandler) ValidateUser(w http.ResponseWriter, r *http.Request) {
	var credentials map[string]string
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		logger.Error("Request Body JSON Error in validation handler")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(error2.NewInternalServerError("Internal Server Error").Code)
		json.NewEncoder(w).Encode(error2.NewInternalServerError("Internal Server Error"))
		return
	}

	username := credentials["username"]
	password := credentials["password"]

	jwtToken, cError := handler.service.Login(username, password)
	if cError != nil {
		logger.Error("Error in Login ,May be wrong user credentials")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(error2.NewInternalServerError("Internal Server Error").Code)
		json.NewEncoder(w).Encode(error2.NewInternalServerError("Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jwtToken)
	logger.Info("Successfully generated token")
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		var secretKey = []byte("secret")
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func CheckToken(authHeader string) *model.Error {
	if authHeader == "" {
		return error2.NewUnauthorizedError("Authorization header missing")
	}

	// Bearer token format
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate the token
	token, err := ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return error2.NewUnauthorizedError("Invalid or expired token")
	}
	return nil
}
