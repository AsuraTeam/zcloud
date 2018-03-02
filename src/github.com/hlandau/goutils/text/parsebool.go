package text

import (
	"regexp"
	"strings"
)

var re_noForgiving = regexp.MustCompile(`^\s*(00*|[nN][oO]?|[fF]([aA][lL][sS][eE])?|)\s*$`)

func ParseBoolForgiving(s string) (value, ok bool) {
	return !re_noForgiving.MatchString(s), true
}

func ParseBoolUser(s string) (value, ok bool) {
	return parseBoolUser(s, nil)
}

var (
	yes = true
	no  = false
)

func ParseBoolUserDefaultYes(s string) (value, ok bool) {
	return parseBoolUser(s, &yes)
}

func ParseBoolUserDefaultNo(s string) (value, ok bool) {
	return parseBoolUser(s, &no)
}

func parseBoolUser(s string, dflt *bool) (value, ok bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "y" || s == "yes" {
		return true, true
	} else if s == "n" || s == "no" {
		return false, true
	} else if s == "" && dflt != nil {
		return *dflt, true
	} else {
		return false, false
	}
}
