package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	dateStr := "2023-05-19"
	hourStr := "14"

	location, _ := time.LoadLocation("Asia/Shanghai")
	// 解析日期字符串
	date, err := time.ParseInLocation("2006-01-02", dateStr, location)
	if err != nil {
		fmt.Println("Invalid date string:", err)
		return
	}

	// 解析小时
	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		fmt.Println("Invalid hour string:", err)
		return
	}

	// 获取当前时间 精确到小时
	currentTime := time.Now().Truncate(time.Hour)
	fmt.Println(currentTime)
	// 创建组合时间
	combinedTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
	fmt.Println(combinedTime)
	// 比较时间
	if combinedTime.Before(currentTime) {
		fmt.Println("Combined time is before current time")
	} else if combinedTime.After(currentTime) {
		fmt.Println("Combined time is after current time")
	} else {
		fmt.Println("Combined time is equal to current time")
	}
}
