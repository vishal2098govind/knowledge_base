package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
)

func main() {
	err := Kill("server.pid")
	if err != nil {
		fmt.Println("ERROR", err)
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("not found")
		}
		// chaining wrapped err
		for e := err; e != nil; e = errors.Unwrap(e) {
			fmt.Println("error: ", e)
		}
	}
}

func Kill(pidFile string) error {
	file, err := os.Open(pidFile)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			slog.Warn("close", "file", pidFile, "error", err)
		}
	}()

	var pid int

	if _, err := fmt.Fscanf(file, "%d", &pid); err != nil {
		return fmt.Errorf("%q - bad pid, %w", pidFile, err)
	}

	slog.Info("killing", "pid", pid)
	if err := os.Remove(pidFile); err != nil {
		slog.Warn("delete", "file", pidFile, "error", err)
	}

	return nil
}
