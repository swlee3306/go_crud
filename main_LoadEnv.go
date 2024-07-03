package main

import (
	"os"

	"baton-om-data-apiservice/internal/sysenv"

	"bitbucket.org/okestrolab/baton-om-sdk/btoutil"
)

func main_LoadEnvDb() error {

	// database
	{
		if val, ok := os.LookupEnv("BATON_DB_DSN"); ok && (len(val) > 0) {
			sysenv.Database.Dsn = val
		}
		if val, ok := os.LookupEnv("BATON_DATABASE_MAX_IDLE_CONNS"); ok && (len(val) > 0) {
			sysenv.Database.MaxIdleConns = btoutil.ToInt(val)
		}
		if val, ok := os.LookupEnv("BATON_DATABASE_MAX_LIFETIME_HOUR"); ok && (len(val) > 0) {
			sysenv.Database.MaxLifetimeHour = btoutil.ToInt(val)
		}
	}

	return nil
}
