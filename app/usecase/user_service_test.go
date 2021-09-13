package usecase

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
				t.Errorf("Create() value is missmatch :%s", diff)
			}
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Verify() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
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
			if tc.dbErr == nil {
				sm.EXPECT().Update(gomock.Any(), map[string]interface{}{
					"email": tc.email,
				}).Return(tc.dbErr)
			}
			us := NewUserService(sm)

			actualErr := us.Update(tc.id, tc.email)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Verify() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestUpdateWithPassword(t *testing.T) {
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
			if tc.dbErr == nil {
				sm.EXPECT().Update(gomock.Any(), map[string]interface{}{
					"email":          tc.email,
					"hashedPassword": tc.hashedPassword,
				}).Return(tc.dbErr)
			}
			us := NewUserService(sm)

			actualErr := us.UpdateWithPassword(tc.id, tc.email, tc.hashedPassword)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Verify() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
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
			if tc.dbErr == nil {
				sm.EXPECT().Delete(gomock.Any()).Return(tc.dbErr)
			}
			us := NewUserService(sm)

			actualErr := us.Delete(tc.id)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Verify() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}
