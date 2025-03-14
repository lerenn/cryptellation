package react

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// AddRoutes serves the static file system for the UI React App.
func AddRoutes(staticFS embed.FS, router gin.IRouter) {
	embeddedBuildFolder := newStaticFileSystem(staticFS)
	fallbackFileSystem := newFallbackFileSystem(embeddedBuildFolder)
	router.Use(static.Serve("/", embeddedBuildFolder))
	router.Use(static.Serve("/", fallbackFileSystem))
}

// ----------------------------------------------------------------------
// staticFileSystem serves files out of the embedded build folder

type staticFileSystem struct {
	staticFS embed.FS
	http.FileSystem
}

var _ static.ServeFileSystem = (*staticFileSystem)(nil)

func newStaticFileSystem(staticFS embed.FS) *staticFileSystem {
	sub, err := fs.Sub(staticFS, "build")

	if err != nil {
		panic(err)
	}

	return &staticFileSystem{
		staticFS:   staticFS,
		FileSystem: http.FS(sub),
	}
}

func (s *staticFileSystem) Exists(_ string, path string) bool {
	buildpath := fmt.Sprintf("build%s", path)

	// support for folders
	if strings.HasSuffix(path, "/") {
		_, err := s.staticFS.ReadDir(strings.TrimSuffix(buildpath, "/"))
		return err == nil
	}

	// support for files
	f, err := s.staticFS.Open(buildpath)
	if f != nil {
		_ = f.Close()
	}
	return err == nil
}

// fallbackFileSystem wraps a staticFileSystem and always serves /index.html.
type fallbackFileSystem struct {
	staticFileSystem *staticFileSystem
}

var _ static.ServeFileSystem = (*fallbackFileSystem)(nil)
var _ http.FileSystem = (*fallbackFileSystem)(nil)

func newFallbackFileSystem(staticFileSystem *staticFileSystem) *fallbackFileSystem {
	return &fallbackFileSystem{
		staticFileSystem: staticFileSystem,
	}
}

func (f *fallbackFileSystem) Open(_ string) (http.File, error) {
	return f.staticFileSystem.Open("/index.html")
}

func (f *fallbackFileSystem) Exists(_ string, _ string) bool {
	return true
}
