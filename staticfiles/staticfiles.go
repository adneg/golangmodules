package staticfiles

import (
	//wymaga errorpage, logtrace
	"io"

	"github.com/adneg/golangmodules/logtrace"

	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var (
	conf       Configuration
	rest       *httprouter.Router
	buffersize int64
)

func Init(configfile string, r *httprouter.Router, b int64) {
	conf = loadconfig(configfile)
	rest = r
	buffersize = b
}

func StartShowDir(place string) {
	// 1) pokazujesz pliki
	// rest.ServeFiles(url+"*filepath", http.Dir("./"))

}

func Start() {
	if conf.DirectoryIndex {
		rest.ServeFiles(conf.Url+"*filepath", http.Dir(conf.Place))
	} else {
		rest.GET(conf.Url+"*file", StaticNoSessionFiles)
	}

	//
}

func StaticNoSessionFiles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	pathstring := strings.Replace(r.URL.Path[:], conf.Url, conf.Place, 1)
	logtrace.Info.Println("START SEND FILE", pathstring)
	sourceFileStat, err := os.Stat(pathstring)
	if err != nil {
		logtrace.Warning.Println("FILE NOT EXIST", pathstring)
		rest.NotFound.ServeHTTP(w, r)
		return
	}

	if !sourceFileStat.Mode().IsRegular() {
		logtrace.Warning.Println("FILE NOT REGULAR", pathstring)
		//..notfoung
		rest.NotFound.ServeHTTP(w, r)
		return

	}

	// zapis z buforem jakby plik to był film
	//  bo ram
	source, err := os.Open(pathstring)
	if err != nil {
		//..notfoung
		logtrace.Warning.Println("NOT ACCESS TO FILE", pathstring)
		rest.NotFound.ServeHTTP(w, r)
		return
	}
	defer source.Close()
	fileStat, _ := source.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)
	fileContentType := "text/plain"
	if strings.HasSuffix(pathstring, ".css") {
		fileContentType = "text/css"
	} else if strings.HasSuffix(pathstring, ".html") {
		fileContentType = "text/html"
	} else if strings.HasSuffix(pathstring, ".js") {
		fileContentType = "application/javascript"
	} else if strings.HasSuffix(pathstring, ".png") {
		fileContentType = "image/png"
	} else if strings.HasSuffix(pathstring, ".svg") {
		fileContentType = "image/svg+xml"
	} else if strings.HasSuffix(pathstring, ".json") {
		fileContentType = "application/json"
	} else {
		fileHeader := make([]byte, 512)
		//Copy the headers into the FileHeader buffer
		source.Read(fileHeader)
		//Get content type of file
		fileContentType = http.DetectContentType(fileHeader)
		source.Seek(0, 0)

	}

	w.Header().Set("Content-Type", fileContentType)
	w.Header().Add("Content-Length", fileSize)
	buf := make([]byte, buffersize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			//"EOF"
			logtrace.Info.Println("FINISH SEND FIND EOF", pathstring)
			return
		}
		if n == 0 {

			logtrace.Info.Println("FINISH SEND 0 BYTES TO SEND", pathstring)
			//"zerobitow StaticNoSessionFiles"
			return
		}

		if _, err := w.Write(buf[:n]); err != nil {
			//"rozlaczyl sie huj"
			//fmt.Println(err.Error())
			logtrace.Warning.Println("FINISH SEND  CLIENT DISCONECT", pathstring)
			return
		}
	}
	//..notfoung
	logtrace.Warning.Println("ERROR FILE SEND", pathstring)
	rest.NotFound.ServeHTTP(w, r)

}
