package comm

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Schedule struct {
	Minute []int
	Hour   []int
	Dom    []int
	Month  []int
	Dow    []int
}

var regtime *regexp.Regexp = regexp.MustCompile(`^((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)$`)
var respace *regexp.Regexp = regexp.MustCompile(`\s+`)
var restar *regexp.Regexp = regexp.MustCompile(`\*+`)

func (_this *Schedule) Parse(timeStr string) (bool, error) {
	if !regtime.MatchString(timeStr) {
		return false, errors.New("Time error")
	}

	r1 := restar.ReplaceAllString(timeStr, "*")
	r2 := respace.ReplaceAllString(r1, " ")
	r3 := strings.SplitN(r2, " ", -1)

	if len(r3) != 5 {
		return false, errors.New("Time error")
	}

	_this.Minute = ParseNumber(&r3[0], 0, 59)
	_this.Hour = ParseNumber(&r3[1], 0, 23)
	_this.Dom = ParseNumber(&r3[2], 1, 31)
	_this.Month = ParseNumber(&r3[3], 1, 12)
	_this.Dow = ParseNumber(&r3[4], 0, 6)
	return true, nil
}

func ParseNumber(s *string, min, max int) []int {
	if *s == "*" {
		return []int{-1}
	}
	v := strings.SplitN(*s, ",", -1)
	result := make([]int, 0)
	for _, vv := range v {
		if vv == "" {
			continue
		}
		vvv := strings.SplitN(vv, "/", -1)
		var step int
		if len(vvv) < 2 || vvv[1] == "" {
			step = 1
		} else {
			step, _ = strconv.Atoi(vvv[1])
		}
		vvvv := strings.SplitN(vvv[0], "-", -1)
		var _min, _max int
		if len(vvvv) == 2 {
			_min, _ = strconv.Atoi(vvvv[0])
			_max, _ = strconv.Atoi(vvvv[1])
		} else {
			if vvv[0] == "*" {
				_min = min
				_max = max
			} else {
				_min, _ = strconv.Atoi(vvv[0])
				_max = _min
			}
		}

		for i := _min; i <= _max; i += step {
			result = append(result, i)
		}
	}
	return result
}
