package jpxservice

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type SentenceFormula struct {
	Form        string `json:"form"`
	Description string `json:"description"`
	Backward    string `json:"backward"`
}

type MinnaLesson struct {
	Minna    int               `json:"minna"`
	Formulas []SentenceFormula `json:"formulas"`
}

func ParseMinnaLessonCfg(cfgFile string) (*[]MinnaLesson, error) {
	yamlFile, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed read cfg file")
	}

	var minaLessons []MinnaLesson
	err = yaml.Unmarshal(yamlFile, &minaLessons)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing cfg file")
	}

	return &minaLessons, nil
}
