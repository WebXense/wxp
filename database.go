package wxp

import (
	"github.com/WebXense/env"
	"github.com/WebXense/ginger/ginger"
	"github.com/WebXense/sql/conn"
	"gorm.io/gorm"
)

var DB *gorm.DB
var MEM *gorm.DB
var DB_ENCRYPT_KEY string

func init() {
	var err error

	if env.String("GIN_MODE", true) == ginger.GIN_MODE_TEST {
		DB, err = conn.SQLite("test.db", false)
	} else {
		DB, err = conn.MySQL(
			env.String("DB_HOST"),
			env.String("DB_PORT"),
			env.String("DB_USERNAME"),
			env.String("DB_PASSWORD"),
			env.String("DB_DATABASE"),
			env.String("GIN_MODE", true) == ginger.GIN_MODE_RELEASE,
		)
	}

	if err != nil {
		panic(err)
	}

	MEM, err = conn.SQLiteInMemory(env.String("GIN_MODE", true) == ginger.GIN_MODE_RELEASE)
	if err != nil {
		panic(err)
	}

	DB_ENCRYPT_KEY = env.String("DB_ENCRYPT_KEY")
}
