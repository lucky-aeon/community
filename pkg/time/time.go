package time

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}
func Now() LocalTime {
	localTime := LocalTime(time.Now())
	return localTime
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}
