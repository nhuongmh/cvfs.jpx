package mazii

import "net/http"

type KanjiResultResp struct {
	Results []KanjiResultEntry `json:"results"`
}

type KanjiResultEntry struct {
	Kanji    string               `json:"kanji"`
	Meaning  string               `json:"mean"`
	Detail   string               `json:"detail"`
	Kun      string               `json:"kun"`
	On       string               `json:"on"`
	Examples []KanjiResultExample `json:"examples"`
	WordId   int                  `json:"mobileId"`
}

type KanjiResultExample struct {
	Meaning       string `json:"m"`
	Word          string `json:"w"`
	Hiragana      string `json:"h"`
	Pronunciation string `json:"p"`
}

type MaziiCommentEntry struct {
	Text string `json:"mean"`
}

type MaziiComments struct {
	Comments []MaziiCommentEntry `json:"result"`
}

type MaziiSearchResult struct {
	Results []MaziiSearchResultEntry `json:"data"`
	Found   bool                     `json:"found"`
}

type MaziiSearchResultEntry struct {
	Phonetic       string                 `json:"phonetic"`
	ShortMean      string                 `json:"short_mean"`
	Word           string                 `json:"word"`
	Pronunciations []MaziiWordPronunEntry `json:"pronunciation"`
	WordId         int                    `json:"mobileId"`
}

type MaziiWordMeanEntry struct {
	Kind string `json:"kind"`
	Mean string `json:"mean"`
	// Example []string `json:"example"`
}

type MaziiWordPronunEntry struct {
	Transcripts []MaziiWordPronunTranscript `json:"transcriptions"`
}

type MaziiWordPronunTranscript struct {
	Romaji string `json:"romaji"`
	Kana   string `json:"kana"`
}

type MaziiFetcher struct {
	httpClient *http.Client
}

func NewMaziiFetcher() *MaziiFetcher {
	return &MaziiFetcher{
		httpClient: &http.Client{},
	}
}
