package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	h := &UploadHandler{
		UploadDir: "../upload",
	}
	http.Handle("/", http.HandlerFunc(h.IndexHandler))
	http.Handle("/upload", http.HandlerFunc(h.UploadFile))
	http.Handle("/files", http.HandlerFunc(h.GetFiles))

	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Upload-Server started on", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func (h *UploadHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}
}

type file struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Size int64  `json:"size"`
}

func (f file) GetExtension() string {
	els := strings.Split(f.Name, ".")
	if len(els) > 1 {
		f.Ext = els[len(els)-1]
	}
	return f.Ext
}

func (h *UploadHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	files := []file{}
	filesMap := make(map[string][]file)

	err := filepath.Walk(h.UploadDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != ".DS_Store" {
			f := file{}
			f.Name = info.Name()
			f.Ext = f.GetExtension()
			f.Size = info.Size()
			f.Path = path
			files = append(files, f)
			filesMap[f.Ext] = append(filesMap[f.Ext], f)
		}
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ext := strings.Trim(r.FormValue("ext"), ".")
	if ext != "" {
		if match, ok := filesMap[ext]; ok {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(match)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, fmt.Sprintf("files with extension \"%s\" are not found", ext), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}
	f, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	fileName := h.UploadDir + "/" + header.Filename
	err = ioutil.WriteFile(fileName, data, 0600)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fileLink := h.HostAddr + "/" + header.Filename
	_, err = fmt.Fprintln(w, fileLink)
}
