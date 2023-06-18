package mysql

import (
	"database/sql"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

var Pool *sql.DB

func init() {
	if !viper.GetBool("drivers.mysql.enabled") {
		return
	}
	var err error
	Pool, err = sql.Open("mysql", viper.GetString("drivers.mysql.dsn"))
	if err != nil {
		panic(err)
	}
	Pool.SetMaxOpenConns(viper.GetInt("drivers.mysql.max_connections"))
	Pool.SetMaxIdleConns(viper.GetInt("drivers.mysql.max_connections"))
}
