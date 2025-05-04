package validator

import (
	"regexp"
)

var (
	swiftCodeRegex = regexp.MustCompile(`^[A-Z0-9]{11}$`)
	iso2Regex      = regexp.MustCompile(`^[A-Z]{2}$`)
)

func IsValidSwiftCode(code string) bool {
	return swiftCodeRegex.MatchString(code)
}

func IsValidISO2(code string) bool {
	return iso2Regex.MatchString(code)
}
