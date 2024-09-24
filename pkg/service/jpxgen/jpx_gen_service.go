package jpxgen

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
	jpxRepo        jp.JpxGeneratorRepository
	env            *bootstrap.Env
}

func NewJpxService(repo jp.JpxGeneratorRepository, timeout time.Duration, env *bootstrap.Env) jp.JpxGeneratorService {
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
	wordList, err := ggService.fetchWords(jps.env.GoogleSpreadSheetId, jps.env.GoogleWordSheetName)
	if err != nil {
		return errors.Wrap(err, "fetching data from google sheet failed")
	}

	for i := range *wordList {
		w := (*wordList)[i]
		logger.Log.Info().Msgf("word: %v, prop: %v, category: %v", w.Name, w.Properties, w.Category)
	}

	formulas, err := ggService.fetchFormulas(jps.env.GoogleSpreadSheetId, jps.env.GoogleFormulaSheetName)
	if err != nil {
		return errors.Wrap(err, "failed init sentence formula")
	}

	for i := range *formulas {
		formula := (*formulas)[i]
		logger.Log.Info().Msgf("form: %v, description: %v, backward: %v", formula.Form, formula.Description, formula.Backward)
	}

	return nil
}

func (jps *jpxService) SyncWordList(ctx context.Context) error {
	return model.ErrNotImplemented
}

func (jps *jpxService) GenSentences(ctx context.Context) error {

	return model.ErrNotImplemented
}

//generate practice sentence

//
