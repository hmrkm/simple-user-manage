package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/hmrkm/simple-user-manage/domain"
)

func TestAuth(t *testing.T) {
	authResponse := AuthResponse{
		User: struct {
			Id    string `json:"id"`
			Email string `json:"email"`
		}{
			Id:    "aaa",
			Email: "aaa@example.com",
		},
	}
	authResponseStr, _ := json.Marshal(authResponse)
	testCases := []struct {
		name         string
		ctx          context.Context
		token        string
		authResponse []byte
		requestErr   error
		requireRead  bool
		dbUser       domain.User
		dbErr        error
		expected     domain.User
		expectedErr  error
	}{
		{
			"正常ケース",
			context.Background(),
			"token",
			[]byte(authResponseStr),
			nil,
			true,
			domain.User{
				Id: "aaa",
			},
			nil,
			domain.User{
				Id: "aaa",
			},
			nil,
		},
		{
			"リクエストエラーの異常ケース",
			context.Background(),
			"token",
			[]byte(authResponseStr),
			domain.ErrNotReaching,
			false,
			domain.User{},
			nil,
			domain.User{},
			domain.ErrNotReaching,
		},
		{
			"リクエストレスポンスの異常ケース",
			context.Background(),
			"token",
			[]byte("aaaa"),
			nil,
			false,
			domain.User{},
			nil,
			domain.User{},
			domain.ErrInvalidValue,
		},
		{
			"DBエラーの異常ケース",
			context.Background(),
			"token",
			[]byte(authResponseStr),
			nil,
			true,
			domain.User{},
			domain.ErrDBAccess,
			domain.User{},
			domain.ErrDBAccess,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cm := domain.NewMockCommunicator(ctrl)
			cm.EXPECT().Request(tc.ctx, "", map[string]string{
				"token": tc.token,
			}).Return(tc.authResponse, tc.requestErr)
			usdm := domain.NewMockUserService(ctrl)
			if tc.requireRead {
				usdm.EXPECT().Read("aaa").Return(tc.dbUser, tc.dbErr)
			}

			auth := NewAuth("", cm, usdm)

			actual, actualErr := auth.Auth(tc.ctx, tc.token)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Auth() user is missmatch :%s", diff)
			}

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Update() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}
