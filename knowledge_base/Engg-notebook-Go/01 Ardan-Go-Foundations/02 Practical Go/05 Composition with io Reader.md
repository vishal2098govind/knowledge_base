#engineers-notebook #io-Reader #io-Writer 

```go
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
```