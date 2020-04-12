package restapi

type Configuration struct {
	Port, Adres    string
	Ssl, Checkcert bool
	Ca, Crt, Key   string
}
