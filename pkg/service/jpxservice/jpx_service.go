package jpxservice

import (
	"context"
	"time"

	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
	"github.com/pkg/errors"
)

type jpxService struct {
	contextTimeout time.Duration
	jpxRepo        jp.JpxRepository
	env            *bootstrap.Env
}

func NewJpxService(repo jp.JpxRepository, timeout time.Duration, env *bootstrap.Env) jp.JpxService {
	jps := &jpxService{
		contextTimeout: timeout,
		jpxRepo:        repo,
		env:            env,
	}

	return jps
}

//parsing word list

func (jps *jpxService) InitData(ctx context.Context) error {
	ggService, err := InitNewGoogleSheetService(jps.env.GoogleKeyBase64)
	if err != nil {
		return errors.Wrap(err, "init google service failed")
	}
	wordList, err := ggService.fetchData(jps.env.GoogleSpreadSheetId, jps.env.GoogleSheetName)
	if err != nil {
		return errors.Wrap(err, "fetching data from google sheet failed")
	}

	for i := range *wordList {
		w := (*wordList)[i]
		logger.Log.Info().Msgf("word: %v, prop: %v, category: %v", w.Name, w.Properties, w.Category)
	}

	minaLessons, err := ParseMinnaLessonCfg("config/sentence_formula.yml")
	if err != nil {
		return errors.Wrap(err, "failed init sentence formula")
	}

	for i := range *minaLessons {
		lesson := (*minaLessons)[i]
		for f := range lesson.Formulas {
			formula := lesson.Formulas[f]
			logger.Log.Info().Msgf("form: %v, description: %v, backward: %v", formula.Form, formula.Description, formula.Backward)
		}
	}

	return nil
}

func (jps *jpxService) SyncWordList(ctx context.Context) error {
	return model.ErrNotImplemented
}

//generate practice sentence

//
