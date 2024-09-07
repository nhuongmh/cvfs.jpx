package jp

import (
	"context"

	"github.com/nhuongmh/cfvs.jpx/pkg/model"
	"github.com/nhuongmh/cfvs.jpx/pkg/tsunami"
)

const (
	KANA     = "kana"
	MEANING  = "meaning"
	HAN_VIE  = "han_viet"
	CATEGORY = "category"
)

type Word struct {
	model.Base
	tsunami.Entry
	Category string
}

func NewWord(w string) Word {
	return Word{
		Entry: tsunami.Entry{
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

type JpxService interface {
	InitData(ctx context.Context) error
	SyncWordList(ctx context.Context) error
}

type JpxRepository interface {
}
