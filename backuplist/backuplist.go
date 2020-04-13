package backuplist

import (
	"encoding/json"
	"time"

	"github.com/adneg/golangmodules/logtrace"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *gorm.DB
	rest *httprouter.Router
)

type Firma struct {
	IdFirma   int    `gorm:"primary_key;AUTO_INCREMENT:index" json:"id_firma"`
	NameFirma string `gorm:"type:VARCHAR(3000);not null;unique" json:"name_firma"`
}

func (Firma) TableName() string {
	return "FIRMY"
}

type Pawilon struct {
	IdPawilon   int    `gorm:"primary_key;AUTO_INCREMENT:index" json:"id_pawilon"`
	IdFirma     int    `gorm:"type:bigint REFERENCES FIRMY(id_firma);unique_index:idx_firma_pawilon" json:"id_firmy"`
	NamePawilon string `gorm:"type:VARCHAR(3000);not null;unique_index:idx_firma_pawilon" json:"name_pawilon"`
}

func (Pawilon) TableName() string {
	return "PAWILONY"
}

type Server struct {
	IdServer   int    `gorm:"primary_key;AUTO_INCREMENT:index" json:"id_server"`
	IdPawilon  int    `gorm:"type:bigint REFERENCES PAWILONY(id_pawilon);unique_index:idx_pawilon_server" json:"id_pawilon"`
	NameServer string `gorm:"type:VARCHAR(3000);not null;unique_index:idx_pawilon_server" json:"name_server"`
}

func (Server) TableName() string {
	return "SERVERY"
}

type BazaDanych struct {
	IdDb     int    `gorm:"primary_key;AUTO_INCREMENT:index" json:"id_db"`
	IdServer int    `gorm:"type:bigint REFERENCES SERVERY(id_server);unique_index:idx_server_db" json:"id_server"`
	NameDb   string `gorm:"type:VARCHAR(3000);not null;unique_index:idx_server_db" json:"name_db"`
}

func (BazaDanych) TableName() string {
	return "BAZYDANYCH"
}

type Kopia struct {
	IdBackup  int       `gorm:"primary_key;AUTO_INCREMENT:index" json:"id_backup"`
	IdFirma   int       `gorm:"type:bigint REFERENCES FIRMY(id_firma)" json:"id_firma"`
	IdPawilon int       `gorm:"type:bigint REFERENCES PAWILONY(id_pawilon)" json:"id_pawilon"`
	IdServer  int       `gorm:"type:bigint REFERENCES SERVERY(id_server)" json:"id_server"`
	IdDb      int       `gorm:"type:bigint REFERENCES BAZYDANYCH(id_db)" json:"id_db"`
	Status    bool      `gorm:"type:BOOLEAN;NOT NULL;default:false"`
	Data      time.Time `gorm:"type:DATETIME;not null;default:CURRENT_TIMESTAMP" json:"data"`
}

type NewRecord struct {
	NameFirma, NamePawilon, NameServer, NameDb string
}

func (Kopia) TableName() string {
	return "KOPIE"
}

func Init(dbcon *gorm.DB, r *httprouter.Router) {
	logtrace.Info.Println("BACKUPLIST INIT")
	rest = r
	db = dbcon
}
func Start() {
	DropAllTable()
	CreateAllTables()
	Createcallhandles()
}

func ClerTableDB() {
	DropAllTable()
	CreateAllTables()

}

func DropAllTable() {

	db.DropTable(&Kopia{}, &BazaDanych{}, &Server{}, &Pawilon{}, &Firma{})

}

func CreateAllTables() {

	db.AutoMigrate(&Firma{}, &Pawilon{}, &Server{}, &BazaDanych{}, &Kopia{})

	n := NewRecord{NameFirma: "nazwa_firmy", NamePawilon: "nazwa_pawilonu", NameServer: "nazwa_serwera", NameDb: "nazwa_bazy"}
	for i := 0; i < 2; i++ {
		AddKopia(n)
	}
	// GetKopiaLimit(30)

}

func selectFirmaAll() (data []byte) {
	k := []Firma{}
	db.Find(&k)

	data, _ = json.Marshal(k)
	return
	// logtrace.Info.Println(string(data))
}

func selectPawilonAll() (data []byte) {
	k := []Pawilon{}
	db.Find(&k)

	data, _ = json.Marshal(k)
	return
	// logtrace.Info.Println(string(data))
}

func selectServerAll() (data []byte) {
	k := []Server{}
	db.Find(&k)

	data, _ = json.Marshal(k)
	return
	// logtrace.Info.Println(string(data))
}

func selectBazaDanychAll() (data []byte) {
	k := []BazaDanych{}
	db.Find(&k)

	data, _ = json.Marshal(k)
	return
	// logtrace.Info.Println(string(data))
}
func selectKopiaAll() (data []byte) {
	k := []Kopia{}
	db.Find(&k)

	data, _ = json.Marshal(k)
	return
	// logtrace.Info.Println(string(data))
}
func selectKopiaLimit(n int) (data []byte) {
	k := []Kopia{}
	db.Limit(n).Find(&k)

	data, _ = json.Marshal(k)
	return
	// logtrace.Info.Println(string(data))
}

func AddKopia(n NewRecord) {
	// 1) "nazwa firmy"
	// 2) "nazwa pawilonu"
	// 3) "nazwa servera"
	// 4) "NAZWA BAZY DANYCH"
	// 5) "kopia bazy danyhc"

	//1)
	f := Firma{NameFirma: n.NameFirma}
	db.Save(&f)
	logtrace.Info.Println(f.IdFirma)
	if f.IdFirma == 0 {

		db.Where(&f).First(&f)
		logtrace.Info.Println(f.IdFirma)
	}

	//2)
	p := Pawilon{IdFirma: f.IdFirma, NamePawilon: n.NamePawilon}
	db.Save(&p)
	if p.IdPawilon == 0 {

		db.Where(&p).First(&p)
		logtrace.Info.Println(p.IdPawilon)
	}

	// 3)
	s := Server{IdPawilon: p.IdPawilon, NameServer: n.NameServer}
	db.Save(&s)
	logtrace.Info.Println(s.IdServer)
	if s.IdServer == 0 {
		db.Where(&s).First(&s)
		logtrace.Info.Println(s.IdServer)
	}
	// 4)
	d := BazaDanych{IdServer: s.IdServer, NameDb: n.NameDb}
	db.Save(&d)
	logtrace.Info.Println(d.IdDb)
	if d.IdDb == 0 {
		db.Where(&d).First(&d)

	}
	// 5)
	b := Kopia{IdDb: d.IdDb, IdFirma: f.IdFirma, IdPawilon: p.IdPawilon, IdServer: s.IdServer, Status: true}
	db.Save(&b)
	logtrace.Info.Println(b.IdBackup)
	if b.IdBackup == 0 {

		logtrace.Info.Println(b.IdBackup)
	}

}
