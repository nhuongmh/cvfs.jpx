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

type EWord struct {
	model.Base
	Word         string `json:"word"`
	RefArticleID uint64 `json:"ref_id"`
}
