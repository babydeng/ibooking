package internal

import (
	"Group4/ibooking-back/store"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func CheckDateFormat(dateString string) bool {
	pattern := `^\d{4}-\d{2}-\d{2}$`
	match, _ := regexp.MatchString(pattern, dateString)
	return match
}

func CheckTime(order store.IbookingOrder) bool {
	dateStr := order.Date
	hourStr := order.StartTime

	location, _ := time.LoadLocation("Asia/Shanghai")
	// 解析日期字符串
	date, err := time.ParseInLocation("2006-01-02", dateStr, location)
	if err != nil {
		fmt.Println("Invalid date string:", err)
		return false
	}

	// 解析小时
	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return false
	}

	// 获取当前时间 精确到小时
	currentTime := time.Now()
	truncateCurrentTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	fmt.Println(truncateCurrentTime)
	fmt.Println(hour)
	// 创建组合时间
	combinedTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	fmt.Println(combinedTime)
	fmt.Println(currentTime.Hour())
	// 比较时间
	if combinedTime.Before(truncateCurrentTime) {
		return false
	} else if combinedTime.After(truncateCurrentTime) {
		return true
	} else {
		if hour-currentTime.Hour() < 2 {
			return false
		} else {
			return true
		}
	}

}
