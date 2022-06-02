package serverHandler

import (
	"net/http"
	"os"
	"runtime"
)

type wrappedHttpFile struct {
	http.File
}

func (f wrappedHttpFile) Readdir(n int) (items []os.FileInfo, err error) {
	return f.File.Readdir(n)
}

type wrappedHttpFileSystem struct {
	http.FileSystem
}

func (fs wrappedHttpFileSystem) Open(name string) (http.File, error) {
	file, err := fs.FileSystem.Open(name)
	return wrappedHttpFile{file}, err
}

func createFsFileServer(root string) http.Handler {
	return http.FileServer(wrappedHttpFileSystem{http.Dir(root)})
}

func serveFsContent(h *handler, w http.ResponseWriter, r *http.Request, info os.FileInfo, file *os.File) {
	h.fileServer.ServeHTTP(w, r)
}

func TryEnableWSL1Fix() bool {
	if runtime.GOOS == "linux" && len(os.Getenv("WSL_DISTRO_NAME")) > 0 && len(os.Getenv("WSL_INTEROP")) == 0 {
		createFileServer = createFsFileServer
		serveContent = serveFsContent
		return true
	}
	return false
}
