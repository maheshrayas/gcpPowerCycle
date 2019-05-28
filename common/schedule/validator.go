package common

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	RE_SCHEDULE_TIME = "((?i)%s_([0-9]{1,2}-[0-9]{1,2}))"
	RE_SCHEDULE_DAYS = "([a-z]{3}-[a-z]{3})"
)

var shortDayNames = []string{
	"mon",
	"tue",
	"wed",
	"thu",
	"fri",
	"sat",
	"sun",
}

type InstaceTimeDetails struct {
	CurrentTime   time.Time
	Localtimezone string //Australia/Sydney
	InsLabel      string
	InstanceName  string
}

func (I *InstaceTimeDetails) getTimeZone() {
	t := time.Now()
	fmt.Println(I.InsLabel)
	localTime, err := time.LoadLocation(I.Localtimezone)
	if err != nil {
		fmt.Println(err)
	}
	I.CurrentTime = t.In(localTime)
}

func (I *InstaceTimeDetails) Validate() bool {
	I.getTimeZone()
	label := I.InsLabel
	response := fmt.Sprintf(RE_SCHEDULE_TIME, "start")
	start := regular(label, response)
	start = strings.Replace(start, "-", ":", 1)
	response = fmt.Sprintf(RE_SCHEDULE_TIME, "stop")
	stop := regular(label, response)
	stop = strings.Replace(stop, "-", ":", 1)
	response = fmt.Sprintf(RE_SCHEDULE_DAYS)
	week := regular(label, response)
	var dayRanges = make([]int, 1)
	dayRanges[0] = int(I.CurrentTime.Weekday())
	for _, day := range strings.Split(week, "-") {
		if idx := getDayIndex(day); idx != 99 {
			dayRanges = append(dayRanges, idx)
		}
	}
	if isDayInRange := inBetween(dayRanges); !isDayInRange {
		fmt.Printf("Instance %s needs to be running between %s to %s\n", I.InstanceName, shortDayNames[dayRanges[1]], shortDayNames[dayRanges[2]])
		return false
	}
	return I.isTimeInRange(start, stop)
}

func inBetween(i []int) bool {
	if (i[0] >= i[1]) && (i[0] <= i[2]) {
		return true
	}
	return false
}

func getDayIndex(day string) int {
	for i, val := range shortDayNames {
		if val == day {
			return i
		}
	}
	return 99
}

func regular(matchstring string, reg string) string {
	re, err := regexp.Compile(reg)
	if err != nil {
		fmt.Println("Something went wrong")
	}
	if re.MatchString(matchstring) {
		return re.FindString(matchstring)
	}
	return "NoMatch"
}

func (I *InstaceTimeDetails) isTimeInRange(start string, stop string) bool {
	localTime, _ := time.LoadLocation(I.Localtimezone)
	formatedTime := I.CurrentTime.Format("15:04")
	t1, _ := time.ParseInLocation("15:04", formatedTime, localTime)
	t2, _ := time.ParseInLocation("15:04", strings.Split(start, "_")[1], localTime)
	t3, _ := time.ParseInLocation("15:04", strings.Split(stop, "_")[1], localTime)
	if t1.After(t2) && t1.Before(t3) {
		return true
	}
	return false
}
