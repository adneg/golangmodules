package gormdb

type Configuration struct {
	DatabaseDebug, DB          bool
	DriverDB, Databasefile     string
	MaxIdleConns, MaxOpenConns int
}
