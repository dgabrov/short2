package data

type ConfigData struct {
	ServerAddress  string
	Context        string
	AuthServerUrl  string
	TokenTtlSecond int
	Db             *DbConfig
}

type DbConfig struct {
	Machine  string
	Port     int
	Database string
	Login    string
	Password string
}
