package service

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
	error2 "todo/error"
	"todo/interfaces"
	"todo/logger"
	"todo/model"
)

type UserService struct {
	repo interfaces.UserRepositoryInterface
}

func NewUserService(repo interfaces.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (service *UserService) Login(username string, password string) (string, *model.Error) {
	//if service.repo.Validate(username, password) {
	//	token := GenerateJwt(username)
	//	return token, nil
	//}
	//return "", error2.NewBadRequestError("You are not authorized to access this resource")

	isValid := service.repo.Validate(username, password)
	if !isValid {
		// Log invalid credentials
		logger.Log.Error("Invalid username or password for user: ")
		return "", error2.NewBadRequestError("Invalid username or password")
	}

	// Generate JWT token
	token := GenerateJwt(username)
	//if err != nil {
	//	// Log error generating JWT
	//	logger.Log.Error("Error while generating JWT for user: ", username, " Error: ", err)
	//	return "", error2.NewInternalServerError("Internal server error")
	//}

	return token, nil
}

func GenerateJwt(username string) string {
	var secretKey = []byte("secret")

	//generate token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logger.Log.Error("Error while signing token")
		return ""
	}
	return tokenString
}
