package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
)

type IUserService interface {
	CreateUser(user *model.User) error
	DeleteUser(id int) error
	PutUser(user *model.User) error
	PutPassword(id int, password string) error
	IsExistsByAttr(attrName, attrVal string) (bool, error)
	GetPasswordByAttr(attrName, attrVal string) (string, error)
	GetStateByAttr(attrName, attrVal string) (string, error)
	GetUserByPage(page, onePageCount int, state string) ([]model.User, error)
	GetUserByID(id int) (model.User, error)
	GetCountByState(state string) (int, error)
	GetIDNameByAtr(attrName, attrVal string) (int, string, error)
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

func (us *UserService) DeleteUser(id int) error {
	return us.UserRepository.Delete(id)
}

func (us *UserService) PutUser(user *model.User) error {
	return us.UserRepository.Put(user)
}

func (us *UserService) PutPassword(id int, password string) error {
	return us.UserRepository.UpdateAttr(id, "password", password)
}

func (us *UserService) IsExistsByAttr(attrName, attrVal string) (bool, error) {
	return us.UserRepository.IsExistsByAttr(attrName, attrVal)
}

func (us *UserService) GetPasswordByAttr(attrName, attrVal string) (string, error) {
	return us.UserRepository.SelectPasswordByAttr(attrName, attrVal)
}

func (us *UserService) GetStateByAttr(attrName, attrVal string) (string, error) {
	return us.UserRepository.SelectStateByAttr(attrName, attrVal)
}

func (us *UserService) GetUserByPage(page, onePageCount int, state string) ([]model.User, error) {
	return us.UserRepository.SelectUsersByPage(page, onePageCount, state)
}

func (us *UserService) GetUserByID(id int) (model.User, error) {
	return us.UserRepository.SelectUserByID(id)
}

func (us *UserService) GetCountByState(state string) (int, error) {
	return us.UserRepository.SelectCountByState(state)
}

func (us *UserService) GetIDNameByAtr(attrName, attrVal string) (int, string, error) {
	return us.UserRepository.SelectIDNameByAtr(attrName, attrVal)
}
