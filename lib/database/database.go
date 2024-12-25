package database

import (
	"gorm.io/gorm"

	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
)

type DatabaseInterface interface {
	Connection() *gorm.DB
}

type connection struct {
	pg PostgresClientInterface
}

func StartDatabaseClient(env constants.Env) DatabaseInterface {
	return &connection{
		pg: NewPostgresClient(env),
	}
}

func (conn connection) Connection() *gorm.DB {
	return conn.pg.Connection()
}
