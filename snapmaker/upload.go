package snapmaker

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

func (sm *Snapmaker) SendGcodeFile(filePath string) error {

	fmt.Printf("starting upload: %s\n", filePath)
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
	io.Copy(tokenPart, strings.NewReader(sm.token))

	filePart, _ := multipartWriter.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(filePart, file)
	multipartWriter.Close()

	multipartRequest, _ := http.NewRequest("POST",
		sm.buildApiUrl(apiGcodeUpload),
		body)

	multipartRequest.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	multipartRequest.Close = true

	resp, err := sm.httpClient.Do(multipartRequest)

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
