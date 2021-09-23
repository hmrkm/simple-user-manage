package domain

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCreateHash(t *testing.T) {
	testCases := []struct {
		name     string
		src      string
		expected string
	}{
		{
			"正常ケース",
			"password",
			"5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			actual := CreateHash(tc.src)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("CreateHash() value is missmatch :%s", diff)
			}
		})
	}
}
func TestCreateULID(t *testing.T) {
	layout := "2006年01月02日 15時04分05秒 (MST)"
	now, _ := time.Parse(layout, "2021年09月08日 10時00分00秒 (JST)")

	testCases := []struct {
		name     string
		now      time.Time
		expected string
	}{
		{
			"正常ケース",
			now,
			"01FF1EPDM0V1J25N30TAC2M9WM",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			actual := CreateULID(tc.now)

			if diff := cmp.Diff(tc.expected, actual); diff != "" {
				t.Errorf("CreateULID() value is missmatch :%s", diff)
			}
		})
	}
}

func TestPager(t *testing.T) {
	testCases := []struct {
		name         string
		count        int
		page         int
		limit        int
		expectedNow  int
		expectedLast int
	}{
		{
			"通常ケース",
			100,
			1,
			10,
			1,
			10,
		},
		{
			"通常ケース2",
			100,
			3,
			10,
			3,
			10,
		},
		{
			"通常ケース3",
			20,
			7,
			3,
			7,
			7,
		},
		{
			"countが小さいケース",
			100,
			11,
			10,
			10,
			10,
		},
		{
			"limitが0のケース",
			100,
			1,
			0,
			1,
			1,
		},
		{
			"limitが0未満のケース",
			100,
			1,
			-1,
			1,
			1,
		},
		{
			"pageが0のケース",
			100,
			0,
			10,
			1,
			10,
		},
		{
			"pageが0未満のケース",
			100,
			-1,
			10,
			1,
			10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			actualNow, actualLast := Pager(tc.count, tc.page, tc.limit)

			if diff := cmp.Diff(tc.expectedNow, actualNow); diff != "" {
				t.Errorf("Pager() now is missmatch :%s", diff)
			}
			if diff := cmp.Diff(tc.expectedLast, actualLast); diff != "" {
				t.Errorf("Pager() last is missmatch :%s", diff)
			}
		})
	}
}
