package snapmaker

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Snapmaker struct {
	ipAdress string
	port     int
	token    string
}

const apiConnect = "connect"
const apiPrinterStatus = "status"
const apiGcodeUpload = "upload"
const apiEnclosureStatus = "enclosure"

const multipartBoundary = "----------------------------268923783128719097072428"

func (sm *Snapmaker) Connect() error {
	client := &http.Client{}

	req, err := http.NewRequest("POST",
		sm.buildApiUrl(apiConnect),
		strings.NewReader(url.Values{"token": {sm.token}}.Encode()))

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)

	if err != nil && err != io.EOF {
		return err
	}

	fmt.Printf("%s\n", string(body[:n]))
	return nil
}

func (sm *Snapmaker) GetStatus() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET",
		sm.buildApiUrl(apiPrinterStatus)+"?token="+sm.token,
		nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)

	if err != nil && err != io.EOF {
		return "", err
	}

	return string(body[:n]), nil
}

func (sm *Snapmaker) SendGcodeFile(filePath string) error {

	fmt.Printf("start upload: %s\n", filePath)
	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(body)
	multipartWriter.SetBoundary(multipartBoundary)

	tokenHeader := make(textproto.MIMEHeader)
	tokenHeader.Set("Content-Disposition", `form-data; name="token"`)

	tokenPart, err := multipartWriter.CreatePart(tokenHeader)
	if err != nil {
		return err
	}
	io.Copy(tokenPart, strings.NewReader(token))

	filePart, _ := multipartWriter.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(filePart, file)
	multipartWriter.Close()

	multipartRequest, _ := http.NewRequest("POST",
		sm.buildApiUrl(apiGcodeUpload),
		body)

	multipartRequest.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	multipartRequest.Close = true
	client := &http.Client{}

	fmt.Println("start request")
	resp, err := client.Do(multipartRequest)

	if err != nil {
		return err
	}

	responseBody := make([]byte, 1024)
	n, err := resp.Body.Read(responseBody)

	if err != nil && err != io.EOF {
		return err
	}

	fmt.Printf("recevied %d bytes\nstatus: %s\n", n, resp.Status)
	fmt.Printf("upload response:\n%s\n", string(responseBody[:n]))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed, status = %s", resp.Status)
	}
	return nil
}

func (sm *Snapmaker) buildApiUrl(api string) string {
	return fmt.Sprintf("http://%s:%d/api/v1/%s", sm.ipAdress, sm.port, api)
}
