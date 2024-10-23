package jpxpractice

import (
	"context"
	"time"

	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
	"github.com/open-spaced-repetition/go-fsrs/v3"
	"github.com/pkg/errors"
)

type jpxPracService struct {
	contextTimeout time.Duration
	repo           langfi.PracticeRepo
	fsrsService    *fsrs.FSRS
}

func NewJpxPracService(timeout time.Duration, repo langfi.PracticeRepo, env *bootstrap.Env) langfi.PracticeService {
	jpa := &jpxPracService{
		contextTimeout: timeout,
		repo:           repo,
		fsrsService:    fsrs.NewFSRS(fsrs.DefaultParam()),
	}
	return jpa
}

func (jps *jpxPracService) GetGroups(ctx context.Context) []string {
	return []string{"jp"}
}

func (jps *jpxPracService) FetchCard(ctx context.Context, groupID string) (*langfi.ReviewCard, error) {
	// get card from repo
	card, err := jps.repo.FetchReviewCard(ctx, groupID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get card")
	}

	return card, nil

}

func (jps *jpxPracService) SubmitCard(ctx context.Context, cardId uint64, rating uint64) error {
	card, err := jps.repo.GetCard(ctx, cardId)
	if err != nil {
		return errors.Wrap(err, "failed to get card")
	}

	if rating <= 0 || rating >= 5 {
		return errors.New("rating must be between 1 and 4")
	}

	schedulingCards := jps.fsrsService.Repeat(card.FsrsData.Card, time.Now())
	card.FsrsData.Card = schedulingCards[fsrs.Rating(rating)].Card
	return jps.repo.UpdateCard(ctx, card)
}

func (jps *jpxPracService) GetCard(ctx context.Context, cardID uint64) (*langfi.ReviewCard, error) {
	return jps.repo.GetCard(ctx, cardID)
}
