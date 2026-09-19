package data

type ConfigData struct {
	ServerAddress  string
	Context        string
	AuthServerUrl  string
	TokenTtlSecond int
	ShortUrlPrefix string
	Db             *DbConfig
}

type DbConfig struct {
	Machine  string
	Port     int
	Database string
	Login    string
	Password string
}
