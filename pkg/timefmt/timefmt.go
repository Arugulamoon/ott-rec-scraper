package timefmt

import (
	"fmt"
	"strconv"
	"strings"
)

// 24H time is 00:00 to 23:59

type TimeFmt struct {
	Start, End string
}

func TranslateEvents(s string) []TimeFmt {
	eventTimes := make([]TimeFmt, 0)
	for _, event := range SplitEvents(SanitizeTimes(s)) {
		times := AppendAMPMToStartTime(SplitEventTimes(event))
		times24h := TimeFmt{
			Start: TranslateTimeStrTo24H(times.Start),
			End:   TranslateTimeStrTo24H(times.End),
		}
		eventTimes = append(eventTimes, times24h)
	}
	return eventTimes
}

func TranslateTimeStrTo24H(s string) string {
	ampm := s[len(s)-2:]
	s1 := s[:len(s)-2]
	hm := strings.Split(s1, ":")
	i, err := strconv.Atoi(hm[0])
	if err != nil {
		panic(err)
	}

	if ampm == "am" {
		if i == 12 {
			return fmt.Sprintf("00:%s", hm[1])
		}
		return s1
	}
	if i == 12 {
		return s1
	}
	return fmt.Sprintf("%d:%s", i+12, hm[1])
}

func AppendAMPMToStartTime(t TimeFmt) TimeFmt {
	if !strings.HasSuffix(t.Start, "m") {
		t.Start += t.End[len(t.End)-2:]
	}
	return t
}

func SplitEventTimes(s string) TimeFmt {
	times := strings.Split(s, "-")
	return TimeFmt{
		Start: TranslateTimeStrToHHMM(times[0]),
		End:   TranslateTimeStrToHHMM(times[1]),
	}
}

func TranslateTimeStrToHHMM(s string) string {
	if !strings.Contains(s, ":") {
		return fmt.Sprintf("%s:00", s)
	}
	return s
}

func SplitEvents(s string) []string {
	return strings.Split(s, ",")
}

func SanitizeName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, `
			`, " ")
	return s
}

func SanitizeTimes(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "	", "")
	s = strings.ReplaceAll(s, `
`, "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
