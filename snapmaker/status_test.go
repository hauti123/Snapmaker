package snapmaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serverGetStatus(t *testing.T, printer *Snapmaker) {

	sampleStatus := `{
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

	select {
	case <-time.After(15 * time.Second):
		assert.FailNow(t, "waiting for status trigger timeout")
	case <-printer.triggerStatusRetrieval:
	}

	printer.waitForNewStatus <- sampleStatus
}

func TestGetStatus(t *testing.T) {
	printer := NewSnapmaker("1.1.1.1", "abcd")

	go serverGetStatus(t, printer)
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
