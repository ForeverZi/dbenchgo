package driver

import (
	"database/sql"
	"dbenchgo/driver/mysql"
)

func GetDriver(name string) (pool *sql.DB) {
	switch name {
	case "mysql":
		return mysql.Pool
	default:
		return nil
	}
}
