package time

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// LocalTime 是自定义的时间类型，基于 time.Time
type LocalTime time.Time

// MarshalJSON 实现了 JSON 序列化接口
func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

// UnmarshalJSON 实现了 JSON 反序列化接口
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	// 去掉 JSON 字符串中的引号
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// 解析字符串为 time.Time
	parsedTime, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}

	// 将解析结果赋值给 LocalTime
	*t = LocalTime(parsedTime)
	return nil
}

// Now 返回当前时间的 LocalTime 类型
func Now() LocalTime {
	return LocalTime(time.Now())
}

// Value 实现了 driver.Valuer 接口，返回数据库中的值
func (t LocalTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现了 sql.Scanner 接口，用于从数据库中扫描值
func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = LocalTime(v)
	case []byte:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*t = LocalTime(parsedTime)
	case string:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*t = LocalTime(parsedTime)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}
