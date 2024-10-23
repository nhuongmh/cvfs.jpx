package langfi

import (
	"context"
	"encoding/json"

	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model"
	"github.com/open-spaced-repetition/go-fsrs/v3"
)

type FSRSData struct {
	model.Base
	fsrs.Card
}

const (
	CARD_NEW     = "New"
	CARD_LEARN   = "Learn"
	CARD_DISCARD = "Discard"
	CARD_SAVE    = "Save"
)

type ReviewCard struct {
	model.Base
	FsrsData   FSRSData
	Front      string                 `json:"front"`
	Back       string                 `json:"back"`
	Properties map[string]interface{} `json:"properties"`
	Status     string                 `json:"status"`
}

func NewReviewCard(front string, back string) ReviewCard {
	return ReviewCard{
		FsrsData:   FSRSData{Card: fsrs.NewCard()},
		Front:      front,
		Back:       back,
		Properties: map[string]interface{}{},
		Status:     CARD_NEW,
	}
}

func (c *ReviewCard) SetProp(key string, value interface{}) {
	c.Properties[key] = value
}

func (c *ReviewCard) GetProp(key string) interface{} {
	return c.Properties[key]
}

func (c *ReviewCard) PropertiesToJson() string {
	props, err := json.Marshal(c.Properties)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("failed to marshal card %v properties", c.ID)
		return "{}"
	}
	return string(props)
}

func (c *ReviewCard) SetPropertiesFromJson(str string) {
	props := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &props)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("failed to unmarshal card %v properties", c.ID)
	} else {
		c.Properties = props
	}
}

type PracticeService interface {
	GetGroups(ctx context.Context) []string
	FetchCard(ctx context.Context, group string) (*ReviewCard, error)
	//newState should be Again, Hard, Good, Easy
	SubmitCard(ctx context.Context, cardID, rating uint64) error
	GetCard(ctx context.Context, cardId uint64) (*ReviewCard, error)
}

type PracticeRepo interface {
	AddCard(ctx context.Context, card *ReviewCard) error
	GetCard(ctx context.Context, cardID uint64) (*ReviewCard, error)
	UpdateCard(ctx context.Context, card *ReviewCard) error
	FetchReviewCard(ctx context.Context, groupID string) (*ReviewCard, error)
	FetchUnProcessCard(ctx context.Context, groupID string) (*ReviewCard, error)
}
