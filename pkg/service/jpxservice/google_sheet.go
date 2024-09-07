package jpxservice

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type ggSheetDatasource struct {
	SheetSrv *sheets.Service
}

func InitNewGoogleSheetService(googleKeyBase64 string) (*ggSheetDatasource, error) {
	logger.Log.Info().Msg("Initializing google sheet service")
	// create api context
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	credBytes, err := base64.StdEncoding.DecodeString(googleKeyBase64)
	if err != nil {
		return nil, errors.Wrap(err, "Failed decode base64 google key")
	}
	// logger.Log.Debug().Msgf("API: %v", string(credBytes))

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to authenticate to google API")
	}

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.Wrap(err, "Failed creating new google service")
	}
	return &ggSheetDatasource{SheetSrv: srv}, nil
}

func (ggs *ggSheetDatasource) fetchData(spreadsheetId, sheetName string) (*[]jp.Word, error) {
	// https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>

	logger.Log.Info().Msgf("Fetching data from google sheet %v (%v)", spreadsheetId, sheetName)
	readRange := fmt.Sprintf("%s!A2:H", sheetName)
	resp, err := ggs.SheetSrv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to retrieve data from sheet")
	}

	if len(resp.Values) == 0 {
		logger.Log.Warn().Msg("No data found.")
		return nil, model.ErrNoData
	}

	wordList := make([]jp.Word, 0, 20)
	for _, row := range resp.Values {
		// fmt.Printf("%s, %s, %s\n", row[3], row[5], row[6])
		if len(row) < 7 {
			continue
		}
		word := row[3].(string)

		w := jp.NewWord(word)
		w.SetProp(jp.KANA, row[2].(string))
		w.SetProp(jp.HAN_VIE, row[4].(string))
		w.SetProp(jp.MEANING, row[5].(string))
		w.Category = row[6].(string)

		wordList = append(wordList, w)
	}

	return &wordList, nil
}
