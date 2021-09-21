package io

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type HTTP struct {
	retryNumber int
	sleepSecond int
}

func NewHTTP(rn int, ss int) domain.Communicator {
	return HTTP{
		retryNumber: rn,
		sleepSecond: ss,
	}
}

// 引数
// to: 宛先
// m: 宛先に送るkey/value構造のオブジェクト
// 戻り値
// []byte: リクエストが成功した際に得られる応答
func (hf HTTP) Request(
	ctx context.Context,
	to string,
	m map[string]string,
) ([]byte, error) {
	payload := url.Values{}
	for k, v := range m {
		payload.Add(k, v)
	}
	payloadString := strings.NewReader(payload.Encode())

	url, err := url.Parse(to)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), payloadString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("Content-Type", echo.MIMEApplicationJSON)

	for i := 0; i < hf.retryNumber; i++ {
		res, err := c.Do(req)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		defer res.Body.Close()

		switch {
		case res.StatusCode >= 200 && res.StatusCode < 400:
			resBytes, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			return resBytes, nil
		case res.StatusCode >= 400 && res.StatusCode < 500:
			return nil, errors.WithStack(errors.New("url is " + to))
		case res.StatusCode >= 500:
			time.Sleep(time.Duration(hf.sleepSecond) * time.Second)
		default:
			return nil, errors.WithStack(errors.New(http.StatusText(res.StatusCode)))
		}
	}

	return nil, errors.WithStack(domain.ErrNotReaching)
}
