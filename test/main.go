package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/ts"
)

var (
	videoSourcePath = "videos/source_video.mp4"
	hlsDirectory    = "temp/hls"
)

func main() {
	http.HandleFunc("/hls/", handleHLSRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHLSRequest(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("id")
	if videoID == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}

	hlsDir := filepath.Join(hlsDirectory, videoID)
	playlistPath := filepath.Join(hlsDir, "index.m3u8")

	// Ensure the HLS directory exists
	if _, err := os.Stat(hlsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(hlsDir, os.ModePerm); err != nil {
			http.Error(w, "Failed to create HLS directory", http.StatusInternalServerError)
			log.Println("Error creating HLS directory:", err)
			return
		}
	}

	// Generate the HLS manifest file if it doesn't exist
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		if err := generateHLSManifest(videoID, hlsDir); err != nil {
			http.Error(w, "Failed to generate HLS manifest", http.StatusInternalServerError)
			log.Println("Error generating HLS manifest:", err)
			return
		}
	}

	// Serve the requested file
	requestedFile := filepath.Join(hlsDir, r.URL.Path[len("/hls/"+videoID+"/"):])
	http.ServeFile(w, r, requestedFile)
}

func generateHLSManifest(videoID, hlsDir string) error {
	inputFile, err := os.Open(videoSourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source video: %v", err)
	}
	defer inputFile.Close()

	// Create a new demuxer for the input video file
	demuxer, err := avutil.NewDemuxer(inputFile)
	if err != nil {
		return fmt.Errorf("failed to create demuxer: %v", err)
	}
	defer demuxer.Close()

	// Create a new segmenter for HLS
	segmenter := ts.NewMuxer(hlsDir)
	defer segmenter.Close()

	// Create a publisher-subscriber for the video stream
	pubsub := pubsub.NewDemuxerPublisher(demuxer, false)
	defer pubsub.Close()

	// Write segments to the segmenter
	for {
		pkt, err := pubsub.ReadPacket()
		if err != nil {
			break
		}
		segmenter.WritePacket(pkt)
	}

	// Finalize the segmenter
	if err := segmenter.Close(); err != nil {
		return fmt.Errorf("failed to finalize segmenter: %v", err)
	}

	return nil
}
