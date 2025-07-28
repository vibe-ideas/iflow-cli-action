package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	err := downloadAndExtract("https://cloud.iflow.cn/iflow-cli/iflow-iflow-cli-0.0.2.tgz", "iflow-cli")
	if err != nil {
		fmt.Printf("Error downloading and extracting iFlow CLI: %s\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("package/bin/iflow", os.Getenv("INPUT_COMMAND"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running iFlow CLI: %s\n", err)
		os.Exit(1)
	}
}

func downloadAndExtract(url, targetDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tr); err != nil {
				return err
			}
		}
	}

	return nil
}
