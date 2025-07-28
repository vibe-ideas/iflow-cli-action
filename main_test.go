package main

import (
	"archive/tar"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadAndExtract(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new gzip writer
		gw := gzip.NewWriter(w)
		defer gw.Close()

		// Create a new tar writer
		tw := tar.NewWriter(gw)
		defer tw.Close()

		// Add a file to the tar archive
		fileHeader := &tar.Header{
			Name: "test.txt",
			Mode: 0600,
			Size: int64(len("hello world")),
		}
		if err := tw.WriteHeader(fileHeader); err != nil {
			t.Fatalf("Failed to write header: %v", err)
		}
		if _, err := tw.Write([]byte("hello world")); err != nil {
			t.Fatalf("Failed to write file content: %v", err)
		}
	}))
	defer server.Close()

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to the temporary directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Call the function
	if err := downloadAndExtract(server.URL, "test-dir"); err != nil {
		t.Fatalf("downloadAndExtract failed: %v", err)
	}

	// Check if the file was created
	content, err := os.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != "hello world" {
		t.Fatalf("File content is not correct: got %s, want %s", string(content), "hello world")
	}
}
