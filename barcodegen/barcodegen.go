package barcodegen

import (
	"io/ioutil"
	"net/http"

	"github.com/adneg/golangmodules/logtrace"
	"github.com/julienschmidt/httprouter"

	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

var (
	rest *httprouter.Router
)

func Init(r *httprouter.Router) {
	rest = r
}

func Start() {
	createcallhandles()

}
func createcallhandles() {
	// wget   --post-data "data" http://localhost:8081/qrgenerator
	rest.POST("/qrgenerator", postStringToQrCodePNG)
	rest.POST("/sumbitqrgenerator", submitStringToQrCodePNG)
	rest.GET("/sumbitqrgenerator", getHtmlStringToQrCodePNG)

}
func postStringToQrCodePNG(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dataByte, _ := ioutil.ReadAll(r.Body)

	dataString := string(dataByte)
	logtrace.Info.Println("GENERATED QR CODE: ", dataString)
	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	png.Encode(w, qrCode)
}
func submitStringToQrCodePNG(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dataString := r.FormValue("dataString")

	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	png.Encode(w, qrCode)
}

func getHtmlStringToQrCodePNG(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	htmlcode := `<h1>QR CODE GENERATOR</h1>
<div>Please enter the string you want to QRCode.</div>
<form action="sumbitqrgenerator" method=post>
    <input type="text" name="dataString">
    <input type="submit" value="Submit">
</form>`

	w.Write([]byte(htmlcode))
}
