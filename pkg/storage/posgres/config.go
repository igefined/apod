package posgres

type postgresConfig struct {
	host     string
	port     string
	database string
	username string
	password string
}

func NewPostgresConfig(host, port, database, username, password string) *postgresConfig {
	return &postgresConfig{
		host:     host,
		port:     port,
		username: username,
		password: password,
		database: database,
	}
}
