package gemini

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type GoogleAI struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGoogleAI(apiKey string) (*GoogleAI, error) {
	ctx := context.Background()
	logger.Log.Info().Msgf("Initializing google ai service with API key")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, errors.Wrap(err, "Failed init google ai client")
	}

	model := client.GenerativeModel("gemini-2.5-pro-exp-03-25")
	model.ResponseMIMEType = "application/json"

	return &GoogleAI{
		client: client,
		model:  model,
	}, nil
}

func (g *GoogleAI) GenerateContent(ctx context.Context, expectedType *genai.Schema, prompt string) (*genai.GenerateContentResponse, error) {
	g.model.ResponseSchema = expectedType
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, errors.Wrap(err, "Failed generate content")
	}
	return resp, nil
}

func (g *GoogleAI) GenerateContentSimp(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, errors.Wrap(err, "Failed generate content")
	}
	return resp, nil
}
