package usecase

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/hmrkm/simple-user-manage/domain"
)

func TestRights(t *testing.T) {
	testCases := []struct {
		name        string
		endporint   string
		ctx         context.Context
		userID      string
		resource    string
		expectedErr error
	}{
		{
			"正常ケース",
			"dummy_endpoint",
			context.Background(),
			"dummy_user_id",
			"dummy_resorce",
			nil,
		},
		{
			"異常ケース",
			"dummy_endpoint",
			context.Background(),
			"dummy_user_id",
			"dummy_resorce",
			domain.ErrRequest,
		},	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cm := domain.NewMockCommunicator(ctrl)
			cm.EXPECT().Request(tc.ctx, tc.endporint, map[string]string{
				"user_id":  tc.userID,
				"resource": tc.resource,
			}).Return(nil, tc.expectedErr)

			ru := NewRights(tc.endporint, cm)

			actualErr := ru.Verify(tc.ctx, tc.userID, tc.resource)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Update() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}
