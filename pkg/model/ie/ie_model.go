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

type ProposeWord struct {
	Word         string  `json:"word"`
	Context      string  `json:"context_sentence"`
	WordFreq     float32 `json:"freq"`
	RefArticleID uint64  `json:"ref_id"`
}

type WordList struct {
	model.Base
}

type Question struct {
	model.Base
	ArticleReadingId uint64 `json:"article_reading_id"` // ID of the article reading
	Type             string `json:"type"`
	QuestionText     string `json:"question"`

	Options   []string `json:"options,omitempty"` // For multiple choice
	Answer    string   `json:"answer"`
	Headings  []string `json:"headings,omitempty"`  // For matching headings
	Paragraph string   `json:"paragraph,omitempty"` // For matching headings
}

type ArticleReading struct {
	model.Base
	ArticleID uint64     `json:"article_id"`
	Status    string     `json:"status"`
	Score     float32    `json:"score"`
	Questions []Question `json:"questions"`
}
