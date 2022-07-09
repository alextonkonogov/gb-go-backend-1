package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var tests = []struct {
	method            string
	wantedStatus      int
	wantedContentType string
}{
	{method: "GET", wantedStatus: http.StatusOK, wantedContentType: "text/html; charset=utf-8"},
}

func TestFileServer(t *testing.T) {
	for i, tt := range tests {
		fmt.Println("test", i, tt)
		req, err := http.NewRequest(tt.method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := &Handler{
			UploadDir: "../upload",
		}
		handler.ServeFiles(rr, req)

		// Проверяем статус-код ответа
		if status := rr.Code; status != tt.wantedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.wantedStatus)
		}
		// Проверяем тело ответа
		if rr.Header().Get("Content-Type") != tt.wantedContentType {
			t.Errorf("handler returned unexpected content type: got %v want %v",
				rr.Header().Get("Content-Type"), tt.wantedContentType)
		}
	}
}
