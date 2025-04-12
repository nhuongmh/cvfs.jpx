package ieservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/pkg/errors"
)

func (ies *IEservice) GenArticleReading(ctx context.Context, article *ie.Article, force bool) (*ie.ArticleReading, error) {
	articleRds, err := ies.repo.FindReadingByArticleId(ctx, article.ID)
	if err == nil && articleRds != nil && !force {
		logger.Log.Warn().Msg("article reading already exists")
		return articleRds, nil
	}
	//delete existing article reading
	if articleRds != nil {
		err = ies.repo.DeleteArticleReading(ctx, articleRds.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to delete existing article reading")
		}
	}

	questions, err := ies.GenerateQuestion(ctx, article.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate questions")
	}
	articleReading := &ie.ArticleReading{
		ArticleID: article.ID,
		Status:    ie.ARTICLE_NEW,
		Score:     0,
		Questions: *questions,
	}
	return ies.repo.SaveArticleReading(ctx, articleReading)
}

func (ies *IEservice) GetArticleReading(ctx context.Context, articleId uint64) (*ie.ArticleReading, error) {
	articleReading, err := ies.repo.FindReadingByArticleId(ctx, articleId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article reading")
	}
	return articleReading, nil
}

func (ies *IEservice) GenerateQuestion(ctx context.Context, id uint64) (*[]ie.Question, error) {
	if ies.gemi == nil {
		return nil, errors.New("gemini client is not initialized")
	}

	article, err := ies.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article")
	}

	if article.Content == "" {
		return nil, errors.New("article content is empty")
	}

	prompt := fmt.Sprintf(`
	From below article, help me generate 13 questions with answer in IELTS style. Question type such as
		- Multiple choice
		- Short answer
		- True/False/Not Given
		- Matching headings
	Response in JSON format with schema as follows:
		Question = {'type': string, 'question': string, 'options': []string, 'answer': string, 'headings': []string, 'paragraph': string}
		Return: Array<Question>
		Where: 
		- type: type of question, can be one of the following: 'multiple_choice', 'short_answer', 'true_false_not_given', 'matching_headings'
		- question: the question text
		- options: the options for multiple choice question, empty for other types
		- answer: the answer for the question
		- headings: the headings for matching headings question, empty for other types
		- paragraph: the paragraph for matching headings question, empty for other types
	-----
	Title: %v
	Article: %v
	`, article.Title, article.Content)

	resp, err := ies.gemi.GenerateContentSimp(context.Background(), prompt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate content")
	}
	var questions []ie.Question
	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				if err := json.Unmarshal([]byte(txt), &questions); err != nil {
					logger.Log.Error().Err(err).Msg("Failed to unmarshal AI generated questions")
				}
				if len(questions) > 0 {
					break
				}
			}
		}
		if len(questions) > 0 {
			break
		}
	}
	if len(questions) == 0 {
		return nil, errors.New("failed to generate questions")
	}
	for i := range questions {
		questions[i].ArticleReadingId = id
		questions[i].ID = uint64(i + 1)
	}

	return &questions, nil
}
