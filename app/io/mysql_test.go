package io

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/hmrkm/simple-user-manage/domain"
	"gorm.io/gorm"
)

func TestCreateDSN(t *testing.T) {
	testCases := []struct {
		name     string
		user     string
		password string
		database string
		expected string
	}{
		{
			"正常ケース",
			"user",
			"passwd",
			"db",
			"user:passwd@tcp(mysql:3306)/db?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CreateDSN(tc.user, tc.password, tc.database)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("CreateDSN() value is missmatch :%s", diff)
			}
		})
	}
}

func TestClose(t *testing.T) {
	mysql, _ := NewMysqlMock()
	testCases := []struct {
		name     string
		msql     Mysql
		err      error
		expected error
	}{
		{
			"正常ケース",
			mysql,
			nil,
			nil,
		},
		{
			"異常ケース",
			mysql,
			gorm.ErrInvalidDB,
			gorm.ErrInvalidDB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mysql := tc.msql
			if tc.err != nil {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mdbc := NewMockGormConn(ctrl)
				mdbc.EXPECT().DB().Return(&sql.DB{}, tc.err)
				mysql.conn = mdbc
			}

			actual := mysql.Close()

			if !errors.Is(actual, tc.expected) {
				t.Errorf("Close() actualErr: %v, ecpectedErr: %v", actual, tc.expected)
			}
		})
	}
}

func TestFind(t *testing.T) {
	testCases := []struct {
		name        string
		conds       string
		params      interface{}
		dbId        string
		dbName      string
		expected    MockTable
		expectedErr error
	}{
		{

			"正常ケース",
			"id=?",
			"1",
			"1",
			"aaa",
			MockTable{
				Id:   "1",
				Name: "aaa",
			},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mysql, sqlMock := NewMysqlMock()

			actual := MockTable{}
			sqlMock.ExpectQuery(regexp.QuoteMeta(
				"SELECT * FROM `mock_tables` WHERE id=(?)",
			)).
				WithArgs(tc.params).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
					AddRow(tc.dbId, tc.dbName))

			actualErr := mysql.Find(&actual, tc.conds, tc.params)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("Find() value is missmatch :%s", diff)
			}
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Find() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	testCases := []struct {
		name        string
		conds       string
		params      interface{}
		dbId        string
		dbName      string
		dbErr       error
		expected    MockTable
		expectedErr error
	}{
		{

			"正常ケース",
			"id=?",
			"1",
			"1",
			"aaa",
			nil,
			MockTable{
				Id:   "1",
				Name: "aaa",
			},
			nil,
		},
		{

			"レコードが見つからない異常ケース",
			"id=?",
			"1",
			"1",
			"aaa",
			gorm.ErrRecordNotFound,
			MockTable{},
			domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mysql, sqlMock := NewMysqlMock()

			actual := MockTable{}
			if tc.dbErr == nil {
				sqlMock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `mock_tables` WHERE id=(?)",
				)).
					WithArgs(tc.params).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(tc.dbId, tc.dbName))
			} else {
				sqlMock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `mock_tables` WHERE id=(?)",
				)).
					WithArgs(tc.params).
					WillReturnError(tc.dbErr)
			}

			actualErr := mysql.First(&actual, tc.conds, tc.params)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("First() value is missmatch :%s", diff)
			}
			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("First() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestFindWithLimitAndOffset(t *testing.T) {
	testCases := []struct {
		name        string
		mockTable   MockTable
		limit       int
		offset      int
		dbId        string
		dbName      string
		expectedErr error
	}{
		{

			"正常ケース",
			MockTable{
				Id:   "1",
				Name: "aaa",
			},
			10,
			10,
			"aaa",
			"name",
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mysql, sqlMock := NewMysqlMock()

			actual := MockTable{}
			sqlMock.ExpectQuery(regexp.QuoteMeta(
				"SELECT * FROM `mock_tables` LIMIT " + fmt.Sprint(tc.limit) + " OFFSET " + fmt.Sprint(tc.offset),
			)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
					AddRow(tc.dbId, tc.dbName))

			actualErr := mysql.FindWithLimitAndOffset(&actual, tc.limit, tc.offset)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("FindWithLimitAndOffset() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	testCases := []struct {
		name        string
		mockTable   MockTable
		expectedErr error
	}{
		{

			"正常ケース",
			MockTable{
				Id:   "1",
				Name: "aaa",
			},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mysql, sqlMock := NewMysqlMock()

			sqlMock.ExpectBegin()
			sqlMock.ExpectExec(regexp.QuoteMeta(
				"INSERT INTO `mock_tables` (`id`,`name`) VALUES (?,?)",
			)).
				WithArgs(tc.mockTable.Id, tc.mockTable.Name).
				WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectCommit()
			sqlMock.ExpectClose()

			actualErr := mysql.Create(tc.mockTable)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Create() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name        string
		mockTable   MockTable
		params      map[string]interface{}
		expectedErr error
	}{
		{

			"正常ケース",
			MockTable{
				Id:   "1",
				Name: "aaa",
			},
			map[string]interface{}{
				"name": "bbb",
			},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mysql, sqlMock := NewMysqlMock()

			sqlMock.ExpectBegin()
			sqlMock.ExpectExec(regexp.QuoteMeta(
				"UPDATE `mock_tables` SET `name`=? WHERE `id` = ?",
			)).
				WithArgs(tc.params["name"], tc.mockTable.Id).
				WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectCommit()
			sqlMock.ExpectClose()

			actualErr := mysql.Update(&tc.mockTable, tc.params)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Update() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name        string
		mockTable   MockTable
		expectedErr error
	}{
		{

			"正常ケース",
			MockTable{
				Id:   "1",
				Name: "aaa",
			},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mysql, sqlMock := NewMysqlMock()

			sqlMock.ExpectBegin()
			sqlMock.ExpectExec(regexp.QuoteMeta(
				"DELETE FROM `mock_tables` WHERE `mock_tables`.`id` = ?",
			)).
				WithArgs(tc.mockTable.Id).
				WillReturnResult(sqlmock.NewResult(1, 1))
			sqlMock.ExpectCommit()
			sqlMock.ExpectClose()

			actualErr := mysql.Delete(&tc.mockTable)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Delete() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestCount(t *testing.T) {
	testCases := []struct {
		name        string
		mockTable   MockTable
		dbCount     int64
		expectedErr error
	}{
		{

			"正常ケース",
			MockTable{
				Id:   "1",
				Name: "aaa",
			},
			10,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mysql, sqlMock := NewMysqlMock()

			sqlMock.ExpectQuery(regexp.QuoteMeta(
				"SELECT count(*) FROM `mock_tables`",
			)).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).
					AddRow(tc.dbCount))

			count := int64(0)
			actualErr := mysql.Count(&tc.mockTable, &count)

			if !errors.Is(actualErr, tc.expectedErr) {
				t.Errorf("Count() actualErr: %v, ecpectedErr: %v", actualErr, tc.expectedErr)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{

			"正常ケース",
			gorm.ErrRecordNotFound,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mysql, _ := NewMysqlMock()

			actual := mysql.IsNotFoundError(tc.err)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("IsNotFoundError() value is missmatch :%s", diff)
			}
		})
	}
}
