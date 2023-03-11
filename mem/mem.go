package mem

import (
	"github.com/WebXense/env"
	"github.com/WebXense/ginger/ginger"
	"github.com/WebXense/sql/conn"
	"gorm.io/gorm"
)

var MEM *gorm.DB

func init() {
	var err error
	MEM, err = conn.SQLiteInMemory(env.String("GIN_MODE", true) == ginger.GIN_MODE_RELEASE)
	if err != nil {
		panic(err)
	}
}
