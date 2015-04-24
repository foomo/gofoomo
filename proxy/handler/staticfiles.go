package handler

import (
	"github.com/foomo/gofoomo/foomo"
	"github.com/foomo/gofoomo/proxy/utils"
	"net/http"
	"os"
	"strings"
	"time"
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
			if strings.HasSuffix(parts[len(parts)-1], ".php") {
				return false
			} else {
				return fileExists(files.foomo.GetModuleHtdocsDir(moduleNameParts[0]) + "/" + strings.Join(parts[4:], "/"))
			}
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
	/*	if compress {
			err := utils.ServeCompressed(w, incomingRequest, func(writer io.Writer) error {
				_, err := io.Copy(writer, f)
				return err
			})
			panicOnErr(err)
		} else {
			_, err := io.Copy(w, f)
			panicOnErr(err)
		}
	*/

}

func getContentType(path string) (string, bool) {
	if strings.HasSuffix(path, ".png") {
		return "image/png", false
	} else if strings.HasSuffix(path, ".jpg") {
		return "image/jpeg", false
	} else if strings.HasSuffix(path, ".jpeg") {
		return "image/jpeg", false
	} else if strings.HasSuffix(path, ".gif") {
		return "image/gif", false
	} else if strings.HasSuffix(path, ".css") {
		return "text/css", true
	} else if strings.HasSuffix(path, ".js") {
		return "application/javascript", true
	} else if strings.HasSuffix(path, ".html") {
		return "text/html", true
	} else if strings.HasSuffix(path, ".") {
		return "", false
	} else {
		return "octet/stream", false
	}
}
