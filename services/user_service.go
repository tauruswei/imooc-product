package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"imooc-product/datamodels"
	"imooc-product/repositories"
)

type UserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

func NewUserServiceManager(repository repositories.UserRepository) UserService {
	return &UserServiceManager{repository}
}

type UserServiceManager struct {
	UserRepository repositories.UserRepository
}

func (u *UserServiceManager) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {

	user, err := u.UserRepository.Select(userName)

	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)

	if !isOk {
		return &datamodels.User{}, false
	}

	return
}

func (u *UserServiceManager) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil

}
