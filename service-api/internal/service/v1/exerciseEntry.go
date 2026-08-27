package service

import (
	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
)

type ExerciseEntryService struct {
	repo repository.ExerciseEntryRepository
}

func NewExerciseEntryService(repo repository.ExerciseEntryRepository) *ExerciseEntryService {
	return &ExerciseEntryService{repo: repo}
}

//ExerciseEntries
/*
func (s *ExerciseEntryService) Get(c *gin.Context, req models.ExerciseEntryGetRequest) (models.ExerciseEntryGetResponse, error) {

}

func (s *ExerciseEntryService) Patch(c *gin.Context, req models.ExerciseEntryPatchRequest) (models.ExerciseEntryPatchResponse, error) {

}

func (s *ExerciseEntryService) Delete(c *gin.Context, req models.ExerciseEntryDeleteRequest) error {

}*/
