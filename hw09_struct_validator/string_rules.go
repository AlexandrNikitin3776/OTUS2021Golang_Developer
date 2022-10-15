package hw09structvalidator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	InvalidStringRule   = errors.New("string tag has invalid value")
	InvalidStringLen    = errors.New("string not equal target len")
	InvalidStringRegexp = errors.New("string doesn't match regexp")
	InvalidStringIn     = errors.New("string is not in the values list")
)

type StringRule = func(string) error
type StringRuleGetter = func(string) (StringRule, error)

type stringRule struct {
}

func ParseStringRule(rule string) (StringRule, error) {
	funcName, controlValue, found := strings.Cut(rule, ":")
	if !found {
		return nil, InvalidStringRule
	}

	sr := stringRule{}
	ruleGetter, found := sr.getRuleGetter(funcName)
	if !found {
		return nil, InvalidStringRule
	}

	checkFunc, err := ruleGetter(controlValue)
	if err != nil {
		return nil, InvalidStringRule
	}
	return checkFunc, nil
}

func (sr stringRule) getRuleGetter(funcName string) (StringRuleGetter, bool) {
	switch funcName {
	case "len":
		return sr.getLenRule, true
	case "regexp":
		return sr.getRegexpRule, true
	case "in":
		return sr.getInRule, true
	default:
		return nil, false
	}
}

func (sr stringRule) getLenRule(controlValue string) (StringRule, error) {
	controlLen, err := strconv.Atoi(controlValue)
	if err != nil {
		return nil, InvalidStringRule
	}

	return func(value string) error {
		if len(value) == controlLen {
			return nil
		}
		return InvalidStringLen
	}, nil
}

func (sr stringRule) getRegexpRule(controlValue string) (StringRule, error) {
	re, err := regexp.Compile(controlValue)
	if err != nil {
		return nil, InvalidStringRule
	}

	return func(value string) error {
		if re.MatchString(value) {
			return nil
		}
		return InvalidStringRegexp
	}, nil
}

func (sr stringRule) getInRule(controlValue string) (StringRule, error) {
	values := strings.Split(controlValue, ",")
	for _, val := range values {
		if val == "" {
			return nil, InvalidStringRule
		}
	}

	return func(value string) error {
		for _, targetValue := range values {
			if value == targetValue {
				return nil
			}
		}
		return InvalidStringIn
	}, nil
}
