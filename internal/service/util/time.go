package util

import (
	"fmt"
	"time"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	cmlib "gitlab.senseauto.com/apcloud/library/common-go/lib"
)

func TimeParse(timeStr string) int64 {
	timeUnix, _ := time.ParseInLocation(constant.FullTimeTemplate, timeStr, cmlib.GetCSTLocation())
	return timeUnix.UnixMilli()
}

func TimeSince(t time.Time) int64 {
	return int64(time.Since(t))
}

func GetTimeRangeTT(start, end int64) []time.Time {
	var dates []time.Time
	startTime := time.UnixMilli(start)
	endTime := time.UnixMilli(end)
	fmt.Println("show GetTimeRangeTT init", startTime, endTime)
	year, month, day := endTime.Date()
	endTime = time.Date(year, month, day, 8, 0, 0, 0, cmlib.GetCSTLocation())
	fmt.Println("show GetTimeRangeTT after", startTime, endTime, year, month, day)
	for t := startTime; !t.After(endTime); t = t.AddDate(0, 0, 1) {
		fmt.Println("show GetTimeRangeTT for", t, endTime, t.AddDate(0, 0, 1))
		dates = append(dates, t.Add(32*time.Hour))
	}
	return dates
}

func GetYesterDateStartAndEndTime(tt time.Time) (string, string) {
	beginUnix, endUnix := GetYesterDateStartAndEndUnix(tt)
	beginTime := time.UnixMilli(beginUnix).Format(constant.FullTimeTemplate)
	endTime := time.UnixMilli(endUnix).Format(constant.FullTimeTemplate)
	return beginTime, endTime
}

func GetCurDateStartAndEndUnix(tt time.Time) (int64, int64) {
	beginStr := tt.Format(constant.DateTemplate) + " 00:00:00"
	endStr := tt.Format(constant.DateTemplate) + " 23:59:59"

	beginUnix, _ := time.ParseInLocation(constant.FullTimeTemplate, beginStr, cmlib.GetCSTLocation())
	endUnix, _ := time.ParseInLocation(constant.FullTimeTemplate, endStr, cmlib.GetCSTLocation())
	return beginUnix.UnixMilli(), endUnix.UnixMilli()
}

func GetYesterDateStartAndEndUnix(tt time.Time) (int64, int64) {
	tt = tt.Add(-time.Hour * 24)
	beginStr := tt.Format(constant.DateTemplate) + " 00:00:00"
	endStr := tt.Format(constant.DateTemplate) + " 23:59:59"

	beginUnix, _ := time.ParseInLocation(constant.FullTimeTemplate, beginStr, cmlib.GetCSTLocation())
	endUnix, _ := time.ParseInLocation(constant.FullTimeTemplate, endStr, cmlib.GetCSTLocation())
	return beginUnix.UnixMilli(), endUnix.UnixMilli()
}

func GetWeekStartAndEndUnix(tt time.Time) (int64, int64) {
	year, week := tt.ISOWeek()
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, cmlib.GetCSTLocation())
	// 移动到指定周数
	weekOffset := (week - 1) * 7
	firstDay = firstDay.AddDate(0, 0, weekOffset)

	// 计算该周的最后一天（周日）
	lastDay := firstDay.AddDate(0, 0, 7)

	// 获取 Unix 时间戳
	firstDayTimestamp := firstDay.UnixMilli()
	lastDayTimestamp := lastDay.UnixMilli()
	return firstDayTimestamp, lastDayTimestamp
}

func GetMonthStartAndEndUnix(tt time.Time) (int64, int64) {
	year := tt.Year()
	month := tt.Month()
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, cmlib.GetCSTLocation())

	nextMonth := firstDay.AddDate(0, 1, 0)
	// lastDay := nextMonth.AddDate(0, 0, -1)

	// 获取 Unix 时间戳
	firstDayTimestamp := firstDay.UnixMilli()
	lastDayTimestamp := nextMonth.UnixMilli()
	return firstDayTimestamp, lastDayTimestamp
}
