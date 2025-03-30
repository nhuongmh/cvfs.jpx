package ieservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/database/postgresdb"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	ierepo "github.com/nhuongmh/cfvs.jpx/pkg/service/ie/repo"
	"github.com/nhuongmh/cfvs.jpx/pkg/service/llm/gemini"
	"github.com/pkg/errors"
)

type IEservice struct {
	contextTimeout time.Duration
	repo           *ierepo.IErepo
	env            *bootstrap.Env
	gemi           *gemini.GoogleAI
}

func NewIEservice(timeout time.Duration, env *bootstrap.Env, db *postgresdb.DB) *IEservice {
	ies := &IEservice{
		contextTimeout: timeout,
		repo:           ierepo.NewIeRepo(db),
		env:            env,
	}
	gemi, err := gemini.NewGoogleAI(ies.env.GoogleAIKey)
	if err != nil {
		logger.Log.Warn().Err(err).Msg("failed to create gemini client")
	} else {
		ies.gemi = gemi
		logger.Log.Info().Msg("successfully created gemini client")
	}
	return ies
}

func (ies *IEservice) SaveArticle(ctx context.Context, article *ie.Article) (*ie.Article, error) {
	article, err := ies.repo.Save(ctx, article)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save article")
	}

	return article, nil
}

func (ies *IEservice) GetArticle(ctx context.Context, id uint64) (*ie.Article, error) {
	article, err := ies.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article")
	}

	return article, nil
}

func (ies *IEservice) GetAllArticles(ctx context.Context, pageSize, page uint64) ([]*ie.Article, int, error) {
	articles, total, err := ies.repo.FindAll(ctx, pageSize, page)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get all articles")
	}

	return articles, total, nil
}

// find article by title
func (ies *IEservice) FindArticleByTitle(ctx context.Context, title string) ([]*ie.Article, error) {
	article, err := ies.repo.FindByTitle(ctx, title)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article by title")
	}

	return article, nil
}

// delete article by id
func (ies *IEservice) DeleteArticle(ctx context.Context, id uint64) error {
	err := ies.repo.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete article")
	}

	return nil
}

func (ies *IEservice) GenerateQuestion(ctx context.Context, id uint64) (*ie.ArticleReading, error) {
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
		Question = {'type': string, 'question_text': string, 'options': []string, 'answer': string, 'headings': []string, 'paragraph': string}
		Return: Array<Question>
		Where: 
		- type: type of question, can be one of the following: 'multiple_choice', 'short_answer', 'true_false_not_given', 'matching_headings'
		- question_text: the question text
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
	articleReading := &ie.ArticleReading{
		ArticleID: article.ID,
		Questions: questions,
		Status:    ie.ARTICLE_NEW,
		Score:     0,
	}

	return articleReading, nil
}
