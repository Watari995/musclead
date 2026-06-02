package trainingusecase

import (
	"context"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/pagination"
	trainingdomain "github.com/Watari995/musclead/internal/training/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
)

type ListTrainingsInput struct {
	UserID valueobject.UserID
	Limit  int
	Offset int
}

type ListTrainingsOutput struct {
	Trainings  []*trainingdomain.Training
	Pagination pagination.OffsetPaginator
}

type ListTrainings struct {
	trainingRepo trainingdomain.TrainingRepository
}

func (uc *ListTrainings) Execute(ctx context.Context, input ListTrainingsInput) (*ListTrainingsOutput, error) {
	trainings, paginator, err := uc.trainingRepo.FindAllByUserIDWithOffsetPagination(ctx, input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, myerror.NewInternalError().Wrap(err)
	}

	return &ListTrainingsOutput{
		Trainings:  trainings,
		Pagination: paginator,
	}, nil
}

func NewListTraining(trainingRepo trainingdomain.TrainingRepository) *ListTrainings {
	return &ListTrainings{trainingRepo: trainingRepo}
}
