package utils

import regexp "regexp"

func IsValidEmail(s string) bool {
	matchString, err := regexp.MatchString(`^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$`, s)
	if err != nil {
		return false
	}
	return matchString
}
