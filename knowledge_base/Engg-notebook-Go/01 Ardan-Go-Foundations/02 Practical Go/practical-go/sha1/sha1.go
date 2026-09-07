package main

import (
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func main() {
	sig, err := SHA1Sign("test.txt")
	if err != nil {
		fmt.Println("error")
		if errors.Is(err, io.EOF) {
			fmt.Println("end of file")
		}
		return
	}
	fmt.Printf("sig1: %s\n", sig)
	sig2, err := SHA1Sign("test-2.txt")
	if err != nil {
		fmt.Println("error")
		if errors.Is(err, io.EOF) {
			fmt.Println("end of file")
		}
		return
	}
	fmt.Printf("sig2: %s\n", sig2)
}

func SHA1Sign(filename string) (string, error) {

	file, err := os.Open(filename)
	if err != nil {
		slog.Error("open file", "file", filename, "error", err)
		return "", fmt.Errorf("open file: %w", err)
	}

	exts := strings.Split(filename, ".")
	ext := exts[len(exts)-1]
	slog.Info("extension", "ext", ext)
	var isCompressed bool
	if ext == ".gz" {
		slog.Info("compressed")
		isCompressed = true
	}

	var r io.ReadCloser
	if isCompressed {
		r, err = gzip.NewReader(file)
		if err != nil {
			slog.Error("gzip NewReader:", "file", filename, "error", err)
			return "", fmt.Errorf("gzip NewReader: %q - %w", filename, err)
		}
	} else {
		r = file
	}
	defer r.Close()

	w := sha1.New()

	if _, err := io.Copy(w, r); err != nil {
		slog.Error("io copy", "file", filename, "error", err)
		return "", fmt.Errorf("io copy: %q - %w", filename, err)
	}

	sig := w.Sum(nil)

	return fmt.Sprintf("%x", string(sig)), nil
}
