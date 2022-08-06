package snapmaker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type connectResponse struct {
	Token        string `json:"token"`
	Readonly     bool   `json:"readonly"`
	Series       string `json:"series"`
	HeadType     int    `json:"headType"`
	HasEnclosure bool   `json:"hasEnclosure"`
}

func (sm *Snapmaker) Connect() error {
	req, err := sm.buildConnectRequest()
	if err != nil {
		return err
	}

	resp, err := sm.httpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connect request failed: %s\n%s", resp.Status, string(body))
	}

	sm.token, err = getToken(body)
	if err != nil {
		return err
	}

	// status loop is needed to avoid connection loss
	go sm.statusLoop()
	go sm.statusLoopBeat()
	return nil
}

func (sm *Snapmaker) WaitForConnection(timeout time.Duration) error {
	timer := time.NewTimer(timeout)

	sm.createStatusSyncChannel()
	defer sm.closeStatusSyncChannel()

	select {
	case <-timer.C:
		return fmt.Errorf("Waiting for connection timeout.")

	case <-sm.waitForNewStatus:
		return nil
	}
}

func (sm *Snapmaker) Close() {
	sm.cancel()
}

func (sm *Snapmaker) buildConnectRequest() (*http.Request, error) {
	var requestBody io.Reader
	if len(sm.token) > 0 {
		requestBody = strings.NewReader(url.Values{"token": {sm.token}}.Encode())
	}

	req, err := http.NewRequest("POST", sm.buildApiUrl(apiConnect), requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	return req, err
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body[:n], nil
}

func getToken(body []byte) (string, error) {
	var jsonResp connectResponse
	err := json.Unmarshal(body, &jsonResp)
	if err != nil {
		return "", err
	}
	return jsonResp.Token, nil
}
