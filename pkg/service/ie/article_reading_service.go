package ieservice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

func (ies *IEservice) GradeQuestionSubmit(ctx context.Context, articleReadingId uint64, answers map[int]string) (*ie.TestResult, error) {
	articleReading, err := ies.repo.FindArticleReadingByID(ctx, articleReadingId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article reading")
	}
	if articleReading == nil {
		return nil, errors.New("article reading not found")
	}
	if len(articleReading.Questions) != len(answers) {
		return nil, errors.New("number of answers does not match number of questions")
	}

	numberCorrect := 0
	testResult := ie.TestResult{
		ArticleReadingId: articleReading.ID,
		Score:            0.0,
		QuestionResults:  []ie.QuestionResult{},
	}
	for i := range articleReading.Questions {
		q := articleReading.Questions[i]
		userAnswer, ok := answers[int(q.ID)]
		if !ok {
			logger.Log.Warn().Msgf("answer for question %d not found", i)
			userAnswer = ""
		}
		questionResult := ie.QuestionResult{
			QuestionID: q.ID,
			Answer:     q.Answer,
			UserAnswer: userAnswer,
			Correct:    false,
		}
		if q.Type == ie.QUESTION_TYPE_MULTIPLE_CHOICE {
			if strings.EqualFold(q.Answer, userAnswer) {
				questionResult.Correct = true
			}
		} else if q.Type == ie.QUESTION_TYPE_SHORT_ANSWER {
			if strings.Contains(strings.ToLower(q.Answer), strings.ToLower(userAnswer)) {
				questionResult.Correct = true
			}
		} else if q.Type == ie.QUESTION_TYPE_TRUE_FALSE {
			if strings.EqualFold(q.Answer, userAnswer) {
				questionResult.Correct = true
			}
		} else {
			logger.Log.Warn().Msgf("unsupport question type %s", q.Type)
			if strings.EqualFold(q.Answer, userAnswer) {
				questionResult.Correct = true
			}
		}

		if questionResult.Correct {
			numberCorrect++
		}
		testResult.QuestionResults = append(testResult.QuestionResults, questionResult)
	}
	testResult.Score = float32(numberCorrect) / float32(len(articleReading.Questions)) * 100.0
	savedTestResult, err := ies.repo.SaveTestSubmission(ctx, &testResult)
	return savedTestResult, err
}

func (ies *IEservice) GetTestSubmission(ctx context.Context, id uint64) (*ie.TestResult, error) {
	testResult, err := ies.repo.GetTestSubmissionById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get test submission")
	}
	return testResult, nil
}

func (ies *IEservice) GetTestSubmissionByReadingId(ctx context.Context, readingId uint64) (*[]ie.TestResult, error) {
	testResults, err := ies.repo.FindSubmissionByReadingId(ctx, readingId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get test submission")
	}
	return testResults, nil
}

func (ies *IEservice) DeleteTestSubmission(ctx context.Context, id uint64) error {
	err := ies.repo.DeleteSubmission(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete test submission")
	}
	return nil
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
	Response in JSON format with schema as follows:
		Question = {'type': string, 'question': string, 'options': []string, 'answer': string}
		Return: Array<Question>
		Where: 
		- type: type of question, can be one of the following: 'multiple_choice', 'short_answer', 'true_false_not_given'
		- question: the question text
		- options: the options for multiple choice question, empty for other types
		- answer: the answer for the question, if the question is multiple choice, the answer should be A, B, C, D,... otherwise it is the answer text
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
