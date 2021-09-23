package usecase

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/hmrkm/simple-user-manage/domain"
)

func TestList(t *testing.T) {
	testCases := []struct {
		name          string
		page          int
		limit         int
		dbUsers       []domain.User
		dbErr         error
		dbCount       int
		dbCountErr    error
		expectedUsers []domain.User
		expectedNow   int
		expectedLast  int
		expectedErr   error
	}{
		{
			"正常ケース",
			1,
			10,
			[]domain.User{
				{
					Id: "aaa",
				},
			},
			nil,
			30,
			nil,
			[]domain.User{
				{
					Id: "aaa",
				},
			},
			1,
			3,
			nil,
		},
		{
			"ReadListDBエラー異常ケース",
			1,
			10,
			[]domain.User{
				{
					Id: "aaa",
				},
			},
			domain.ErrDBAccess,
			30,
			nil,
			nil,
			0,
			0,
			domain.ErrDBAccess,
		},
		{
			"CountDBエラー異常ケース",
			1,
			10,
			[]domain.User{
				{
					Id: "aaa",
				},
			},
			nil,
			30,
			domain.ErrDBAccess,
			nil,
			0,
			0,
			domain.ErrDBAccess,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usm := domain.NewMockUserService(ctrl)
			usm.EXPECT().ReadList(tc.page, tc.limit).Return(tc.dbUsers, tc.dbErr)
			if tc.dbErr == nil {
				usm.EXPECT().Count().Return(tc.dbCount, tc.dbCountErr)
			}
			sm := domain.NewMockStore(ctrl)

			users := NewUsers(usm, sm)

			actualUsers, actualNow, actualLast, actualErr := users.List(tc.page, tc.limit)

			if diff := cmp.Diff(tc.expectedUsers, actualUsers); diff != "" {
				t.Errorf("List() users is missmatch :%s", diff)
			}
			if diff := cmp.Diff(tc.expectedNow, actualNow); diff != "" {
				t.Errorf("List() now is missmatch :%s", diff)
			}
			if diff := cmp.Diff(tc.expectedLast, actualLast); diff != "" {
				t.Errorf("List() last is missmatch :%s", diff)
			}
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("List() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	testCases := []struct {
		name         string
		email        string
		password     string
		passwordConf string
		dbErr        error
		expectedErr  error
	}{
		{
			"正常ケース",
			"aaa@email.com",
			"password",
			"password",
			nil,
			nil,
		},
		{
			"パスワード不一致異常ケース",
			"aaa@email.com",
			"password",
			"passwd",
			nil,
			domain.ErrInvalidValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usm := domain.NewMockUserService(ctrl)
			usm.EXPECT().VerifyPassword(tc.password, tc.passwordConf).DoAndReturn(
				func(input string, conf string) bool {
					return input == conf
				},
			)
			if tc.expectedErr == nil {
				usm.EXPECT().Create(gomock.Any(), tc.email, domain.CreateHash(tc.password)).Return(tc.dbErr)
			}
			sm := domain.NewMockStore(ctrl)

			users := NewUsers(usm, sm)

			actualErr := users.Create(tc.email, tc.password, tc.passwordConf)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Create() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestDetail(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		dbUser      domain.User
		dbErr       error
		expected    domain.User
		expectedErr error
	}{
		{
			"正常ケース",
			"aaa",
			domain.User{
				Id: "aaa",
			},
			nil,
			domain.User{
				Id: "aaa",
			},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usm := domain.NewMockUserService(ctrl)
			usm.EXPECT().Read(tc.id).Return(tc.dbUser, tc.dbErr)
			sm := domain.NewMockStore(ctrl)

			users := NewUsers(usm, sm)

			actual, actualErr := users.Detail(tc.id)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Detail() user is missmatch :%s", diff)
			}
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Detail() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	password := "password"
	passwd := "passwd"
	testCases := []struct {
		name         string
		id           string
		email        string
		password     *string
		passwordConf *string
		dbUser       domain.User
		dbErr        error
		dbUpdateErr  error
		expectedErr  error
	}{
		{
			"正常ケース",
			"aaa",
			"aaa@example.com",
			&password,
			&password,
			domain.User{
				Id: "aaa",
			},
			nil,
			nil,
			nil,
		},
		{
			"パスワード変更無しの正常ケース",
			"aaa",
			"aaa@example.com",
			nil,
			nil,
			domain.User{
				Id: "aaa",
			},
			nil,
			nil,
			nil,
		},
		{
			"パスワード不一致の異常ケース",
			"aaa",
			"aaa@example.com",
			&password,
			&passwd,
			domain.User{
				Id: "aaa",
			},
			nil,
			nil,
			domain.ErrInvalidValue,
		},
		{
			"DBエラーの異常ケース",
			"aaa",
			"aaa@example.com",
			&password,
			&password,
			domain.User{
				Id: "aaa",
			},
			domain.ErrDBAccess,
			nil,
			domain.ErrDBAccess,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usm := domain.NewMockUserService(ctrl)
			sm := domain.NewMockStore(ctrl)
			if tc.password != nil {
				usm.EXPECT().VerifyPassword(*tc.password, *tc.passwordConf).DoAndReturn(
					func(input string, conf string) bool {
						return input == conf
					},
				)
			}
			if tc.expectedErr != domain.ErrInvalidValue {
				usm.EXPECT().Read(tc.id).Return(tc.dbUser, tc.dbErr)
				if tc.dbErr == nil {
					if tc.password == nil {
						sm.EXPECT().Update(gomock.Any(), map[string]interface{}{
							"email": tc.email,
						}).Return(tc.dbUpdateErr)
					} else {
						sm.EXPECT().Update(gomock.Any(), map[string]interface{}{
							"email":    tc.email,
							"password": domain.CreateHash(*tc.password),
						}).Return(tc.dbUpdateErr)
					}
				}
			}
			users := NewUsers(usm, sm)

			actualErr := users.Update(tc.id, tc.email, tc.password, tc.passwordConf)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Update() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		dbErr       error
		dbUser      domain.User
		dbDeleteErr error
		expectedErr error
	}{
		{
			"正常ケース",
			"aaa",
			nil,
			domain.User{
				Id: "aaa",
			},
			nil,
			nil,
		},
		{
			"DBエラー異常ケース",
			"aaa",
			domain.ErrDBAccess,
			domain.User{
				Id: "aaa",
			},
			nil,
			domain.ErrDBAccess,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			usm := domain.NewMockUserService(ctrl)
			sm := domain.NewMockStore(ctrl)
			usm.EXPECT().Read(tc.id).Return(tc.dbUser, tc.dbErr)
			if tc.dbErr == nil {
				sm.EXPECT().Delete(gomock.Any()).Return(tc.dbDeleteErr)
			}

			users := NewUsers(usm, sm)

			actualErr := users.Delete(tc.id)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Delete() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}
