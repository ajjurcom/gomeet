package service

import "com/mittacy/gomeet/repository"

type IEmailService interface {
	GetContent(string) (string, error)
}

func NewEmailService(repository repository.IEmailRepository) IEmailService {
	return &EmailService{repository}
}

type EmailService struct {
	EmailRepository repository.IEmailRepository
}

func (es *EmailService) GetContent(name string) (string, error) {
	return es.EmailRepository.GetContentByName(name)
}
