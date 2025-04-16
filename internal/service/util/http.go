package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	cmerror "gitlab.senseauto.com/apcloud/library/common-go/error"
)

func GetProjectCollectorURL() string {
	dominURL := "localhost:8088"
	if v, ok := os.LookupEnv("DATACOLLECTOR_SERVICE"); ok {
		dominURL = v
		fmt.Println("show dominURL:", dominURL)
	}
	return fmt.Sprintf("http://%s", dominURL)
}

func DoRequest(method, url string, data []byte, headers map[string]string) (string, error) {
	var e error = nil

	ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultHTTPTimeout*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))

	if err != nil {
		e = cmerror.BadRequest.WithError(err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e = cmerror.BadRequest.WithError(err)
		return "", e
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)

	return string(result), e
}
