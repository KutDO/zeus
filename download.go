package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var progress = struct {
	sync.RWMutex
	data map[string]int
}{data: make(map[string]int)}

type DownloadTask struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id,omitempty"`
	URL        string `json:"url"`
	Filename   string `json:"filename,omitempty"`
	Decompress bool   `json:"decompress"`
}

// @Summary Create download task
// @Description Create a new download task
// @Tags downloads
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Param url formData string true "URL of the file to download"
// @Param filename formData string false "Optional filename for the downloaded file"
// @Param decompress formData bool false "Decompress the file if true"
// @Param user_id formData string false "Optional user ID for the download"
// @Success 201 {object} DownloadTask
// @Router /downloads [post]
func CreateDownloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create download task")

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	url := r.FormValue("url")
	filename := r.FormValue("filename")
	decompress := r.FormValue("decompress") == "on"
	userID := r.FormValue("user_id")

	if url == "" {
		log.Println("URL is missing")
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	task := DownloadTask{
		ID:         uuid.New().String(),
		UserID:     userID,
		URL:        url,
		Filename:   filename,
		Decompress: decompress,
	}

	log.Printf("Parsed form data: URL=%s, Filename=%s, Decompress=%v, UserID=%s", task.URL, task.Filename, task.Decompress, task.UserID)

	// Set default filename if not provided
	if task.Filename == "" {
		task.Filename = filepath.Base(task.URL)
	}

	// Create user-specific or web-specific download folder
	downloadFolder := config.Paths.WebDirectory
	if task.UserID != "" {
		downloadFolder = filepath.Join(config.Paths.DownloadDirectory, task.UserID)
	}

	// Verify the downloadFolder is not empty
	if downloadFolder == "" {
		log.Println("Download folder is empty")
		http.Error(w, "Internal server error: download folder is not set", http.StatusInternalServerError)
		return
	}

	// Ensure the parent directory exists
	if err := os.MkdirAll(downloadFolder, os.ModePerm); err != nil {
		log.Printf("Error creating download directory: %v", err)
		http.Error(w, "Error creating download directory", http.StatusInternalServerError)
		return
	}

	task.Filename = filepath.Join(downloadFolder, task.Filename)

	// Check if file already exists and append timestamp if necessary
	if _, err := os.Stat(task.Filename); err == nil {
		task.Filename = fmt.Sprintf("%s-%d%s", task.Filename, time.Now().Unix(), filepath.Ext(task.Filename))
	}

	EnqueueTask(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)

	log.Printf("Download task created: %+v", task)
}

func DownloadFile(task *DownloadTask) {
	progress.Lock()
	progress.data[task.ID] = 0
	progress.Unlock()

	resp, err := http.Get(task.URL)
	if err != nil {
		log.Printf("Error downloading file: %v\n", err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(task.Filename)
	if err != nil {
		log.Printf("Error creating file: %v\n", err)
		return
	}
	defer out.Close()

	size := resp.ContentLength
	var downloaded int64
	buf := make([]byte, 1024)

	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Error reading response body: %v\n", err)
			return
		}
		if n == 0 {
			break
		}

		out.Write(buf[:n])
		downloaded += int64(n)

		progress.Lock()
		progress.data[task.ID] = int(float64(downloaded) / float64(size) * 100)
		progress.Unlock()
	}

	if task.Decompress {
		err = Unzip(task.Filename, filepath.Dir(task.Filename))
		if err != nil {
			log.Printf("Error decompressing file: %v\n", err)
			return
		}
	}

	log.Printf("Download completed: %s", task.Filename)
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}

// @Summary Get download progress
// @Description Get the progress of a download task
// @Tags downloads
// @Accept  json
// @Produce  json
// @Param id path string true "Download Task ID"
// @Success 200 {object} map[string]int
// @Router /progress/{id} [get]
func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	progress.RLock()
	prog, ok := progress.data[id]
	progress.RUnlock()

	if !ok {
		http.Error(w, "Progress not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"progress": prog})
}
