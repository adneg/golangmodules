package main

import (
	"os"

	"os/signal"

	"github.com/adneg/golangmodules/backuplist"
	"github.com/adneg/golangmodules/barcodegen"
	"github.com/adneg/golangmodules/errorpage"
	"github.com/adneg/golangmodules/gormdb"
	"github.com/adneg/golangmodules/logtrace"
	"github.com/adneg/golangmodules/restapi"
	"github.com/adneg/golangmodules/staticfiles"

	"net/http"
	"runtime/debug"
	"sync"

	"github.com/julienschmidt/httprouter"
)

const (
	serviceName = "EXAMPLE SERVER"
)

var (
	wg   sync.WaitGroup
	sigs = make(chan os.Signal, 1)
)

func init() {
	debug.SetMaxThreads(999999999999999999)
	logtrace.Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	logtrace.Trace.Println("LOGING ON")

}

func main() {

	//////////////////////// MODULY START
	//REST
	rest := restapi.Init("./config/api.json")
	//ERRORPAGE
	errorpage.Init(rest)
	errorpage.Start()
	//STATICFILES
	staticfiles.Init("./config/staticfiles.json", rest, 4096)
	staticfiles.Start()
	//SERVICE NAME... NO MODOL
	rest.GET("/", ShowServiceName)
	//GORM
	gormdb.Init("./config/gormdb.json")

	db := gormdb.Start()
	//BACKUP LIST
	backuplist.Init(db, rest)
	backuplist.Start()
	// backuplist.CreateAllTables()

	barcodegen.Init(rest)
	barcodegen.Start()
	//////////////////////// MODULY STOP
	srv, stopchan := restapi.Start()
	wg.Add(1)
	signal.Notify(sigs, os.Interrupt)
	go func() {

		sig := <-sigs
		db.Close()
		srv.Close()

		logtrace.Info.Print(sig, "ADMIN CLOSE SERVER")
		//wysyla:
		// close(*stopchan)

		//jaki sygnal

	}()

	go func() {
		// koniec gdy sygnal konca jest rozeslany
		<-*stopchan
		wg.Done()
	}()

	wg.Wait()

}

func ShowServiceName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte(serviceName))
}
