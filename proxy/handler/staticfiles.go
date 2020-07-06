package handler

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/gofoomo/proxy/utils"
)

// StaticFiles handles serving static files from the local file system. It knows about
// foomos hierarchy and serves files from the htdocs directories of modules.
// Currently it will also serve files of disabled modules.
type StaticFiles struct {
	foomo *foomo.Foomo
}

// NewStaticFiles constructor
func NewStaticFiles(foomo *foomo.Foomo) *StaticFiles {
	sf := new(StaticFiles)
	sf.foomo = foomo
	return sf
}

// HandlesRequest request handler implementation
func (files *StaticFiles) HandlesRequest(incomingRequest *http.Request) bool {
	if strings.HasPrefix(incomingRequest.URL.Path, "/foomo/modulesVar/") {
		return true
	}
	if strings.HasPrefix(incomingRequest.URL.Path, "/foomo/modules/") {
		parts := strings.Split(incomingRequest.URL.Path, "/")
		if len(parts) > 3 {
			moduleNameParts := strings.Split(parts[3], "-")
			if strings.HasSuffix(parts[len(parts)-1], ".php") {
				return false
			}
			return fileExists(files.foomo.GetModuleHtdocsDir(moduleNameParts[0]) + "/" + strings.Join(parts[4:], "/"))
		}
		return false
	}
	return false
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (files *StaticFiles) ServeHTTP(w http.ResponseWriter, incomingRequest *http.Request) {
	parts := strings.Split(incomingRequest.URL.Path, "/")
	path := strings.Join(parts[4:], "/")
	moduleNameParts := strings.Split(parts[3], "-")
	moduleName := moduleNameParts[0]
	var moduleDir string
	if strings.HasPrefix(incomingRequest.URL.Path, "/foomo/modules/") {
		moduleDir = files.foomo.GetModuleHtdocsDir(moduleName)
	} else {
		moduleDir = files.foomo.GetModuleHtdocsVarDir(moduleName)
	}
	fullName := filepath.Join(moduleDir, path)
	// validate path
	absPath, errAbs := filepath.Abs(fullName)
	if errAbs != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if !strings.HasPrefix(absPath, moduleDir) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	// check file
	fileInfo, errStat := os.Stat(fullName)
	if errStat != nil {
		switch true {
		case os.IsNotExist(errStat):
			http.Error(w, "not found", http.StatusNotFound)
		case os.IsPermission(errStat):
			http.Error(w, "forbidden", http.StatusForbidden)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	// open it
	f, errOpen := os.Open(fullName)
	defer f.Close()
	if errOpen != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	// compression support
	_, compress := getContentType(path)
	w.Header().Set("Expires", time.Now().Add(time.Hour*24*365).Format(http.TimeFormat))
	if compress && strings.Contains(incomingRequest.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		crw := utils.NewCompressedResponseWriter(w)
		defer crw.Close()
		w = crw
	}
	// passing it on to std library
	http.ServeContent(w, incomingRequest, f.Name(), fileInfo.ModTime(), f)
}

func getContentType(path string) (mimeType string, compress bool) {
	parts := strings.Split(path, ".")
	suffix := parts[len(parts)-1]

	compress = false

	switch suffix {
	case "css", "js", "html", "htm", "ttf", "eot", "svg", "txt", "csv":
		compress = true
	}
	mimeType = mime.TypeByExtension("." + suffix)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	return
}
