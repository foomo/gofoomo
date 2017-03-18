package handler

import (
	"mime"
	"net/http"
	"os"
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

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}

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
	f, err := os.Open(moduleDir + "/" + path)
	panicOnErr(err)
	defer f.Close()
	_, compress := getContentType(path)
	//w.Header().Set("Content-Type", mime)
	fileInfo, err := f.Stat()
	panicOnErr(err)
	//const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
	w.Header().Set("Expires", time.Now().Add(time.Hour*24*365).Format(http.TimeFormat))
	if compress && strings.Contains(incomingRequest.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		crw := utils.NewCompressedResponseWriter(w)
		defer crw.Close()
		w = crw
	}
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
