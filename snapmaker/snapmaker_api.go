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

//const token = "aaf66ddf-bf4e-4d0e-b608-bb4fa3ff2fb7"
const token = "6a545dbb-3634-44cc-8508-f11a04cbebfe" // Luban token
const snapmakerPort = 8080

func ConnectToPrinter(printerIp string) error {
	client := &http.Client{}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("http://%s:%d/api/v1/connect", printerIp, snapmakerPort),
		strings.NewReader(url.Values{"token": {token}}.Encode()))

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

func GetPrinterStatus(printerIp string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET",
		fmt.Sprintf("http://%s:%d/api/v1/status?token=%s", printerIp, snapmakerPort, token),
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

func SendGcodeFile(printerIp string, filePath string) error {

	fmt.Printf("start upload: %s\n", filePath)
	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(body)
	multipartWriter.SetBoundary("----------------------------268923783128719097072428")

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
		fmt.Sprintf("http://%s:%d/api/v1/upload", printerIp, snapmakerPort),
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
