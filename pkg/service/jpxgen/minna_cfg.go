package jpxgen

const (
	RANDOM_GENERATOR_STRATEGY   = "random"
	ITERATOR_GENERATOR_STRATEGY = "onebyone"
)

// func ParseMinnaLessonCfg(cfgFile string) (*[]jp.MinnaLesson, error) {
// 	yamlFile, err := os.ReadFile(cfgFile)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed read cfg file")
// 	}

// 	var minaLessons []jp.MinnaLesson
// 	err = yaml.Unmarshal(yamlFile, &minaLessons)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed parsing cfg file")
// 	}

// 	return &minaLessons, nil
// }
