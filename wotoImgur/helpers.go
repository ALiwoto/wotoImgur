package wotoImgur

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func createUploadForm(image []byte, album, dType, title, description string) url.Values {
	form := url.Values{}

	form.Add("image", string(image[:]))
	form.Add("type", dType)

	if album != "" {
		form.Add("album", album)
	}
	if title != "" {
		form.Add("title", title)
	}
	if description != "" {
		form.Add("description", description)
	}

	return form
}

func extractRateLimits(h http.Header) (*RateLimit, error) {
	rl := new(RateLimit)
	var err error

	userLimitStr := h.Get("X-RateLimit-UserLimit")
	if userLimitStr != "" {
		rl.UserLimit, err = strconv.ParseInt(userLimitStr, 10, 32)
	}

	userRemainingStr := h.Get("X-RateLimit-UserRemaining")
	if userRemainingStr != "" {
		rl.UserRemaining, err = strconv.ParseInt(userRemainingStr, 10, 32)
	}

	unixTimeStr := h.Get("X-RateLimit-UserReset")
	if unixTimeStr != "" {
		var userReset int64
		userReset, err = strconv.ParseInt(unixTimeStr, 10, 64)
		rl.UserReset = time.Unix(userReset, 0)
	}

	clientLimitStr := h.Get("X-RateLimit-ClientLimit")
	if clientLimitStr != "" {
		rl.ClientLimit, err = strconv.ParseInt(clientLimitStr, 10, 32)
	}

	clientRemainingStr := h.Get("X-RateLimit-ClientRemaining")
	if clientRemainingStr != "" {
		rl.ClientRemaining, err = strconv.ParseInt(clientRemainingStr, 10, 32)
	}

	return rl, err
}

func getErr(status int, value string) *ImgurError {
	return &ImgurError{
		Status:  status,
		Message: value,
	}
}

func getErrF(status int, format string, args ...any) *ImgurError {
	return &ImgurError{
		Err:    fmt.Errorf(format, args...),
		Status: status,
	}
}
