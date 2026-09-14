package config

import (
	"fmt"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
	gnrconfig "github.com/vandad1901/p3s/packages/go/envutil"
)

func getDSN() string {
	host := gnrconfig.MustGetString("API_PG_HOST")
	port := gnrconfig.MustGetString("API_PG_PORT")
	user := gnrconfig.MustGetString("API_PG_USER")
	password := gnrconfig.MustGetString("API_PG_PASSWORD")
	dbName := gnrconfig.MustGetString("API_PG_DATABASE")

	return fmt.Sprintf(dbpattern.PostgresDSNFormat, host, user, password, dbName, port)
}
