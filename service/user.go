package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
)

type IUserService interface {
	CreateUser(user *model.User) error
	//PutUser(user *model.User) error
	IsExistsByAttr(attrName, attrVal string) (bool, error)
	GetPasswordByAttr(attrName, attrVal string) (string, error)
	GetStateByAttr(attrName, attrVal string) (string, error)
	//DeleteUser(id int) error
	//GetUserByPage(page, onePage int) ([]model.User, error)
	//CheckAdminByAttr(attrName, attrVal string) (bool, error)
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

//func (us *UserService) PutUser(user *model.User) error {}

func (us *UserService) IsExistsByAttr(attrName, attrVal string) (bool, error) {
	return us.UserRepository.IsExistsByAttr(attrName, attrVal)
}

func (us *UserService) GetPasswordByAttr(attrName, attrVal string) (string, error) {
	return us.UserRepository.SelectPasswordByAttr(attrName, attrVal)
}

func (us *UserService) GetStateByAttr(attrName, attrVal string) (string, error) {
	return us.UserRepository.SelectStateByAttr(attrName, attrVal)
}

//func (us *UserService) DeleteUser(id int) error {}

//func (us *UserService) CheckAdminByAttr(attrName, attrVal string) (bool, error) {
//	return us.UserRepository.SelectIsAdminByAttr(attrName, attrVal)
//}
