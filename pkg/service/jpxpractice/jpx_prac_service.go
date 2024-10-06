package jpxpractice

import (
	"time"

	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
	"github.com/open-spaced-repetition/go-fsrs/v3"
)

type jpxPracService struct {
	contextTimeout time.Duration
	repo           jp.JpxPracticeRepository
	fsrsService    *fsrs.FSRS
}

func NewJpxPracService(timeout time.Duration, repo jp.JpxPracticeRepository) jp.JpxPracticeService {
	jpa := &jpxPracService{
		contextTimeout: timeout,
		repo:           repo,
		fsrsService:    fsrs.NewFSRS(fsrs.DefaultParam()),
	}
	return jpa
}
