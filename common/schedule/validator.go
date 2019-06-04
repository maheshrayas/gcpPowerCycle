package schedule

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

// make this as a const
var shortDayNames = []string{
	"mon",
	"tue",
	"wed",
	"thu",
	"fri",
	"sat",
	"sun",
}

var DaysMap = map[string]int{
	"Monday":    0,
	"Tuesday":   1,
	"Wednesday": 2,
	"Thursday":  3,
	"Friday":    4,
	"Saturday":  5,
	"Sunday":    6,
}

type InstaceTimeDetails struct {
	CurrentTime   time.Time
	Localtimezone string //Australia/Sydney
	InsLabel      string
	InstanceName  string
}

// set the time based on the local/set timezone
func (I *InstaceTimeDetails) getTimeZone() {
	t := time.Now()
	localTime, err := time.LoadLocation(I.Localtimezone)
	if err != nil {
		fmt.Println(err)
	}
	I.CurrentTime = t.In(localTime)
}

//Validate : Validates whether the instace/nodepool must be up or stopped
func (I *InstaceTimeDetails) Validate() bool {
	I.getTimeZone()
	label := I.InsLabel
	// for the entered label RE to format entered in Label
	response := fmt.Sprintf(RE_SCHEDULE_TIME, "start")
	start := regular(label, response)
	start = strings.Replace(start, "-", ":", 1)
	response = fmt.Sprintf(RE_SCHEDULE_TIME, "stop")
	stop := regular(label, response)
	stop = strings.Replace(stop, "-", ":", 1)
	week := regular(label, RE_SCHEDULE_DAYS)
	var dayRanges = make([]int, 1)
	// get the current day of week
	dayRanges[0] = DaysMap[I.CurrentTime.Weekday().String()]

	// get the up time days of the week
	// dayRanges = [currentdayofweek, startDay, endDay]
	for _, day := range strings.Split(week, "-") {
		if idx := getDayIndex(day); idx != 99 {
			dayRanges = append(dayRanges, idx)
		}
	}

	// check if day is the range of weeks
	if isDayInRange := inBetween(dayRanges); !isDayInRange {
		fmt.Printf("Instance %s needs to be running between %s to %s\n", I.InstanceName, shortDayNames[dayRanges[1]], shortDayNames[dayRanges[2]])
		return false
	}
	// check if current time is within the described range
	return I.isTimeInRange(start, stop)
}

func inBetween(i []int) bool {
	fmt.Println(i)
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
