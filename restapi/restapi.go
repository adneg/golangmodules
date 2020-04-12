package restapi

import (
	"github.com/adneg/golangmodules/logtrace"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// func init() {

// }
var (
	conf     Configuration
	rest     *httprouter.Router
	srv      *http.Server
	stopchan = make(chan struct{})
)

// *http.Server,
//stopchanP *chan struct{}
func Init(configfile string) (rest *httprouter.Router) {

	conf = loadconfig(configfile)

	// loging.Info.Println(a.Adres)

	rest = httprouter.New()

	if conf.Checkcert {
		clientCA, err := ioutil.ReadFile(conf.Ca)
		if err != nil {
			logtrace.Error.Fatalln(err.Error())
			close(stopchan)
		}
		clientCAPool := x509.NewCertPool()
		clientCAPool.AppendCertsFromPEM(clientCA)
		logtrace.Trace.Println(conf.Ca, "LOADED")
		// log.Println("ClientCA loaded")
		srv = &http.Server{
			Addr: conf.Port, Handler: rest,
			TLSConfig: &tls.Config{
				ClientCAs:  clientCAPool,
				ClientAuth: tls.RequireAndVerifyClientCert,
				GetCertificate: func(info *tls.ClientHelloInfo) (certificate *tls.Certificate, e error) {
					c, err := tls.LoadX509KeyPair(conf.Crt, conf.Key)
					if err != nil {
						logtrace.Error.Fatalln(err.Error())
						close(stopchan)

					}
					logtrace.Trace.Println(conf.Crt, conf.Key, "LOADED")
					return &c, nil
				},
			},
		}

	} else {
		srv = &http.Server{Addr: conf.Port, Handler: rest}
	}
	return rest
}
func Start() (server *http.Server, stopchanP *chan struct{}) {

	go func() {

		var err error
		logtrace.Trace.Println("SERVER ON")
		if conf.Ssl {

			err = srv.ListenAndServeTLS(conf.Crt, conf.Key)

		} else {
			err = srv.ListenAndServe()
		}
		if err != nil {
			logtrace.Error.Println("SERWER OFF ", err.Error())
		}

		close(stopchan)

	}()
	return srv, &stopchan
}
