package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var IndexHandlerTests = []struct {
	method       string
	wantedStatus int
	wantedBody   string
}{
	{method: "GET", wantedStatus: http.StatusOK, wantedBody: ""},
	{method: "POST", wantedStatus: http.StatusMethodNotAllowed, wantedBody: "bad method"},
}

func TestIndexHandler(t *testing.T) {
	for i, tt := range IndexHandlerTests {
		fmt.Println("test", i, tt)
		req, err := http.NewRequest(tt.method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := &UploadHandler{}
		handler.IndexHandler(rr, req)

		// Проверяем статус-код ответа
		if status := rr.Code; status != tt.wantedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.wantedStatus)
		}
		// Проверяем тело ответа
		if strings.Trim(rr.Body.String(), "\n") != tt.wantedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), tt.wantedBody)
		}
	}
}

var UploadHandlerTests = []struct {
	method           string
	fileName         string
	wantedStatus     int
	wantedBodySubStr string
}{
	{
		method:           "POST",
		fileName:         "uploadserver.go",
		wantedStatus:     http.StatusOK,
		wantedBodySubStr: "uploadserver.go",
	},
	{
		method:           "GET",
		fileName:         "uploadserver.go",
		wantedStatus:     http.StatusMethodNotAllowed,
		wantedBodySubStr: "bad method",
	},
}

func TestUploadHandler(t *testing.T) {
	for i, tt := range UploadHandlerTests {
		fmt.Println("test", i, tt)

		f, err := os.Open(tt.fileName)
		if err != nil {
			t.Errorf("cant open file %v", tt.fileName)
			return
		}
		body := &bytes.Buffer{}

		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", filepath.Base(f.Name()))
		io.Copy(part, f)
		writer.Close()
		f.Close()

		req, _ := http.NewRequest(tt.method, "/upload", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "ok!")
		}))
		defer ts.Close()
		uploadHandler := &UploadHandler{
			UploadDir: "../upload",
			HostAddr:  ts.URL,
		}
		uploadHandler.UploadFile(rr, req)
		if status := rr.Code; status != tt.wantedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.wantedStatus)
		}
		if !strings.Contains(rr.Body.String(), tt.wantedBodySubStr) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), tt.wantedBodySubStr)
		}
	}
}

var GetFilesTests = []struct {
	method            string
	params            string
	wantedStatus      int
	wantedContentType string
}{
	{
		method:            "GET",
		params:            "",
		wantedStatus:      http.StatusOK,
		wantedContentType: "application/json",
	},
	{
		method:            "GET",
		params:            "?ext=",
		wantedStatus:      http.StatusOK,
		wantedContentType: "application/json",
	},
	{
		method:            "GET",
		params:            "?ext=txt",
		wantedStatus:      http.StatusOK,
		wantedContentType: "application/json",
	},
	{
		method:            "GET",
		params:            "?ext=sdf123",
		wantedStatus:      http.StatusNotFound,
		wantedContentType: "text/plain; charset=utf-8",
	},
	{
		method:            "POST",
		params:            "?ext=",
		wantedStatus:      http.StatusMethodNotAllowed,
		wantedContentType: "text/plain; charset=utf-8",
	},
}

func TestGetFiles(t *testing.T) {
	for i, tt := range GetFilesTests {
		fmt.Println("test", i, tt)
		req, err := http.NewRequest(tt.method, "/files"+tt.params, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "ok!")
		}))
		handler := &UploadHandler{
			UploadDir: "../upload",
			HostAddr:  ts.URL,
		}
		handler.GetFiles(rr, req)

		if status := rr.Code; status != tt.wantedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.wantedStatus)
		}

		if rr.Header().Get("Content-Type") != tt.wantedContentType {
			t.Errorf("handler returned unexpected content type: got %v want %v",
				rr.Header().Get("Content-Type"), tt.wantedContentType)
		}
	}
}
