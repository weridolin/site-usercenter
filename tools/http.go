package tools

import (
	"io"
	"net/http"
	"time"
)

func HttpGet(url string, params map[string]string, headers map[string]string) ([]byte, int, error) {
	// os.Setenv("HTTP_PROXY", "http://127.0.0.1:10809")
	// os.Setenv("HTTPS_PROXY", "http://127.0.0.1:10809")
	client := &http.Client{Timeout: time.Second * 60 * 180}
	if params != nil {
		url = url + "?" + joinQueryParams(params)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, err
	}
	return body, resp.StatusCode, nil
}

func joinQueryParams(params map[string]string) string {
	var query string
	for k, v := range params {
		query += k + "=" + v + "&"
	}
	return query[:len(query)-1] //去掉结尾的&
}

func HttpPost(url string, headers map[string]string, body io.Reader) ([]byte, int, error) {
	// os.Setenv("HTTP_PROXY", "http://127.0.0.1:10809")
	// os.Setenv("HTTPS_PROXY", "http://127.0.0.1:10809")
	client := &http.Client{Timeout: time.Second * 60 * 180}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, err
	}
	return res, resp.StatusCode, nil
}
