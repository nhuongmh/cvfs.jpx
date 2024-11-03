package jp

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nhuongmh/cfvs.jpx/pkg/model"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
)

const (
	KANA            = "kana"
	MEANING         = "meaning"
	HAN_VIE         = "han_viet"
	CATEGORY        = "category"
	MINNA           = "minna"
	MARKED_TO_LEARN = "marked_to_learn"
	FORM_VAR_REGEX  = `\[([a-zA-Z_]+[@]?[1-9]?)\]`
)

// card state

const (
	MAX_CARDS_PER_FORMULA = 8
)

type Word struct {
	model.Base
	model.Entry
	Category string
}

type SentenceFormula struct {
	Minna       string `json:"minna"`
	Form        string `json:"form"`
	Description string `json:"description"`
	Backward    string `json:"backward"`
}

func (s *SentenceFormula) IsValid() error {
	//get all vars in Form
	sentenceVarRegex := regexp.MustCompile(FORM_VAR_REGEX)
	formVars := sentenceVarRegex.FindAllStringSubmatch(s.Form, -1)
	backVars := sentenceVarRegex.FindAllStringSubmatch(s.Backward, -1)

	//all vars in backward must available in form
	fVarsMap := map[string]bool{}
	for _, fvar := range formVars {
		fVarsMap[fvar[1]] = true
	}
	for _, bvar := range backVars {
		sfVar := bvar[1]
		if _, ok := fVarsMap[sfVar]; !ok {
			return fmt.Errorf("var [%v] found in backward but not found in form", sfVar)
		}
	}
	return nil

}

// type CardProposal struct {
// 	model.Base
// 	Front     string `json:"front"`
// 	Back      string `json:"back"`
// 	State     string `json:"state"`
// 	FormulaID int    `json:"formula_id"`
// }

func NewWord(w string) Word {
	return Word{
		Entry: model.Entry{
			Name:       w,
			Properties: map[string]string{},
		},
	}
}

func (w *Word) SetProp(key, value string) {
	w.Properties[key] = value
}

func (w *Word) GetKana() string {
	return w.GetPropOrEmpty(KANA)
}

func (w *Word) GetMeaning() string {
	return w.GetPropOrEmpty(MEANING)
}

func (w *Word) GetPropOrEmpty(key string) string {
	if value, ok := w.Properties[key]; ok {
		return value
	}
	return ""
}

type JpxGeneratorService interface {
	InitData(ctx context.Context) error
	DeleteNewCards(ctx context.Context) error
	GetWordList(ctx context.Context) *[]Word
	BuildCards(ctx context.Context) (*[]langfi.ReviewCard, error)
	FetchProposal(ctx context.Context) (*langfi.ReviewCard, error)
	SubmitProposal(ctx context.Context, cardID uint64, status string) error
	// GetProcessGroups(ctx context.Context) []string
	EditCardText(ctx context.Context, newCard *langfi.ReviewCard) (*langfi.ReviewCard, error)
}

// type JpxGeneratorRepository interface {
// 	AddCardProposals(ctx context.Context, cards *[]CardProposal) (*[]CardProposal, error)
// 	ModifyCard(ctx context.Context, card *CardProposal) error
// }
