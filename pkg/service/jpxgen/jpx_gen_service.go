package jpxgen

import (
	"context"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"
	"time"

	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
	"github.com/pkg/errors"
)

type jpxService struct {
	contextTimeout time.Duration
	repo           langfi.PracticeRepo
	env            *bootstrap.Env
	ggService      *ggSheetDatasource
	wordList       *[]jp.Word
}

func NewJpxService(repo langfi.PracticeRepo, timeout time.Duration, env *bootstrap.Env) jp.JpxGeneratorService {
	jps := &jpxService{
		contextTimeout: timeout,
		repo:           repo,
		env:            env,
	}
	jps.InitData(context.Background())

	return jps
}

//parsing word list

func (jps *jpxService) InitData(ctx context.Context) error {
	ggService, err := InitNewGoogleSheetService(jps.env.GoogleKeyBase64)
	if err != nil {
		return errors.Wrap(err, "init google service failed")
	}
	logger.Log.Info().Msg("init google service success, now trying to fetch data")
	wordList, err := ggService.fetchWords(jps.env.GoogleSpreadSheetId, jps.env.GoogleWordSheetName)
	if err != nil {
		return errors.Wrap(err, "fetching data from google sheet failed")
	}

	if len(*wordList) == 0 {
		return errors.New("no data fetched from google sheet")
	}

	jps.ggService = ggService
	jps.wordList = wordList
	logger.Log.Info().Msg("tried fetching data success")

	return nil
}

func (jps *jpxService) SyncWordList(ctx context.Context) error {
	return model.ErrNotImplemented
}

func (jps *jpxService) GetWordList(ctx context.Context) *[]jp.Word {
	return jps.wordList
}

// using google service to build cards based on words and setences formula from google sheet
func (jps *jpxService) BuildCards(ctx context.Context) (*[]langfi.ReviewCard, error) {
	if !jps.checkInitialized() {
		return nil, model.ErrServiceIsNotInitialized
	}

	formulas, err := jps.ggService.fetchFormulas(jps.env.GoogleSpreadSheetId, jps.env.GoogleFormulaSheetName)
	if err != nil {
		return nil, errors.Wrap(err, "failed init sentence formula")
	}

	proposalList := []langfi.ReviewCard{}

	sentenceVarRegex := regexp.MustCompile(`\[(\w+)\]`)
	for i := range *formulas {
		formula := (*formulas)[i]
		sentence := formula.Form
		meaning := formula.Backward
		//sample formula: [Subject] は [Job] です
		//sample output sentence: わたし は せんせい です

		//parse formula to find all variables and fill it with correct word
		sentenceVars := sentenceVarRegex.FindAllStringSubmatch(formula.Form, -1)
		buildSuccess := true
		for _, svar := range sentenceVars {
			rvar := svar[1]
			w, err := jps.randomWordFromCategory(rvar)
			if err != nil {
				logger.Log.Warn().Msgf("formula: %v => error not found any word for category %v: %v", formula, rvar, err)
				buildSuccess = false
				break
			}
			logger.Log.Debug().Msgf("replacing %v , [%v] with %v", sentence, rvar, w.Name)
			sentence = strings.Replace(sentence, fmt.Sprintf("[%v]", rvar), w.Name, 1)
			meaning = strings.Replace(meaning, fmt.Sprintf("[%v]", rvar), w.GetMeaning(), 1)
		}

		if buildSuccess {
			newCard := langfi.NewReviewCard(sentence, meaning)

			proposalList = append(proposalList, newCard)
		}
	}

	if len(proposalList) == 0 {
		return nil, errors.Wrap(model.ErrNoData, "No proposal card generated")
	}

	// check if card exist -> if not, insert to db

	for i := range proposalList {
		card := proposalList[i]
		_, err := jps.repo.FetchReviewCard(ctx, card.Front)
		if err != nil {
			logger.Log.Info().Msgf("card %v not exist, trying to insert", card.Front)
			err = jps.repo.AddCard(ctx, &card)
			if err != nil {
				logger.Log.Error().Err(err).Msgf("failed to insert card %v", card.Front)
			}
		}
	}

	return &proposalList, nil
}

func (jps *jpxService) randomWordFromCategory(cat string) (*jp.Word, error) {
	wordCat := []jp.Word{}
	for i := range *jps.wordList {
		w := (*jps.wordList)[i]
		if w.Category == cat {
			wordCat = append(wordCat, w)
		}
	}

	if len(wordCat) == 0 {
		return nil, errors.Wrapf(model.ErrNoData, "no such word for category %v", cat)
	}

	return &wordCat[rand.IntN(len(wordCat))], nil
}

func (jps *jpxService) checkInitialized() bool {
	return jps.ggService != nil && jps.wordList != nil
}

func (jps *jpxService) FetchProposal(ctx context.Context) (*langfi.ReviewCard, error) {
	return jps.repo.FetchUnProcessCard(ctx, "")
}

func (jps *jpxService) SubmitProposal(ctx context.Context, cardID uint64, newStatus string) error {
	card, err := jps.repo.GetCard(ctx, cardID)
	if err != nil {
		return errors.Wrap(err, "failed to get card")
	}

	card.Status = newStatus
	return jps.repo.UpdateCard(ctx, card)
}

//generate practice sentence

//
