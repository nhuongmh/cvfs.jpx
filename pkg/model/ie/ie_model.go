package ie

import "github.com/nhuongmh/cfvs.jpx/pkg/model"

type Article struct {
	model.Base
	Title       string `json:"title"`
	Content     string `json:"content"`
	Origin      string `json:"origin"`
	Author      string `json:"author"`
	Image       string `json:"image"`
	PublishDate string `json:"publish_date"`
}

const (
	ARTICLE_NEW       = "NEW"
	ARTICLE_TESTED    = "TESTED"
	ARTICLE_ANALYZED  = "ANALYZED"
	ARTICLE_REVIEWED  = "REVIEWED"
	ARTICLE_DISCARDED = "DISCARDED"
	ARTICLE_LEARNING  = "ARTICLE_LEARNING"
)

const (
	QUESTION_TYPE_MULTIPLE_CHOICE = "multiple_choice"
	QUESTION_TYPE_MATCHING        = "matching_headings"
	QUESTION_TYPE_SHORT_ANSWER    = "short_answer"
	QUESTION_TYPE_TRUE_FALSE      = "true_false_not_given"
)

type ProposeWord struct {
	Word         string  `json:"word"`
	Context      string  `json:"context_sentence"`
	WordFreq     float32 `json:"freq"`
	RefArticleID uint64  `json:"ref_id"`
}

type Question struct {
	model.Base
	ArticleReadingId uint64 `json:"article_reading_id"` // ID of the article reading
	Type             string `json:"type"`
	QuestionText     string `json:"question"`

	Options []string `json:"options,omitempty"` // For multiple choice
	Answer  string   `json:"answer"`
	// Headings  []string `json:"headings,omitempty"`  // For matching headings
	// Paragraph string   `json:"paragraph,omitempty"` // For matching headings
}

type QuestionResult struct {
	QuestionID uint64 `json:"question_id"`
	Answer     string `json:"answer"`
	UserAnswer string `json:"user_answer"`
	Correct    bool   `json:"correct"`
}

type TestResult struct {
	model.Base
	ArticleReadingId uint64           `json:"article_reading_id"`
	QuestionResults  []QuestionResult `json:"question_results"`
	Score            float32          `json:"score"`
}

type ArticleReading struct {
	model.Base
	ArticleID uint64     `json:"article_id"`
	Status    string     `json:"status"`
	Questions []Question `json:"questions"`
}

type LearningWord struct {
	Word       string   `json:"word"`
	Position   string   `json:"position"`
	Definition string   `json:"definition"`
	Examples   []string `json:"examples"`

	WordFreq float32 `json:"freq"`
}

type IeExample struct {
	Text string `json:"text"`
}

type IeDefinition struct {
	Text     string      `json:"text"`
	Position string      `json:"pos"`
	Examples []IeExample `json:"example"`
}

type IePronunciation struct {
	Lang string `json:"lang"`
	Pron string `json:"pron"`
	Url  string `json:"url"`
}

type IeVocab struct {
	model.Base
	Word          string            `json:"word"`
	VocabListId   uint64            `json:"vocab_list_id"`
	Pronunciation []IePronunciation `json:"pronunciation"`
	Definitions   []IeDefinition    `json:"definition"`
	Properties    map[string]string `json:"properties"`
}

type IeVocabList struct {
	model.Base
	Name         string    `json:"name"`
	Vocabs       []IeVocab `json:"vocabs"`
	RefArticleID uint64    `json:"ref_article_id"`
}
