package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func Test_gemini(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		logger.InitLog()
		env := bootstrap.NewEnv()
		gemini, err := NewGoogleAI(env.GoogleAIKey)

		if err != nil {
			t.Errorf("Failed to init gemini: %v", err)
			return
		}
		expectedType := &genai.Schema{
			Type:  genai.TypeArray,
			Items: &genai.Schema{Type: genai.TypeString},
		}
		resp, err := gemini.GenerateContent(context.Background(),
			expectedType,
			"List 10 top songs of Greenday")
		if err != nil {
			t.Errorf("Failed to generate content: %v", err)
			return
		}
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				var recipes []string
				if err := json.Unmarshal([]byte(txt), &recipes); err != nil {
					logger.Log.Error().Err(err).Msg("Failed to unmarshal AI generated recipes")
				}
				for _, recipe := range recipes {
					fmt.Println(recipe)
				}
			}
		}
	})
}
