package mazii

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMaziiFetcher_FetchMaziiWord(t *testing.T) {
	type args struct {
		kj string
	}
	fetcher := NewMaziiFetcher()
	tests := []struct {
		name    string
		m       *MaziiFetcher
		args    args
		want    *MaziiSearchResult
		wantErr bool
	}{
		{
			name: "Valid Word",
			m:    fetcher,
			args: args{
				kj: "ふるさと",
			},
			want: &MaziiSearchResult{
				Results: []MaziiSearchResultEntry{
					{
						Phonetic:  "ふるさと",
						ShortMean: "quê hương; nơi chôn nhau cắt rốn",
						Word:      "古里",
					},
				},
				Found: true,
			},
			wantErr: false,
		},
		{
			name: "Invalid Kanji",
			m:    fetcher,
			args: args{
				kj: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.SearchMaziiWord(tt.args.kj)
			if (err != nil) != tt.wantErr {
				if got != nil {
					fmt.Println(*got)
				}
				t.Errorf("MaziiFetcher.SearchMaziiWord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// compare kanji and meaning of the first result only
			if got != nil && tt.want != nil {
				if len(got.Results) > 0 && len(tt.want.Results) > 0 {
					jsonGot, _ := json.Marshal(got)
					fmt.Printf("Got: %v", string(jsonGot))
					if got.Results[0].Phonetic != tt.want.Results[0].Phonetic {
						t.Errorf("MaziiFetcher.SearchMaziiWord() = %v, want %v", got.Results[0].Phonetic, tt.want.Results[0].Phonetic)
					}
					if got.Results[0].Word != tt.want.Results[0].Word {
						t.Errorf("MaziiFetcher.SearchMaziiWord() = %v, want %v", got.Results[0].Word, tt.want.Results[0].Word)
					}
				}
			}
		})
	}
}
