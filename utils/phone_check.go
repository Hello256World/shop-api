package utils

import (
	"regexp"
)

func CheckPhoneNum(phone string) bool {
	phoneRegex := `^(\+?\d{1,3}[-. ]?)?(\(?\d{3}\)?[-. ]?)?\d{3}[-. ]?\d{4}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}