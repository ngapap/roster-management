package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

const Driver = "postgres"

func newDBStringFromConfig(config *viper.Viper) string {
	var dbParams []string
	dbParams = append(dbParams, fmt.Sprintf("user=%s", config.Get("POSTGRES_USER")))
	dbParams = append(dbParams, fmt.Sprintf("host=%s", config.Get("POSTGRES_HOST")))
	dbParams = append(dbParams, fmt.Sprintf("port=%s", config.Get("POSTGRES_PORT")))
	dbParams = append(dbParams, fmt.Sprintf("dbname=%s", config.Get("POSTGRES_DB")))
	if password := config.Get("POSTGRES_PASSWORD"); password != "" {
		dbParams = append(dbParams, fmt.Sprintf("password=%s", password))
	}
	dbParams = append(dbParams, fmt.Sprintf("sslmode=%s", config.Get("POSTGRES_SSL_MODE")))
	dbParams = append(dbParams, fmt.Sprintf("connect_timeout=%d", config.GetInt("POSTGRES_CONNECT_TIMEOUT")))
	dbParams = append(dbParams, fmt.Sprintf("statement_timeout=%d", config.GetInt("POSTGRES_STATEMENT_TIMEOUT")))
	dbParams = append(dbParams, fmt.Sprintf("idle_in_transaction_session_timeout=%d", config.GetInt("POSTGRES_IDLE_IN_TIME_COUNT")))

	return strings.Join(dbParams, " ")
}

func Open(config *viper.Viper) (*sql.DB, error) {
	dbString := newDBStringFromConfig(config)
	logrus.Println("dbstring", dbString)
	db, err := otelsql.Open(Driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Connectx(config *viper.Viper) (*sqlx.DB, error) {
	dbString := newDBStringFromConfig(config)

	db, err := otelsqlx.Connect(Driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
