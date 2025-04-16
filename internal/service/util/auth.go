package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func hmacSha256(secret string, message string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func GetURLPathWithParams(urlPath string, paramMap map[string]interface{}) (string, error) {
	params := url.Values{}
	parseURL, err := url.Parse(urlPath)

	if err != nil {
		return "", err
	}

	for k, v := range paramMap {
		// datautils.ToString(v)
		params.Set(k, v.(string))
	}

	parseURL.RawQuery = params.Encode()

	return parseURL.String(), nil
}

func GetAuthorization(accessKey, secretKey string) (string, string) {
	xdate := time.Now().UTC().Format(http.TimeFormat)

	strXDate := fmt.Sprintf("x-date: %s", xdate)
	strSignature := hmacSha256(secretKey, strXDate)
	sign := fmt.Sprintf(`hmac accesskey="%s", algorithm="hmac-sha256", headers="x-date", signature="%s"`, accessKey, strSignature)

	return xdate, sign
}
