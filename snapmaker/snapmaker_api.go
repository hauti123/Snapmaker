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

func connect() {
	client := &http.Client{}
	//	c.Post(url, ,)

	req, err := http.NewRequest("POST",
		"http://192.168.188.130:8080/api/v1/connect",
		strings.NewReader(url.Values{"token": {token}}.Encode()))

	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)

	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Printf("%s\n", string(body[:n]))
}

func status() string {
	client := &http.Client{}

	req, err := http.NewRequest("GET",
		"http://192.168.188.130:8080/api/v1/status?token="+token, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body := make([]byte, 1024)
	n, err := resp.Body.Read(body)

	if err != nil && err != io.EOF {
		panic(err)
	}

	return string(body[:n])
}

func sendFile(filePath string) {

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
		panic(err)
	}
	io.Copy(tokenPart, strings.NewReader(token))

	filePart, _ := multipartWriter.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(filePart, file)
	multipartWriter.Close()

	multipartRequest, _ := http.NewRequest("POST", "http://192.168.188.130:8080/api/v1/upload", body)
	multipartRequest.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	multipartRequest.Close = true
	client := &http.Client{}

	fmt.Println("start request")
	resp, err := client.Do(multipartRequest)

	if err != nil {
		panic(err)
	}

	responseBody := make([]byte, 1024)
	n, err := resp.Body.Read(responseBody)

	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Printf("recevied %d bytes\nstatus: %s\n", n, resp.Status)
	fmt.Printf("upload response:\n%s\n", string(responseBody[:n]))
}
