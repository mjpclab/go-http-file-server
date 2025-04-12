package serverHandler

import (
	"archive/tar"
	"io"
	"net/http"
	"os"
	"strings"
)

// tarArchiver implements the Archiver interface for tar format archives.
// It provides methods to create a tar writer and specify the format's extension and MIME type.
type tarArchiver struct{}

// Create initializes a new tar archive writer for the given output writer.
// Returns a tar-specific ArchiveWriter or an error if creation fails.
func (t *tarArchiver) Create(w io.Writer) (ArchiveWriter, error) {
	return &tarWriter{tw: tar.NewWriter(w)}, nil
}

// Extension returns the file extension for tar archives (".tar").
func (t *tarArchiver) Extension() string {
	return ".tar"
}

// MimeType returns the MIME type for tar archives ("application/x-tar").
func (t *tarArchiver) MimeType() string {
	return "application/x-tar"
}

// tarWriter implements the ArchiveWriter interface for tar and tar.gz archives.
// It handles adding files to a tar archive and closing the writer, with optional gzip compression.
type tarWriter struct {
	tw     *tar.Writer
	closer io.Closer // for gzip writer
}

// AddFile adds a file or directory to the tar archive at the specified path.
// It normalizes the path, writes a tar header, and copies file contents if applicable.
// Returns an error if the operation fails.
func (w *tarWriter) AddFile(file *os.File, fileInfo os.FileInfo, archivePath string) error {
	archivePath = strings.TrimPrefix(archivePath, "/")

	var typeFlag byte
	var mode int64
	var size int64
	if fileInfo.IsDir() {
		archivePath += "/"
		typeFlag = tar.TypeDir
		mode = 0755
	} else {
		typeFlag = tar.TypeReg
		mode = 0644
		size = fileInfo.Size()
	}

	header := &tar.Header{
		Name:       archivePath,
		Typeflag:   typeFlag,
		Mode:       mode,
		Size:       size,
		ModTime:    fileInfo.ModTime(),
		AccessTime: fileInfo.ModTime(),
		ChangeTime: fileInfo.ModTime(),
	}

	if err := w.tw.WriteHeader(header); err != nil {
		return err
	}

	if fileInfo.IsDir() || file == nil || size == 0 {
		return nil
	}

	_, err := io.Copy(w.tw, file)
	return err
}

// Close finalizes the tar archive and, if applicable, the gzip compression layer.
// Returns an error if closing fails.
func (w *tarWriter) Close() error {
	if err := w.tw.Close(); err != nil {
		return err
	}
	if w.closer != nil {
		return w.closer.Close()
	}
	return nil
}

// tar handles HTTP requests to create a tar archive of selected files.
// It delegates to baseArchiveHandler with a tar-specific archiver.
// Returns false if validation fails, true on success.
func (h *aliasHandler) tar(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) bool {
	return h.baseArchiveHandler(w, r, session, data, &tarArchiver{})
}
