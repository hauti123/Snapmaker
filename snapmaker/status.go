package snapmaker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Status struct {
	Status                     string  `json:"status"`
	X                          float64 `json:"x"`
	Y                          float64 `json:"y"`
	Z                          float64 `json:"z"`
	Homed                      bool    `json:"homed"`
	OffsetX                    float64 `json:"offsetX"`
	OffsetY                    float64 `json:"offsetY"`
	OffsetZ                    float64 `json:"offsetZ"`
	ToolHead                   string  `json:"toolHead"`
	NozzleTemperature          int     `json:"nozzleTemperature"`
	NozzleTargetTemperature    int     `json:"nozzleTargetTemperature"`
	HeatedBedTemperature       int     `json:"heatedBedTemperature"`
	HeatedBedTargetTemperature int     `json:"heatedBedTargetTemperature"`
	IsFilamentOut              bool    `json:"isFilamentOut"`
	WorkSpeed                  int     `json:"workSpeed"`
	PrintStatus                string  `json:"printStatus"`
	ModuleList                 Modules `json:"moduleList"`
}

type Modules struct {
	Enclosure           bool `json:"enclosure"`
	RotaryModule        bool `json:"rotaryModule"`
	EmergencyStopButton bool `json:"emergencyStopButton"`
	AirPurifier         bool `json:"airPurifier"`
}

func (sm *Snapmaker) GetStatus(timeout time.Duration) (Status, error) {
	timer := time.NewTimer(timeout)

	sm.createStatusSyncChannel()
	defer sm.closeStatusSyncChannel()

	// trigger immediate status retrieval
	sm.triggerStatusRetrieval <- time.Now()

	var statusJson string
	select {
	case <-timer.C:
		return Status{}, fmt.Errorf("Status timeout.")

	case statusJson = <-sm.waitForNewStatus:
	}

	status := Status{}
	err := json.Unmarshal([]byte(statusJson), &status)

	if err != nil {
		return Status{}, err
	}

	return status, nil
}

func (sm *Snapmaker) getStatusFromPrinter() (string, error) {

	timestamp := time.Now().UnixMilli()
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s?token=%s&%d", sm.buildApiUrl(apiPrinterStatus), sm.token, timestamp),
		nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	resp, err := sm.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {

	case http.StatusOK:
		body := make([]byte, 1024)
		n, err := resp.Body.Read(body)

		if err != nil && err != io.EOF {
			return "", err
		}

		return string(body[:n]), nil

	case http.StatusNoContent:
		return "", nil

	default:
		body := make([]byte, 1024)
		n, err := resp.Body.Read(body)

		if err != nil && err != io.EOF {
			return "", err
		}
		return string(body[:n]), fmt.Errorf("Getting status failed: %s", resp.Status)
	}
}

func (sm *Snapmaker) createStatusSyncChannel() {

	sm.waitChanMutex.Lock()
	if sm.waitForNewStatus == nil {
		sm.waitForNewStatus = make(chan string)
	}
	sm.waitChanMutex.Unlock()
}

func (sm *Snapmaker) closeStatusSyncChannel() {
	sm.waitChanMutex.Lock()
	close(sm.waitForNewStatus)
	sm.waitForNewStatus = nil
	sm.waitChanMutex.Unlock()
}

// this allows for triggering status requests in between the usual "heartbeat"
func (sm *Snapmaker) statusLoopBeat() {
	ticker := time.NewTicker(sm.statusLoopInterval)

	for {
		select {
		case <-sm.ctx.Done():
			return
		case tick := <-ticker.C:
			sm.triggerStatusRetrieval <- tick
		}
	}
}

func (sm *Snapmaker) statusLoop() {
	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-sm.triggerStatusRetrieval:
			status, err := sm.getStatusFromPrinter()
			if err != nil {
				fmt.Printf("Status retrieval failed: %v\n", err)
			} else {
				if len(status) > 0 {
					sm.waitChanMutex.Lock()
					if sm.waitForNewStatus != nil {
						sm.waitForNewStatus <- status
					}
					sm.waitChanMutex.Unlock()
				}
			}
		}
	}
}
