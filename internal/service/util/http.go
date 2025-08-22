package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fabxu/datacollector-service/internal/lib/constant"
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

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)

	return string(result), e
}
