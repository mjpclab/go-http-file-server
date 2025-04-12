package serverHandler

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
)

// tgzArchiver implements the Archiver interface for gzip-compressed tar archives.
// It provides methods to create a tar.gz writer and specify the format's extension and MIME type.
type tgzArchiver struct{}

// Create initializes a new tar.gz archive writer with gzip compression for the given output writer.
// Returns a tar-specific ArchiveWriter or an error if creation fails.
func (t *tgzArchiver) Create(w io.Writer) (ArchiveWriter, error) {
	gzipWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	return &tarWriter{
		tw:     tar.NewWriter(gzipWriter),
		closer: gzipWriter,
	}, nil
}

// Extension returns the file extension for tar.gz archives (".tar.gz").
func (t *tgzArchiver) Extension() string {
	return ".tar.gz"
}

// MimeType returns the MIME type for tar.gz archives ("application/gzip").
func (t *tgzArchiver) MimeType() string {
	return "application/gzip"
}

// tgz handles HTTP requests to create a gzip-compressed tar archive of selected files.
// It delegates to baseArchiveHandler with a tar.gz-specific archiver.
// Returns false if validation fails, true on success.
func (h *aliasHandler) tgz(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData) bool {
	return h.baseArchiveHandler(w, r, session, data, &tgzArchiver{})
}
