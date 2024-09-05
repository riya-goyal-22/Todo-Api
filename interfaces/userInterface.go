package interfaces

import "todo/model"

type UserRepositoryInterface interface {
	Validate(username, password string) bool
}
type UserServiceInterface interface {
	Login(username string, password string) (string, *model.Error)
}
