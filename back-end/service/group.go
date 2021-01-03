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
	GetGroupByID(id int) (model.Group, error)
	GetGroupsByCreatorAndPage(creator, page, onePageCount int) ([]model.Group, error)
	GetAllGroupsByCreator(creator int) ([]model.Group, error)
	GetGroupCountByCreator(creator int) (int, error)
	GetMembersByGroups(ids string) ([]model.Group, error)
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

func (gs *GroupService) GetGroupByID(id int) (model.Group, error) {
	return gs.GroupRepository.SelectGroupByID(id)
}

func (gs *GroupService) GetGroupsByCreatorAndPage(creator, page, onePageCount int) ([]model.Group, error) {
	return gs.GroupRepository.SelectGroupsByCreator(creator, page, onePageCount)
}

func (gs *GroupService) GetAllGroupsByCreator(creator int) ([]model.Group, error) {
	return gs.GroupRepository.SelectGroupsByCreator(creator)
}

func (gs *GroupService) GetGroupCountByCreator(creator int) (int, error) {
	return gs.GroupRepository.SelectGroupCountByCreator(creator)
}

func (gs *GroupService) GetMembersByGroups(ids string) (groups []model.Group, err error) {
	return gs.GroupRepository.SelectMembersByGroups(ids)
}
