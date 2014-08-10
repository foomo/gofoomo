package handler

import (
	"github.com/foomo/gofoomo/foomo"
	"io"
	"net/http"
	"os"
	"strings"
)

// Handles serving static files from the local file system. It knows about
// foomos hierarchy and serves files from the htdocs directories of modules.
// Currently it will also serve files of disabled modules.

type StaticFiles struct {
	foomo *foomo.Foomo
}

func NewStaticFiles(foomo *foomo.Foomo) *StaticFiles {
	sf := new(StaticFiles)
	sf.foomo = foomo
	return sf
}

func (files *StaticFiles) HandlesRequest(incomingRequest *http.Request) bool {
	if strings.HasPrefix(incomingRequest.URL.Path, "/foomo/modules/") {
		parts := strings.Split(incomingRequest.URL.Path, "/")
		if len(parts) > 3 {
			moduleNameParts := strings.Split(parts[3], "-")
			return fileExists(files.foomo.GetModuleHtdocsDir(moduleNameParts[0]) + "/" + strings.Join(parts[4:], "/"))
		} else {
			return false
		}
	} else if strings.HasPrefix(incomingRequest.URL.Path, "/foomo/modulesVar/") {
		return true
	} else {
		return false
	}
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
	f, err := os.Open(moduleDir + "/" + path)
	if err != nil {
		panic(err)
	} else {
		defer f.Close()
		w.Header().Set("Content-Type", getContentType(path))
		io.Copy(w, f)
	}
}

func getContentType(path string) string {
	if strings.HasSuffix(path, ".png") {
		return "image/png"
	} else if strings.HasSuffix(path, ".jpg") {
		return "image/jpeg"
	} else if strings.HasSuffix(path, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(path, ".gif") {
		return "image/gif"
	} else if strings.HasSuffix(path, ".css") {
		return "text/css"
	} else if strings.HasSuffix(path, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(path, ".html") {
		return "text/html"
	} else if strings.HasSuffix(path, ".") {
		return ""
	} else {
		return "octet/stream"
	}
}
