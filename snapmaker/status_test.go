package snapmaker

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleConnectResponse = `{"token":"mytoken","readonly":false,"series":"Snapmaker 2.0 A350","headType":1,"hasEnclosure":false}`

var sampleStatus = `{
	"status": "IDLE",
	"x": -19,
	"y": 347,
	"z": 326.95,
	"homed": true,
	"offsetX": 1.7,
	"offsetY": 3,
	"offsetZ": 4.23,
	"toolHead": "TOOLHEAD_3DPRINTING_1",
	"nozzleTemperature": 22,
	"nozzleTargetTemperature": 220,
	"heatedBedTemperature": 21,
	"heatedBedTargetTemperature": 80,
	"isFilamentOut": true,
	"workSpeed": 1500,
	"printStatus": "Idle",
	"moduleList": {
		"enclosure": true,
		"rotaryModule": true,
		"emergencyStopButton": true,
		"airPurifier": true
	}
}`

func serveGetStatus(t *testing.T, printer *Snapmaker) {

	select {
	case <-time.After(15 * time.Second):
		assert.FailNow(t, "waiting for status trigger timeout")
	case <-printer.triggerStatusRetrieval:
	}

	printer.waitForNewStatus <- sampleStatus
}

func TestGetStatus(t *testing.T) {
	printer := NewSnapmaker("1.1.1.1", "abcd")

	go serveGetStatus(t, printer)
	status, err := printer.GetStatus(5 * time.Second)
	require.NoError(t, err)

	assert.Equal(t, "IDLE", status.Status)
	assert.Equal(t, float64(-19), status.X)
	assert.Equal(t, float64(347), status.Y)
	assert.Equal(t, float64(326.95), status.Z)
	assert.Equal(t, true, status.Homed)
	assert.Equal(t, float64(1.7), status.OffsetX)
	assert.Equal(t, float64(3), status.OffsetY)
	assert.Equal(t, float64(4.23), status.OffsetZ)
	assert.Equal(t, "TOOLHEAD_3DPRINTING_1", status.ToolHead)
	assert.Equal(t, 22, status.NozzleTemperature)
	assert.Equal(t, 220, status.NozzleTargetTemperature)
	assert.Equal(t, 21, status.HeatedBedTemperature)
	assert.Equal(t, 80, status.HeatedBedTargetTemperature)
	assert.Equal(t, true, status.IsFilamentOut)
	assert.Equal(t, 1500, status.WorkSpeed)
	assert.Equal(t, "Idle", status.PrintStatus)
	assert.Equal(t, true, status.ModuleList.Enclosure)
	assert.Equal(t, true, status.ModuleList.RotaryModule)
	assert.Equal(t, true, status.ModuleList.EmergencyStopButton)
	assert.Equal(t, true, status.ModuleList.AirPurifier)
}

func TestGetStatusE2E(t *testing.T) {

	now := time.Now()
	patch := monkey.Patch(time.Now, func() time.Time { return now })
	defer patch.Unpatch()

	mockCtrl := gomock.NewController(t)
	mockHttpClient := NewMockHttpClient(mockCtrl)

	printer := NewSnapmaker("1.2.3.4", "mytoken").WithHttpClient(mockHttpClient)

	connectRequest, err := http.NewRequest("POST", "http://1.2.3.4:8080/api/v1/connect", nil)
	require.NoError(t, err)
	connectRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	connectCall := mockHttpClient.EXPECT().Do(NewHttpRequestMatcher(connectRequest)).Times(1).DoAndReturn(
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(sampleConnectResponse)),
			}, nil
		},
	)

	statusRequest, err := http.NewRequest("GET",
		fmt.Sprintf("http://1.2.3.4:8080/api/v1/status?token=mytoken&%d", now.UnixMilli()), nil)
	require.NoError(t, err)
	statusRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	statusCall := mockHttpClient.EXPECT().
		Do(NewHttpRequestMatcher(statusRequest)).
		AnyTimes().
		DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(sampleStatus)),
				}, nil
			},
		)

	gomock.InOrder(connectCall, statusCall)

	err = printer.Connect()
	require.NoError(t, err)

	_, err = printer.GetStatus(5 * time.Second)
	require.NoError(t, err)
}
