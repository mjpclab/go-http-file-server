package serverHandler

import (
	"io"
	"net/http"
	"os"
)

// ChunkSize defines the buffer size (32KB) for chunked file copying during archiving.
const ChunkSize = 32 << 10

// Archiver defines the interface for creating archive writers and specifying format details.
// Implementations provide format-specific writers, file extensions, and MIME types.
type Archiver interface {
	Create(w io.Writer) (ArchiveWriter, error)
	Extension() string
	MimeType() string
}

// ArchiveWriter defines the interface for writing files to an archive and closing it.
// Implementations handle format-specific file addition and cleanup.
type ArchiveWriter interface {
	AddFile(file *os.File, fileInfo os.FileInfo, archivePath string) error
	Close() error
}

// baseArchiveHandler handles common archive operations for HTTP requests.
// It validates the request, sets up the archive writer, and processes files using the provided Archiver.
// Returns false if validation fails, true on success.
func (h *aliasHandler) baseArchiveHandler(w http.ResponseWriter, r *http.Request, session *sessionContext, data *responseData, archiver Archiver) bool {
	if !data.CanArchive {
		data.Status = http.StatusBadRequest
		return false
	}

	selections, ok := h.normalizeArchiveSelections(r)
	if !ok {
		data.Status = http.StatusBadRequest
		return false
	}

	archiveWriter, err := archiver.Create(w)
	if err != nil {
		h.logError(err)
		data.Status = http.StatusInternalServerError
		return false
	}
	defer func() {
		if err := archiveWriter.Close(); err != nil {
			h.logError(err)
		}
	}()

	ctx := r.Context()
	h.archiveFiles(
		w,
		r,
		session,
		data,
		selections,
		archiver.Extension(),
		archiver.MimeType(),
		func(file *os.File, fileInfo os.FileInfo, relPath string) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return archiveWriter.AddFile(file, fileInfo, relPath)
			}
		},
	)
	return true
}
