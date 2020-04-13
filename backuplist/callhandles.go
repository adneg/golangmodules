package backuplist

import (
	"net/http"
	"strconv"

	"github.com/adneg/golangmodules/logtrace"
	"github.com/julienschmidt/httprouter"
)

const (
	fileContentType = "application/json"
)

func Createcallhandles() {
	// rest.GET("/backuplist/backups/:limit", getLimitBackups)
	rest.GET("/backuplist/firms", getAllFirms)
	rest.GET("/backuplist/pawilons", getAllPawilons)
	rest.GET("/backuplist/databases", getAllDatabases)
	rest.GET("/backuplist/servers", getAllServers)
	rest.GET("/backuplist/backups", getAllBackups)

	// rest.POST("/backuplist/backup", postBackup)
	// rest.GET("/backuplist/firma/:id", get)
	// rest.GET("/backuplist/firms", getFirms)
	// rest.GET("/backuplist/firma/:id", getFirma)
	// rest.GET("/backuplist/firms", getFirms)

	// REST.GET("/", BasicAuth(Protected))
	// REST.POST("/zaloguj", BasicAuth(Logowanie))
	// REST.GET("/wymiary", BasicAuth(Wymiary))
	// REST.DELETE("/wymiar/:id", BasicAuth(RmWymiar))

}

func getAllFirms(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	l := selectFirmaAll()
	logtrace.Info.Println("GET ALL FIRMS LIMIT")
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Add("Content-Length", strconv.Itoa(len(l)))
	w.Write(l)
}
func getAllPawilons(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	l := selectPawilonAll()
	logtrace.Info.Println("GET ALL PAWILONS LIMIT")
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Add("Content-Length", strconv.Itoa(len(l)))
	w.Write(l)
}

func getAllServers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	l := selectServerAll()
	logtrace.Info.Println("GET ALL SERVERS LIMIT")
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Add("Content-Length", strconv.Itoa(len(l)))
	w.Write(l)

}

func getAllDatabases(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	l := selectBazaDanychAll()
	logtrace.Info.Println("GET ALL DATABASES LIMIT")
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Add("Content-Length", strconv.Itoa(len(l)))
	w.Write(l)

}

func getAllBackups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	l := selectKopiaAll()
	logtrace.Info.Println("GET ALLBACKUPS LIMIT")
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Add("Content-Length", strconv.Itoa(len(l)))
	w.Write(l)

}

func getLimitBackups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	i := ps.ByName("limit")
	limit, err := strconv.Atoi(i)
	if err != nil {
		logtrace.Error.Println(err.Error())
		rest.NotFound.ServeHTTP(w, r)
		return
	}

	l := selectKopiaLimit(limit)
	logtrace.Info.Println("GET LIST BACKUPS LIMIT: ", i)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Add("Content-Length", strconv.Itoa(len(l)))
	w.Write(l)

}
