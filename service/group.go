package service

import (
	"com/mittacy/gomeet/model"
	"com/mittacy/gomeet/repository"
)

type IGroupService interface {
	CreateGroup(group model.Group) error
	DeleteGroup(id int) error
	PutName(group model.Group) error
	PutMember(group model.Group) error
}

func NewGroupService(repository repository.IGroupRepository) IGroupService {
	return &GroupService{repository}
}

type GroupService struct {
	GroupRepository repository.IGroupRepository
}

func (gs *GroupService) CreateGroup(group model.Group) error {
	return gs.GroupRepository.Add(group)
}

func (gs *GroupService) DeleteGroup(id int) error {
	return gs.GroupRepository.Delete(id)
}

func (gs *GroupService) PutName(group model.Group) error {
	return gs.GroupRepository.PutName(group)
}

func (gs *GroupService) PutMember(group model.Group) error {
	return gs.GroupRepository.PutMember(group)
}
