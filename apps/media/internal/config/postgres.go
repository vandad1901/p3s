package config

import (
	"fmt"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
	gnrconfig "github.com/vandad1901/p3s/packages/go/envutil"
)

func getDSN() string {
	host := gnrconfig.MustGetString("UPLOAD_PG_HOST")
	port := gnrconfig.MustGetString("UPLOAD_PG_PORT")
	user := gnrconfig.MustGetString("UPLOAD_PG_USER")
	password := gnrconfig.MustGetString("UPLOAD_PG_PASSWORD")
	dbName := gnrconfig.MustGetString("UPLOAD_PG_DATABASE")

	return fmt.Sprintf(dbpattern.PostgresDSNFormat, host, user, password, dbName, port)
}
