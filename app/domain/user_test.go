package domain

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestUserTableName(t *testing.T) {
	tableName := User{}.TableName()

	if tableName != "users" {
		t.Errorf("TableName is not tokens")
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name        string
		user        User
		value       User
		dbErr       error
		expectedErr error
	}{
		{
			"正常ケース",
			User{
				Email: "aaa@aaa.com",
			},
			User{
				Email: "aaa@aaa.com",
			},
			nil,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := NewMockStore(ctrl)
			sm.EXPECT().Update(gomock.Any(), map[string]interface{}{
				"email": tc.user.Email,
			}).Return(tc.dbErr)

			actualErr := tc.user.Update(sm, tc.value)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Update() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestUpdateWithPassword(t *testing.T) {
	testCases := []struct {
		name        string
		user        User
		value       User
		dbErr       error
		expectedErr error
	}{
		{
			"正常ケース",
			User{
				Email:          "aaa@aaa.com",
				HashedPassword: "aaa",
			},
			User{
				Email:          "aaa@aaa.com",
				HashedPassword: "aaa",
			},
			nil,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := NewMockStore(ctrl)
			sm.EXPECT().Update(gomock.Any(), map[string]interface{}{
				"email":    tc.user.Email,
				"password": tc.user.HashedPassword,
			}).Return(tc.dbErr)

			actualErr := tc.user.UpdateWithPassword(sm, tc.value)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Update() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name        string
		user        User
		value       User
		dbErr       error
		expectedErr error
	}{
		{
			"正常ケース",
			User{},
			User{},
			nil,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := NewMockStore(ctrl)
			sm.EXPECT().Delete(gomock.Any()).Return(tc.dbErr)

			actualErr := tc.user.Delete(sm)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Update() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}
