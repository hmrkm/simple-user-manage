package domain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
