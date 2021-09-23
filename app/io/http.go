package io

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/hmrkm/simple-user-manage/domain"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type HTTP struct {
	retryNumber int
	sleepSecond int
}

func NewHTTP(retryNumber int, sleepSecond int) domain.Communicator {
	return HTTP{
		retryNumber: retryNumber,
		sleepSecond: sleepSecond,
	}
}

// 引数
// to: 宛先
// body: 宛先に送るオブジェクト
func (h HTTP) Request(
	ctx context.Context,
	to string,
	body interface{},
) (response []byte, err error) {
	jsn, err := json.Marshal(body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	url, err := url.Parse(to)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewBuffer(jsn))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("Content-Type", echo.MIMEApplicationJSON)

	for i := 0; i < h.retryNumber; i++ {
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
			time.Sleep(time.Duration(h.sleepSecond) * time.Second)
		default:
			return nil, errors.WithStack(errors.New(http.StatusText(res.StatusCode)))
		}
	}

	return nil, errors.WithStack(domain.ErrNotReaching)
}
