package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
)

type IUserService interface {
	CreateUser(user *model.User) error
	GetPasswordByAttr(attrName, attrVal string) (string, error)
	IsExistsByAttr(attrName, attrVal string) (bool, error)
	GetStateByAttr(attrName, attrVal string) (string, error)
	CheckAdminByAttr(attrName, attrVal string) (bool, error)
}

func NewUserService(repository repository.IUserRepository) IUserService {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repository.IUserRepository
}

func (us *UserService) CreateUser(user *model.User) error {
	return us.UserRepository.Add(user)
}

func (us *UserService) GetPasswordByAttr(attrName, attrVal string) (string, error) {
	return us.UserRepository.SelectPasswordByAttr(attrName, attrVal)
}

func (us *UserService) IsExistsByAttr(attrName, attrVal string) (bool, error) {
	return us.UserRepository.IsExistsByAttr(attrName, attrVal)
}

func (us *UserService) GetStateByAttr(attrName, attrVal string) (string, error) {
	return us.UserRepository.SelectStateByAttr(attrName, attrVal)
}

func (us *UserService) CheckAdminByAttr(attrName, attrVal string) (bool, error) {
	return us.UserRepository.SelectIsAdminByAttr(attrName, attrVal)
}
