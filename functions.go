package template

import (
	"fmt"
	"time"
)

func RMB(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

func ChineseDate(t time.Time) string {
	return t.Format("2006 年 01 月 02 日")
}
