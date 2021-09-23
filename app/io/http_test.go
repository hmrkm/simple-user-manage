package io

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/jarcoal/httpmock"
)

func TestRequest(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           context.Context
		url           string
		mockStatus    int
		mockResBody   string
		mockLocation  string
		body          interface{}
		expectResBody []byte
		expectErr     error
	}{
		{
			"HTTPリクエスト正常ケース",
			context.Background(),
			"http://auth/v1/verify",
			200,
			"ok",
			"location",
			map[string]string{
				"token": "aaa",
			},
			[]byte("ok"),
			nil,
		},
		{
			"HTTPリクエスト400異常ケース",
			context.Background(),
			"http://auth/v1/verify",
			400,
			"ng",
			"location",
			map[string]string{
				"token": "aaa",
			},
			nil,
			errors.New("url is http://auth/v1/verify"),
		},
		{
			"HTTPリクエスト500異常ケース",
			context.Background(),
			"http://auth/v1/verify",
			500,
			"ng",
			"location",
			map[string]string{
				"token": "aaa",
			},
			nil,
			domain.ErrNotReaching,
		},
		{
			"HTTPリクエスト100異常ケース",
			context.Background(),
			"http://auth/v1/verify",
			100,
			"ng",
			"location",
			map[string]string{
				"token": "aaa",
			},
			nil,
			errors.New("Continue"),
		},
		{
			"リクエストボディMarshalエラー異常ケース",
			context.Background(),
			"http://auth/v1/verify",
			500,
			"ng",
			"location",
			math.NaN(),
			nil,
			errors.New("json: unsupported value: NaN"),
		},
		{
			"urlパースエラー異常ケース",
			context.Background(),
			"%%%%%%",
			500,
			"ng",
			"location",
			111,
			nil,
			errors.New("parse \"%%%%%%\": invalid URL escape \"%%%\""),
		},
		{
			"コンテキストエラー異常ケース",
			nil,
			"http://auth/v1/verify",
			500,
			"ng",
			"location",
			111,
			nil,
			errors.New("net/http: nil Context"),
		},
		{
			"Locationヘッダーエラー異常ケース",
			context.Background(),
			"http://auth/v1/verify",
			301,
			"ok",
			"",
			111,
			nil,
			errors.New("Post \"http://auth/v1/verify\": 301 response missing Location header"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("POST", tc.url, newMockResponder(tc.mockStatus, tc.mockResBody, tc.mockLocation))

			http := NewHTTP(1, 1)

			actualResBody, actualErr := http.Request(tc.ctx, tc.url, tc.body)

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

func newMockResponder(status int, body string, location string) httpmock.Responder {
	resp := httpmock.NewStringResponse(status, body)
	resp.Header.Set("Location", location)
	return httpmock.ResponderFromResponse(resp)
}
