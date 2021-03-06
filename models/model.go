package models

import (
	"database/sql/driver"
	"easy-go-web/pkg/global"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

// 由于gorm提供的base model没有json tag, 使用自定义
type Model struct {
	Id        uint           `gorm:"primary_key;comment:'自增编号'" json:"id"`
	CreatedAt LocalTime      `gorm:"comment:'创建时间'" json:"createdAt"`
	UpdatedAt LocalTime      `gorm:"comment:'更新时间'" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"comment:'删除时间(软删除)'" sql:"index" json:"deletedAt"`
}

// 表名设置
func (Model) TableName(name string) string {
	// 添加表前缀
	return fmt.Sprintf("%s_%s", global.Conf.Mysql.TablePrefix, name)
}

// 本地时间
type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	str := strings.Trim(string(data), "\"")
	// ""空值不进行解析
	// 避免环包调用, 不再调用utils
	if str == "null" || strings.TrimSpace(str) == "" {
		*t = LocalTime{Time: time.Time{}}
		return
	}

	// 设置str
	t.SetString(str)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	s := t.Format(global.SecLocalTimeFormat)
	// 处理时间0值
	if t.IsZero() {
		s = ""
	}
	output := fmt.Sprintf("\"%s\"", s)
	return []byte(output), nil
}

// gorm 写入 mysql 时调用
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// gorm 检出 mysql 时调用
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to LocalTime", v)
}

// 用于 fmt.Println 和后续验证场景
func (t LocalTime) String() string {
	return t.Format(global.SecLocalTimeFormat)
}

// 只需要日期
func (t LocalTime) DateString() string {
	return t.Format(global.DateLocalTimeFormat)
}

// 设置字符串
func (t *LocalTime) SetString(str string) *LocalTime {
	if t != nil {
		// 指定解析的格式(设置转为本地格式)
		now, err := time.ParseInLocation(global.SecLocalTimeFormat, str, time.Local)
		if err == nil {
			*t = LocalTime{Time: now}
		}
	}
	return t
}
