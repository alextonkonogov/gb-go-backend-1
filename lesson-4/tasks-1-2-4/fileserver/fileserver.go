package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	UploadDir string
}

func (h *Handler) ServeFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, h.UploadDir)
}

func main() {
	filesrv := &Handler{
		UploadDir: "../upload",
	}

	http.Handle("/", http.HandlerFunc(filesrv.ServeFiles))

	fs := &http.Server{
		Addr:         ":9985",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("File-Server is started. Try it on http://localhost%s\n", fs.Addr)
	err := fs.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
