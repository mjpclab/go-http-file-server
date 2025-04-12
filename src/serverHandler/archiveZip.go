package serverHandler

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"strings"
)

// zipArchiver implements the Archiver interface for ZIP format archives.
// It provides methods to create a ZIP writer and specify the format's extension and MIME type.
type zipArchiver struct{}

// Create initializes a new ZIP archive writer for the given output writer.
// Returns a ZIP-specific ArchiveWriter or an error if creation fails.
func (z *zipArchiver) Create(w io.Writer) (ArchiveWriter, error) {
	return &zipWriter{zw: zip.NewWriter(w)}, nil
}

// Extension returns the file extension for ZIP archives (".zip").
func (z *zipArchiver) Extension() string {
	return ".zip"
}

// MimeType returns the MIME type for ZIP archives ("application/zip").
func (z *zipArchiver) MimeType() string {
	return "application/zip"
}

// zipWriter implements the ArchiveWriter interface for ZIP format archives.
// It handles adding files to a ZIP archive and closing the writer.
type zipWriter struct {
	zw *zip.Writer
}

// AddFile adds a file or directory to the ZIP archive at the specified path.
// It normalizes the path, creates a ZIP entry, and copies file contents in chunks if applicable.
// Returns an error if the operation fails.
func (w *zipWriter) AddFile(file *os.File, fileInfo os.FileInfo, archivePath string) error {
	archivePath = strings.TrimPrefix(archivePath, "/")
	if fileInfo.IsDir() {
		archivePath += "/"
	}

	entryWriter, err := w.zw.Create(archivePath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() || file == nil || fileInfo.Size() == 0 {
		return nil
	}

	remaining := fileInfo.Size()
	for remaining > 0 {
		chunk := remaining
		if chunk > ChunkSize {
			chunk = ChunkSize
		}
		written, err := io.CopyN(entryWriter, file, chunk)
		if err != nil {
			return err
		}
		remaining -= written
	}
	return nil
}

// Close finalizes the ZIP archive and releases resources.
// Returns an error if closing fails.
func (w *zipWriter) Close() error {
	return w.zw.Close()
}

// zip handles HTTP requests to create a ZIP archive of selected files.
// It delegates to baseArchiveHandler with a ZIP-specific archiver.
// Returns false if validation fails, true on success.
func (h *aliasHandler) zip(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) bool {
	return h.baseArchiveHandler(w, r, session, data, &zipArchiver{})
}
