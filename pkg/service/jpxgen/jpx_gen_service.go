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

const (
	MOST_CARD_PER_FORMULA = 5
)

type jpxService struct {
	contextTimeout time.Duration
	repo           langfi.PracticeRepo
	env            *bootstrap.Env
	ggService      *ggSheetDatasource
	wordList       *[]jp.Word
	formulaList    *[]jp.SentenceFormula
}

var SENTENCE_VAR_REGEX = regexp.MustCompile(jp.FORM_VAR_REGEX)

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

	jps.ggService = ggService
	return nil
}

func (jps *jpxService) SyncGoogleSheet(ctx context.Context) error {
	wordList, err := jps.ggService.fetchWords(jps.env.GoogleSpreadSheetId, jps.env.GoogleWordSheetName)
	if err != nil {
		return errors.Wrap(err, "fetching data from google sheet failed")
	}

	if len(*wordList) == 0 {
		return errors.New("no data fetched from google sheet")
	}

	jps.wordList = wordList
	logger.Log.Info().Msgf("fetched %v words from google sheet %v", len(*jps.wordList), jps.env.GoogleSpreadSheetId)

	formulas, err := jps.ggService.fetchFormulas(jps.env.GoogleSpreadSheetId, jps.env.GoogleFormulaSheetName)
	if err != nil {
		return errors.Wrap(err, "failed init sentence formula")
	}
	jps.formulaList = formulas
	logger.Log.Info().Msgf("fetched %v formula from google sheet %v", len(*jps.formulaList), jps.env.GoogleSpreadSheetId)

	return nil
}

func (jps *jpxService) GetWordList(ctx context.Context) *[]jp.Word {
	return jps.wordList
}

func (jps *jpxService) DeleteNewCards(ctx context.Context) error {
	return jps.repo.DeleteNewCard(ctx)
}

// using google service to build cards based on words and setences formula from google sheet
func (jps *jpxService) BuildCards(ctx context.Context) (*[]langfi.ReviewCard, error) {
	if !jps.checkInitialized() {
		return nil, model.ErrServiceIsNotInitialized
	}

	err := jps.SyncGoogleSheet(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sync data from google sheet failed")
	}

	proposalList, err := jps.genMinnaCards()
	if err != nil {
		return nil, errors.Wrap(err, "build cards failed")
	}

	if len(*proposalList) == 0 {
		return nil, errors.Wrap(model.ErrNoData, "No proposal card generated")
	} else {
		logger.Log.Info().Msgf("Successfully built %v cards", len(*proposalList))
	}

	// check if card exist -> if not, insert to db

	for i := range *proposalList {
		card := (*proposalList)[i]
		existed, err := jps.repo.GetCardByFront(ctx, card.Front)
		if err != nil || len(*existed) == 0 {
			logger.Log.Info().Msgf("card %v not exist, trying to insert", card.Front)
			err = jps.repo.AddCard(ctx, &card)
			if err != nil {
				logger.Log.Error().Err(err).Msgf("failed to insert card %v", card.Front)
			}
		} else {
			logger.Log.Warn().Msgf("card %v is already existed, skip inserting into database", card.Front)
		}
	}

	return proposalList, nil
}

func (jps *jpxService) genMinnaCards() (*[]langfi.ReviewCard, error) {
	//fetch words -> gen words card -> get formula by word's lessons -> gen sentence cards
	// gen word without lesson
	// gen sentence for formula that not associate with any lesson
	proposalList := []langfi.ReviewCard{}
	minnaList := jps.getMinnaList()
	logger.Log.Debug().Msgf("Minna lesson list: %v", minnaList)
	for _, minna := range minnaList {
		// wordCards, err := jps.buildWordCards(minna)
		// if err != nil {
		// 	logger.Log.Warn().Err(err).Msg("failed to build word cards")
		// } else {
		// 	proposalList = append(proposalList, *wordCards...)
		// }
		sentences, err := jps.buildSentenceCards(minna)
		if err != nil {
			logger.Log.Warn().Err(err).Msg("failed to build sentence cards")
		} else {
			proposalList = append(proposalList, *sentences...)
		}
	}

	//build cards without minna
	wordCards, err := jps.buildWordCards("")
	if err != nil {
		logger.Log.Warn().Err(err).Msg("failed to build word cards")
	} else {
		proposalList = append(proposalList, *wordCards...)
	}
	sentences, err := jps.buildSentenceCards("")
	if err != nil {
		logger.Log.Warn().Err(err).Msg("failed to build sentence cards")
	} else {
		proposalList = append(proposalList, *sentences...)
	}

	return &proposalList, nil
}

func (jps *jpxService) buildSentenceCards(minna string) (*[]langfi.ReviewCard, error) {

	proposalList := []langfi.ReviewCard{}
	sentenceVarRegex := regexp.MustCompile(jp.FORM_VAR_REGEX)
	for i := range *jps.formulaList {
		if (*jps.formulaList)[i].Minna != minna {
			continue
		}
		for c := 0; c < MOST_CARD_PER_FORMULA; c++ {
			formula := (*jps.formulaList)[i]
			sentence := formula.Form
			meaning := formula.Backward

			err := formula.IsValid()
			if err != nil {
				logger.Log.Warn().Err(err).Msgf("formula: %v => is invalid", formula)
				continue
			}

			//sample formula: [Subject] は [Job] です
			//sample output sentence: わたし は せんせい です

			//parse formula to find all variables and fill it with correct word
			sentenceVars := sentenceVarRegex.FindAllStringSubmatch(formula.Form, -1)
			buildSuccess := true
			collectiveProps := map[string]interface{}{}
			for _, svar := range sentenceVars {
				rvar := svar[1] //raw sentence var
				prsvar := rvar  //to process raw sentence var
				// var can contains id like [Job_1]
				if strings.Contains(prsvar, "@") {
					rvarSplit := strings.Split(prsvar, "@")
					prsvar = rvarSplit[0]
				}
				w, err := jps.randomWordFromCategory(prsvar)
				if err != nil {
					logger.Log.Warn().Msgf("formula: %v => error not found any word for category %v: %v", formula, prsvar, err)
					buildSuccess = false
					break
				}
				logger.Log.Debug().Msgf("replacing %v , [%v] with %v", sentence, rvar, w.Name)
				sentence = strings.Replace(sentence, fmt.Sprintf("[%v]", rvar), w.Name, 1)
				meaning = strings.Replace(meaning, fmt.Sprintf("[%v]", rvar), w.GetMeaning(), -1)
				for k, v := range w.Properties {
					collectiveProps[k] = v
				}
			}

			if buildSuccess {
				newCard := langfi.NewReviewCard(sentence, meaning)
				newCard.Properties = collectiveProps
				if minna != "" {
					newCard.Group = minna
				} else {
					newCard.Group = "NA"
				}
				proposalList = append(proposalList, newCard)
			}
		}
	}

	return &proposalList, nil
}

func (jps *jpxService) buildWordCards(minna string) (*[]langfi.ReviewCard, error) {
	proposalList := []langfi.ReviewCard{}

	for i := range *jps.wordList {
		word := (*jps.wordList)[i]
		if minna != word.GetPropOrEmpty(jp.MINNA) {
			continue
		}
		if tolearn, ok := word.Properties[jp.MARKED_TO_LEARN]; !ok || tolearn == "" {
			continue
		}
		newCard := langfi.NewReviewCard(word.Name, word.GetMeaning())
		for k, v := range word.Properties {
			newCard.SetProp(k, v)
		}
		proposalList = append(proposalList, newCard)
	}

	return &proposalList, nil
}

// because words may contain variables, we need to replace them with correct word
func (jps *jpxService) processWord(word *jp.Word) (*jp.Word, error) {
	//copy word to avoid changing original word
	processedWord := *word
	meaning := processedWord.GetMeaning()
	vars := SENTENCE_VAR_REGEX.FindAllStringSubmatch(word.Name, -1)
	collectiveProps := map[string]interface{}{}
	for _, svar := range vars {
		rvar := svar[1]
		w, err := jps.randomWordFromCategory(rvar)
		if err != nil {
			return nil, errors.Wrapf(err, "formula: %v => error not found any word for category %v: %v", processedWord.Name, rvar, err)
		}
		logger.Log.Debug().Msgf("replacing %v , [%v] with %v", processedWord.Name, rvar, w.Name)
		processedWord.Name = strings.Replace(processedWord.Name, fmt.Sprintf("[%v]", rvar), w.Name, 1)
		meaning = strings.Replace(meaning, fmt.Sprintf("[%v]", rvar), w.GetMeaning(), 1)
		for k, v := range w.Properties {
			collectiveProps[k] = v
		}
	}
	processedWord.SetProp(jp.MEANING, meaning)
	return &processedWord, nil
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

	word := &wordCat[rand.IntN(len(wordCat))]
	return jps.processWord(word)
}

func (jps *jpxService) getMinnaList() []string {
	minnaMap := map[string]bool{}
	for i := range *jps.wordList {
		w := (*jps.wordList)[i]
		minna := w.GetPropOrEmpty(jp.MINNA)
		if minna != "" {
			minnaMap[minna] = true
		}
	}

	minnas := make([]string, 0, len(minnaMap))

	for key, _ := range minnaMap {
		minnas = append(minnas, key)
	}
	return minnas
}

func (jps *jpxService) checkInitialized() bool {
	return jps.ggService != nil
}

func (jps *jpxService) FetchProposal(ctx context.Context, group string) (*langfi.ReviewCard, error) {
	return jps.repo.FetchUnProcessCard(ctx, group)
}

func (jps *jpxService) SubmitProposal(ctx context.Context, cardID uint64, newStatus string) error {
	card, err := jps.repo.GetCard(ctx, cardID)
	if err != nil {
		return errors.Wrap(err, "failed to get card")
	}

	for _, stat := range langfi.ALL_CARD_STATUS {
		if strings.EqualFold(newStatus, stat) {
			card.Status = stat
			return jps.repo.UpdateCard(ctx, card)
		}
	}
	return errors.Errorf("invalid status %v", newStatus)
}

func (jps *jpxService) EditCardText(ctx context.Context, newCard *langfi.ReviewCard) (*langfi.ReviewCard, error) {
	card, err := jps.repo.GetCard(ctx, newCard.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get card with id=%v", newCard.ID)
	}

	card.Front = newCard.Front
	card.Back = newCard.Back
	for k, v := range newCard.Properties {
		card.SetProp(k, v)
	}
	return card, jps.repo.UpdateCard(ctx, card)
}

//generate practice sentence

//
