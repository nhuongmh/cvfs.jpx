package mazii

import (
	"fmt"
	"testing"
)

func TestMaziiFetcher_FetchMaziiKanji(t *testing.T) {
	type args struct {
		kj string
	}
	fetcher := NewMaziiFetcher()
	tests := []struct {
		name    string
		m       *MaziiFetcher
		args    args
		want    *KanjiResultResp
		wantErr bool
	}{
		{
			name: "Valid Kanji",
			m:    fetcher,
			args: args{
				kj: "君",
			},
			want: &KanjiResultResp{
				Results: []KanjiResultEntry{
					{
						Kanji:   "君",
						Meaning: "QUÂN",
					},
				},
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
			got, err := tt.m.FetchMaziiKanji(tt.args.kj)
			if (err != nil) != tt.wantErr {
				if got != nil {
					fmt.Println(*got)
				}
				t.Errorf("MaziiFetcher.FetchMaziiKanji() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// compare kanji and meaning of the first result only
			if got != nil && tt.want != nil {
				if len(got.Results) > 0 && len(tt.want.Results) > 0 {
					if got.Results[0].Kanji != tt.want.Results[0].Kanji {
						t.Errorf("MaziiFetcher.FetchMaziiKanji() = %v, want %v", got.Results[0].Kanji, tt.want.Results[0].Kanji)
					}
					if got.Results[0].Meaning != tt.want.Results[0].Meaning {
						t.Errorf("MaziiFetcher.FetchMaziiKanji() = %v, want %v", got.Results[0].Meaning, tt.want.Results[0].Meaning)
					}
				}
			}
		})
	}
}

func TestMaziiFetcher_FetchBestComment(t *testing.T) {
	type args struct {
		wordID int
	}
	fetcher := NewMaziiFetcher()
	tests := []struct {
		name    string
		m       *MaziiFetcher
		args    args
		wantErr bool
	}{
		{
			name: "Valid WordID",
			m:    fetcher,
			args: args{
				wordID: 673,
			},
			wantErr: false,
		},
		{
			name: "Invalid WordID",
			m:    fetcher,
			args: args{
				wordID: -1,
			},
			wantErr: true,
		},
		{
			name: "No Comments Found",
			m:    fetcher,
			args: args{
				wordID: 67890,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.FetchBestComment(tt.args.wordID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaziiFetcher.FetchBestComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if len(got.Text) == 0 {
					t.Errorf("MaziiFetcher.FetchBestComment() = %v, want not empty", got.Text)
				} else {
					fmt.Printf("Got: %v", got.Text)
				}
			}
		})
	}
}
