package bootstrap

import (
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
)

type Env struct {
	AppMode                string `mapstructure:"APP_MODE"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	SqliteDBUrl            string `mapstructure:"SQLITE_DB_URL"`
	GoogleKeyBase64        string `mapstructure:"GOOGLE_API_KEY_BASE64"`
	GoogleSpreadSheetId    string `mapstructure:"GOOGLE_SPREADSHEET_ID"`
	GoogleWordSheetName    string `mapstructure:"GOOGLE_WORD_SHEET_NAME"`
	GoogleFormulaSheetName string `mapstructure:"GOOGLE_FORMULA_SHEET_NAME"`
}

func NewEnv() *Env {
	logger.Log.Info().Msg("Reading ENV")
	envData := Env{}
	var k = koanf.New(".")
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed load .env config file")
	}
	if err := k.Load(file.Provider(".private.env"), dotenv.Parser()); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed load .private.env config file")
	}

	k.Load(env.Provider("", ".", nil), nil)
	if err := mapstructure.WeakDecode(k.All(), &envData); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed parse Env variable")
	}

	// logger.Log.Debug().Msgf("ENV: %v", envData)

	if strings.Contains(envData.AppMode, "dev") {
		logger.Log.Info().Msg("The application is runnning in development mode")
	}

	return &envData
}
