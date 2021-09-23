package io

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/jarcoal/httpmock"
)

func TestRequest(t *testing.T) {
	url := "http://auth/v1/verify"
	testCases := []struct {
		name          string
		mockStatus    int
		mockResBody   string
		m             map[string]string
		expectResBody []byte
		expectErr     error
	}{
		{
			"HTTPリクエスト正常ケース",
			200,
			"ok",
			map[string]string{
				"token": "aaa",
			},
			[]byte("ok"),
			nil,
		},
		{
			"HTTPリクエスト400異常ケース",
			400,
			"ng",
			map[string]string{
				"token": "aaa",
			},
			nil,
			errors.New("url is " + url),
		},
		{
			"HTTPリクエスト500異常ケース",
			500,
			"ng",
			map[string]string{
				"token": "aaa",
			},
			nil,
			domain.ErrNotReaching,
		},
		{
			"HTTPリクエスト100異常ケース",
			100,
			"ng",
			map[string]string{
				"token": "aaa",
			},
			nil,
			errors.New("Continue"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("POST", url, httpmock.NewStringResponder(tc.mockStatus, tc.mockResBody))

			hc := NewHTTP(1, 1)

			ctx := context.Background()

			actualResBody, actualErr := hc.Request(ctx, url, tc.m)

			if diff := cmp.Diff(actualResBody, tc.expectResBody); diff != "" {
				t.Errorf("Request() response is missmatch %s", diff)
			}
			if tc.expectErr != nil {
				if !errors.As(actualErr, &tc.expectErr) {
					t.Errorf("Request() error = %v, expectErr %v", actualErr, tc.expectErr)
				}
			} else {
				if !errors.Is(actualErr, tc.expectErr) {
					t.Errorf("Request() error = %v, expectErr %v", actualErr, tc.expectErr)
				}
			}
		})
	}
}
