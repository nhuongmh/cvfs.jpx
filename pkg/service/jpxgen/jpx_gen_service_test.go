package jpxgen

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
)

type mockRepo struct {
}

func (m *mockRepo) AddCard(ctx context.Context, card *langfi.ReviewCard) error {
	return nil
}
func (m *mockRepo) GetCard(ctx context.Context, cardID uint64) (*langfi.ReviewCard, error) {
	return &langfi.ReviewCard{}, nil
}
func (m *mockRepo) UpdateCard(ctx context.Context, card *langfi.ReviewCard) error {
	return nil
}
func (m *mockRepo) FetchReviewCard(ctx context.Context, groupID string) (*langfi.ReviewCard, error) {
	return &langfi.ReviewCard{}, nil
}
func (m *mockRepo) GetCardByFront(ctx context.Context, front string) (*[]langfi.ReviewCard, error) {
	return &[]langfi.ReviewCard{}, nil
}
func (m *mockRepo) FetchUnProcessCard(ctx context.Context, groupID string) (*langfi.ReviewCard, error) {
	return &langfi.ReviewCard{}, nil
}

func (m *mockRepo) DeleteNewCard(ctx context.Context) error {
	return nil
}

func (m *mockRepo) GetGroupStats(ctx context.Context) (*[]langfi.GroupSummaryDto, error) {
	return &[]langfi.GroupSummaryDto{}, nil
}

func newMockPracticeRepo() {

}

func Test_jpxService_BuildCards(t *testing.T) {
	repo := mockRepo{}
	jps := NewJpxService(&repo, time.Second, bootstrap.NewEnv())
	tests := []struct {
		name    string
		jps     jp.JpxGeneratorService
		ctx     context.Context
		want    *[]langfi.ReviewCard
		wantErr bool
	}{
		{
			name:    "test 1",
			jps:     jps,
			ctx:     context.Background(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.jps.BuildCards(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("jpxService.BuildCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jpxService.BuildCards() = %v, want %v", got, tt.want)
			}
		})
	}
}
