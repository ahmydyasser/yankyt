package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/text/unicode/norm"
)

type VideoEntry struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

type Playlist struct {
	Channel string       `json:"channel"`
	Title   string       `json:"title"`
	Entries []VideoEntry `json:"entries"`
}

// Normalize the title by removing special characters and spaces
func normalizeTitle(title string) string {
	// Remove leading and trailing spaces
	title = strings.TrimSpace(title)
	// Normalize Unicode characters
	title = norm.NFKC.String(title)
	// Convert to lowercase
	title = strings.ToLower(title)
	// Replace spaces with empty string
	title = strings.ReplaceAll(title, " ", "")
	// Remove special characters
	title = strings.ReplaceAll(title, "|", "")
	title = strings.ReplaceAll(title, "/", "")
	title = strings.ReplaceAll(title, "\\", "")
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, "*", "")
	title = strings.ReplaceAll(title, "?", "")
	title = strings.ReplaceAll(title, "\"", "")
	title = strings.ReplaceAll(title, "<", "")
	title = strings.ReplaceAll(title, ">", "")
	title = strings.ReplaceAll(title, ".", "")
	title = strings.ReplaceAll(title, ",", "")
	return title
}

// Generate a hash of the normalized title
func hashTitle(title string) string {
	hasher := sha1.New()
	hasher.Write([]byte(title))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Check if a file is completely downloaded
func isFileComplete(filePath string) bool {
	if strings.HasSuffix(filePath, ".part") || strings.HasSuffix(filePath, ".aria2") {
		return false
	}
	return true
}

// Remove all extensions from a filename
func removeAllExtensions(filename string) string {
	for ext := filepath.Ext(filename); ext != ""; ext = filepath.Ext(filename) {
		filename = strings.TrimSuffix(filename, ext)
	}
	return filename
}

func main() {
	// Define flags for the playlist URL and output directory
	playlistURLFlag := flag.String("u", "", "URL of the YouTube playlist")
	outputDirFlag := flag.String("o", "", "Output directory (default: current directory)")
	flag.Parse()

	// Check if the playlist URL flag is provided
	if *playlistURLFlag == "" {
		log.Fatal("Please provide the URL of the YouTube playlist using the -u flag.")
	}

	// Set the output directory
	outputDir := *outputDirFlag
	if outputDir == "" {
		outputDir, _ = os.Getwd()
	}

	// Run yt-dlp command to get the playlist JSON
	cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", *playlistURLFlag)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running yt-dlp: %v", err)
	}

	// Parse the JSON output
	var playlist Playlist
	err = json.Unmarshal(output, &playlist)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Make directory called by Playlist.Title
	dirName := playlist.Title + "_BY_" + playlist.Channel
	dirPath := filepath.Join(outputDir, dirName)

	// Check if the directory already exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create the directory
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
		log.Printf("Directory created successfully: %s", dirPath)
	} else {
		log.Printf("Directory already exists: %s", dirPath)
	}

	// Get the list of files in the directory using filepath.WalkDir
	var files []os.DirEntry
	err = filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, d)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	// Create a map of existing complete video hashes
	existingVideos := make(map[string]bool)
	var mu sync.Mutex
	for _, file := range files {
		if isFileComplete(file.Name()) {
			name := file.Name()
			title := removeAllExtensions(name)
			normalizedTitle := normalizeTitle(title)
			hash := hashTitle(normalizedTitle)
			mu.Lock()
			existingVideos[hash] = true
			mu.Unlock()
			log.Printf("Detected existing file: %s (normalized title: %s, hash: %s)", name, normalizedTitle, hash)
		}
	}

	// Function to download a video
	downloadVideo := func(entry VideoEntry, wg *sync.WaitGroup) {
		defer wg.Done()

		// Normalize and hash the title
		normalizedTitle := normalizeTitle(entry.Title)
		hash := hashTitle(normalizedTitle)

		mu.Lock()
		_, exists := existingVideos[hash]
		mu.Unlock()

		// Check if the video is already downloaded or partially downloaded
		if !exists {
			log.Printf("Downloading: %s (normalized title: %s, hash: %s)", entry.Title, normalizedTitle, hash)
			// Download the video using yt-dlp with aria2c
			cmd := exec.Command("yt-dlp", "-f", "bv*[height<=1080]+ba/b", "--merge-output-format", "mp4", "--downloader", "aria2c", "--downloader-args", "aria2c:-x 16 -s 16 -k 1M", "-o", filepath.Join(dirPath, "%(title)s.%(ext)s"), entry.Url)
			err := cmd.Run()
			if err != nil {
				log.Printf("Error downloading video %s: %v", entry.Title, err)
			}
		} else {
			log.Printf("Already downloaded: %s (normalized title: %s, hash: %s)", entry.Title, normalizedTitle, hash)
		}
	}

	// Use a WaitGroup to manage concurrency
	var wg sync.WaitGroup

	// Loop through the playlist entries and download missing videos
	for _, entry := range playlist.Entries {
		wg.Add(1)
		go downloadVideo(entry, &wg)
	}

	// Wait for all downloads to complete
	wg.Wait()
	log.Println("All downloads complete.")

	// Send a desktop notification
	err = exec.Command("notify-send", "Download Complete", "All videos have been downloaded.").Run()
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}
}

