package snapmaker

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUploadE2E(t *testing.T) {

	now := time.Now()
	patch := monkey.Patch(time.Now, func() time.Time { return now })
	defer patch.Unpatch()

	mockCtrl := gomock.NewController(t)
	mockHttpClient := NewMockHttpClient(mockCtrl)

	printer := NewSnapmaker("1.2.3.4", "mytoken").WithHttpClient(mockHttpClient)

	connectRequest, err := http.NewRequest("POST", "http://1.2.3.4:8080/api/v1/connect", nil)
	require.NoError(t, err)
	connectRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	connectCall := mockHttpClient.EXPECT().Do(NewHttpRequestMatcher(connectRequest)).DoAndReturn(
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(sampleConnectResponse)),
			}, nil

		},
	)

	uploadRequest, err := http.NewRequest("POST", "http://1.2.3.4:8080/api/v1/upload", nil)
	require.NoError(t, err)
	uploadRequest.Header.Add("Content-Type", "multipart/form-data; boundary=----------------------------268923783128719097072428")

	mockHttpClient.EXPECT().
		Do(NewHttpRequestMatcher(uploadRequest)).
		Times(1).
		After(connectCall).
		DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(sampleStatus)),
				}, nil
			},
		)

	statusRequest, err := http.NewRequest("GET",
		fmt.Sprintf("http://1.2.3.4:8080/api/v1/status?token=mytoken&%d", now.UnixMilli()), nil)
	require.NoError(t, err)
	statusRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	mockHttpClient.EXPECT().
		Do(NewHttpRequestMatcher(statusRequest)).
		AnyTimes().
		After(connectCall).
		DoAndReturn(
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(sampleStatus)),
				}, nil
			},
		)

	err = printer.Connect()
	require.NoError(t, err)

	file, err := createTmpGcodeFile()
	defer os.Remove(file.Name())

	err = printer.SendGcodeFile(file.Name())
	require.NoError(t, err)

	//TODO: It would make sense to also check if the upload body contains the correct token and GCODE data within a correct multi-part-form
}

func createTmpGcodeFile() (*os.File, error) {
	file, err := ioutil.TempFile("", "snaptest-")
	if err != nil {
		return nil, err
	}

	file.Write([]byte(`FLAVOR:Marlin
	;TIME:52383
	;Filament used: 13.0944m
	;Layer height: 0.08
	;MINX:123.287
	;MINY:140.882
	;MINZ:0.15
	;MAXX:196.712
	;MAXY:209.939
	;MAXZ:39.03
	;Generated with Cura_SteamEngine 5.0.0
	M82 ;absolute extrusion mode
	M104 S220 ;Set Hotend Temperature
	M140 S70 ;Set Bed Temperature
	G28 ;home
	G90 ;absolute positioning
	G1 X-10 Y-10 F3000 ;Move to corner 
	G1 Z0 F1800 ;Go to zero offset
	M109 S220 ;Wait for Hotend Temperature
	M190 S70 ;Wait for Bed Temperature
	G92 E0 ;Zero set extruder position
	G1 E20 F200 ;Feed filament to clear nozzle
	G92 E0 ;Zero set extruder position
	G92 E0
	G92 E0
	G1 F2700 E-6
	;LAYER_COUNT:487
	;LAYER:0
	M107
	G0 F1440 X135.381 Y145.993 Z0.15
	G0 X129.301 Y153.091
	;TYPE:SKIRT
	G1 F2700 E0
	G1 F1080 X130.063 Y151.976 E0.04248
	G1 X130.85 Y150.91 E0.08416
	G1 X131.664 Y149.883 E0.12538
	G1 X132.548 Y148.852 E0.1681
	G1 X133.429 Y147.893 E0.20907`))

	return file, err
}
