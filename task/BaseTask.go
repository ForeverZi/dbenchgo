package task

import (
	"database/sql"
	"dbenchgo/driver"
)

type BaseTask struct {
	Driver   string
	Parallel int
}

func (task BaseTask) DB() *sql.DB {
	return driver.GetDriver(task.Driver)
}
