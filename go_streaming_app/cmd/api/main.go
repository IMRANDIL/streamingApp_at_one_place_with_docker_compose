package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Video represents the video metadata.
type Video struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func isVideoFileType(filename string) bool {
	// Add more video formats here if needed
	videoFormats := []string{".mp4", ".avi"}

	ext := strings.ToLower(filepath.Ext(filename))
	for _, format := range videoFormats {
		if ext == format {
			return true
		}
	}
	return false
}

// ...

func uploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a multipart request
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("video")

	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Log file metadata
	log.Printf("File Name: %s", handler.Filename)
	log.Printf("File Size: %d bytes", handler.Size)
	log.Printf("Content Type: %s", handler.Header.Get("Content-Type"))

	if handler.Size > 200*1000000 {
		http.Error(w, "Invalid file size. Only 10MB files are allowed.", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the video (you can use a proper ID generator in production)
	videoID := uuid.New().String()
	videoFilename := fmt.Sprintf("%s%s", videoID, filepath.Ext(handler.Filename))

	// Save the uploaded video to a directory with the videoID as the filename
	videoDir := "./videos"
	if _, err := os.Stat(videoDir); os.IsNotExist(err) {
		os.Mkdir(videoDir, 0755)
	}

	filePath := filepath.Join(videoDir, videoFilename)
	outputFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	// Copy the file data to the output file using a buffer
	buffer := make([]byte, 8192) // 8 KB buffer size
	for {
		// Read the data from the file
		n, err := file.Read(buffer)
		if err == io.EOF {
			// End of file, break the loop
			break
		} else if err != nil {
			http.Error(w, "Error reading the file", http.StatusInternalServerError)
			return
		}

		// Write the data to the output file
		_, err = outputFile.Write(buffer[:n])
		if err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}
	}

	// Create the Video object with video metadata
	video := Video{ID: videoID, Name: handler.Filename}

	// Store video metadata (e.g., ID and Name) in a database or other storage
	// Here, we are just logging it for simplicity
	log.Printf("Uploaded video: %s", videoID)

	// Respond with the video metadata
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

func getAllVideosMetadataHandler(w http.ResponseWriter, r *http.Request) {
	// Get a list of all video files in the videos directory
	videoDir := "./videos"
	files, err := os.ReadDir(videoDir)
	if err != nil {
		http.Error(w, "Error reading video directory", http.StatusInternalServerError)
		return
	}

	// Create a slice to store video metadata
	videos := []Video{}

	// Loop through the files and create video metadata for each file
	for _, file := range files {
		if !file.IsDir() && isVideoFileType(file.Name()) {
			videoID := uuid.New().String()
			videoURL := fmt.Sprintf("http://localhost:8080/videos/%s", file.Name())
			video := Video{ID: videoID, Name: file.Name(), URL: videoURL}
			videos = append(videos, video)
		}
	}

	// Respond with the list of video metadata
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func streamVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the video filename from the URL parameters
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Open the video file
	videoDir := "./videos"
	filePath := filepath.Join(videoDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get the video file size
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error reading video file", http.StatusInternalServerError)
		return
	}
	fileSize := fileInfo.Size()

	// Set the content type as video/mp4 for MP4 videos
	w.Header().Set("Content-Type", "video/mp4")

	// Check if the request supports seeking
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		// Request supports seeking, so set the Accept-Ranges and Content-Range headers
		w.Header().Set("Accept-Ranges", "bytes")

		// Parse the range header to get the start and end positions for byte range
		rangeValue := strings.Split(rangeHeader, "=")[1]
		byteRange := strings.Split(rangeValue, "-")
		start, _ := strconv.ParseInt(byteRange[0], 10, 64)

		// If the end position is not specified, set it to the end of the file
		end := fileSize - 1
		if byteRange[1] != "" {
			end, _ = strconv.ParseInt(byteRange[1], 10, 64)
		}

		// Set the Content-Range header to specify the byte range being served
		contentRange := fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize)
		w.Header().Set("Content-Range", contentRange)

		// Set the status code to 206 Partial Content to indicate a partial response
		w.WriteHeader(http.StatusPartialContent)

		// Seek to the start position in the video file
		_, err = file.Seek(start, io.SeekStart)
		if err != nil {
			http.Error(w, "Error seeking video file", http.StatusInternalServerError)
			return
		}
	}

	// Stream the video content to the response writer
	io.Copy(w, file)
}

// ...

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadVideoHandler).Methods("POST")
	r.HandleFunc("/videos", getAllVideosMetadataHandler).Methods("GET")   // Add this line
	r.HandleFunc("/videos/{filename}", streamVideoHandler).Methods("GET") // New route for video streaming

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	handler := c.Handler(r)

	port := ":8080"
	fmt.Println("Server running on port", port)
	err := http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
