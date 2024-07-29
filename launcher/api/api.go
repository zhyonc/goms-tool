package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	serverAddr                    string = "http://127.0.0.1"
	PathJWT                       string = "/jwt"
	PathAccount                   string = "/account"
	PathAccountPassword           string = PathAccount + "/password"
	PathAuthGroup                 string = "/auth"
	PathAuthAccount               string = PathAuthGroup + "/account"
	PathAuthAccountPassword       string = PathAuthAccount + "/password"
	PathAuthAccountSecondPassword string = PathAuthAccount + "/second-password"
	PathAuthGame                  string = PathAuthGroup + "/game"
	PathAuthGameSkipSDOAuth       string = PathAuthGame + "/skip-sdo-auth"
	PathAuthGameKick              string = PathAuthGame + "/kick"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Timeout: 3 * time.Second,
	}
}

func SendRequest(method string, path string, domain any, accessToken string) ([]byte, error) {
	data, err := json.Marshal(domain)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, serverAddr+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", string(bytes))
	}
	return bytes, nil
}
