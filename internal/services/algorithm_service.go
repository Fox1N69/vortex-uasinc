package service

import "test-task/internal/repository"

type AlgorithmService interface {
}

type algorithmService struct {
	repository repository.AlgorithmRepository
}

func NewAlgorithmService(algorithmRepo repository.AlgorithmRepository) AlgorithmService {
	return &algorithmService{repository: algorithmRepo}
}
