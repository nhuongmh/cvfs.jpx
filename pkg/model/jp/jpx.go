package jp

import (
	"context"

	"github.com/nhuongmh/cfvs.jpx/pkg/model"
)

const (
	KANA     = "kana"
	MEANING  = "meaning"
	HAN_VIE  = "han_viet"
	CATEGORY = "category"
)

// card state
const (
	CARD_PROPOSAL_NEW     = "New"
	CARD_PROPOSAL_TOLEARN = "ToLearn"
	CARD_PROPOSAL_DISCARD = "Discard"
	CARD_PROPOSAL_SAVED   = "Save"
)

const (
	MAX_CARDS_PER_FORMULA = 8
)

type Word struct {
	model.Base
	model.Entry
	Category string
}

type SentenceFormula struct {
	Minna       int    `json:"minna"`
	Form        string `json:"form"`
	Description string `json:"description"`
	Backward    string `json:"backward"`
}

type CardProposal struct {
	model.Base
	Front     string `json:"front"`
	Back      string `json:"back"`
	State     string `json:"state"`
	FormulaID int    `json:"formula_id"`
}

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
	return w.getPropOrEmpty(KANA)
}

func (w *Word) GetMeaning() string {
	return w.getPropOrEmpty(MEANING)
}

func (w *Word) getPropOrEmpty(key string) string {
	if value, ok := w.Properties[key]; ok {
		return value
	}
	return ""
}

type JpxGeneratorService interface {
	InitData(ctx context.Context) error
	SyncWordList(ctx context.Context) error
	GetWordList(ctx context.Context) *[]Word
	BuildCards(ctx context.Context) (*[]CardProposal, error)
}

type JpxGeneratorRepository interface {
	AddCardProposals(ctx context.Context, cards *[]CardProposal) (*[]CardProposal, error)
	ModifyCard(ctx context.Context, card *CardProposal) error
}

type JpxPracticeService interface {
}

type JpxPracticeRepository interface {
}
