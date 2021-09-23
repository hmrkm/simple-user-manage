package domain

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		email          string
		hashedPassword string
		dbErr          error
		expectedErr    error
	}{
		{
			"正常ケース",
			"aaa",
			"example@mail.com",
			"bbb",
			nil,
			nil,
		},
		{
			"異常ケース",
			"aaa",
			"example@mail.com",
			"bbb",
			ErrDBAccess,
			ErrDBAccess,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := NewMockStore(ctrl)
			sm.EXPECT().Create(gomock.Any()).Return(tc.dbErr)
			us := NewUserService(sm)

			actualErr := us.Create(tc.id, tc.email, tc.hashedPassword)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Create() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestRead(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		email          string
		hashedPassword string
		dbErr          error
		expected       User
		expectedErr    error
	}{
		{
			"正常ケース",
			"aaa",
			"example@mail.com",
			"bbb",
			nil,
			User{
				Id:             "aaa",
				Email:          "example@mail.com",
				HashedPassword: "bbb",
			},
			nil,
		},
		{
			"異常ケース",
			"aaa",
			"example@mail.com",
			"bbb",
			ErrDBAccess,
			User{},
			ErrDBAccess,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := NewMockStore(ctrl)
			sm.EXPECT().First(gomock.Any(), "id=?", tc.id).DoAndReturn(
				func(target *User, cond string, params string) error {
					if tc.dbErr == nil {
						*target = User{
							Id:             tc.id,
							Email:          tc.email,
							HashedPassword: tc.hashedPassword,
						}
					}
					return tc.dbErr
				},
			)
			us := NewUserService(sm)

			actual, actualErr := us.Read(tc.id)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Read() value is missmatch :%s", diff)
			}
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Read() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestReadList(t *testing.T) {
	testCases := []struct {
		name           string
		page           int
		limit          int
		offset         int
		id             string
		email          string
		hashedPassword string
		dbErr          error
		expected       []User
		expectedErr    error
	}{
		{
			"正常ケース",
			1,
			10,
			0,
			"aaa",
			"example@mail.com",
			"bbb",
			nil,
			[]User{
				{
					Id:             "aaa",
					Email:          "example@mail.com",
					HashedPassword: "bbb",
				},
			},
			nil,
		},
		{
			"limitマイナス正常ケース",
			1,
			-1,
			0,
			"aaa",
			"example@mail.com",
			"bbb",
			nil,
			[]User{
				{
					Id:             "aaa",
					Email:          "example@mail.com",
					HashedPassword: "bbb",
				},
			},
			nil,
		},
		{
			"pageマイナス正常ケース",
			-1,
			10,
			0,
			"aaa",
			"example@mail.com",
			"bbb",
			nil,
			[]User{
				{
					Id:             "aaa",
					Email:          "example@mail.com",
					HashedPassword: "bbb",
				},
			},
			nil,
		},
		{
			"異常ケース",
			1,
			10,
			0,
			"aaa",
			"example@mail.com",
			"bbb",
			ErrDBAccess,
			[]User{},
			ErrDBAccess,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := NewMockStore(ctrl)
			sm.EXPECT().FindWithLimitAndOffset(gomock.Any(), tc.limit, tc.offset).DoAndReturn(
				func(target *[]User, cond int, params int) error {
					if tc.dbErr == nil {
						*target = []User{
							{
								Id:             tc.id,
								Email:          tc.email,
								HashedPassword: tc.hashedPassword,
							},
						}
					}
					return tc.dbErr
				},
			)
			us := NewUserService(sm)

			actual, actualErr := us.ReadList(tc.page, tc.limit)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("ReadList() value is missmatch :%s", diff)
			}
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("ReadList() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestCount(t *testing.T) {
	testCases := []struct {
		name        string
		dbCount     int64
		dbErr       error
		expected    int
		expectedErr error
	}{
		{
			"正常ケース",
			10,
			nil,
			10,
			nil,
		},
		{
			"異常ケース",
			10,
			ErrDBAccess,
			0,
			ErrDBAccess,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := NewMockStore(ctrl)
			sm.EXPECT().Count(gomock.Any(), gomock.Any()).DoAndReturn(
				func(target *User, count *int64) error {
					if tc.dbErr == nil {
						*count = tc.dbCount
					}
					return tc.dbErr
				},
			)
			us := NewUserService(sm)

			actual, actualErr := us.Count()

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Count() value is missmatch :%s", diff)
			}
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Count() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		confirm  string
		expected bool
	}{
		{
			"不一致ケース",
			"aaa",
			"bbb",
			false,
		},
		{
			"一致ケース",
			"aaa",
			"aaa",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			sm := NewMockStore(ctrl)
			us := NewUserService(sm)

			actual := us.VerifyPassword(tc.input, tc.confirm)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("VerifyPassword() value is missmatch :%s", diff)
			}
		})
	}
}
