package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	InvalidIntMin = errors.New("int value is lower than min value")
	InvalidIntMax = errors.New("int value is greater than min value")
	InvalidIntIn  = errors.New("int value is not in the values list")
)

type IntRule = func(int) error
type IntRuleGetter = func(string) (IntRule, error)

type intRule struct {
}

func ParseIntRules(tag string) ([]IntRule, error) {
	rules := strings.Split(tag, tagDivider)
	result := make([]IntRule, len(rules))
	for _, rule := range rules {
		if checkFunc, err := ParseIntRule(rule); err != nil {
			return nil, err
		} else {
			result = append(result, checkFunc)
		}
	}
	return result, nil
}

func ParseIntRule(rule string) (IntRule, error) {
	funcName, controlValue, found := strings.Cut(rule, ":")
	if !found {
		return nil, fmt.Errorf("tag must contain %v", tagDefinder)
	}

	ir := intRule{}
	ruleGetter, found := ir.getRuleGetter(funcName)
	if !found {
		return nil, fmt.Errorf("tag rule %q isn't supported", funcName)
	}

	checkFunc, err := ruleGetter(controlValue)
	if err != nil {
		return nil, fmt.Errorf("erorr while getting rule: %w", err)
	}
	return checkFunc, nil
}

func (ir intRule) getRuleGetter(funcName string) (IntRuleGetter, bool) {
	switch funcName {
	case "min":
		return ir.getMinRule, true
	case "max":
		return ir.getMaxRule, true
	case "in":
		return ir.getInRule, true
	default:
		return nil, false
	}
}

func (ir intRule) getMinRule(controlValue string) (IntRule, error) {
	controlMin, err := strconv.Atoi(controlValue)
	if err != nil {
		return nil, err
	}

	return func(value int) error {
		if value < controlMin {
			return InvalidIntMin
		}
		return nil
	}, nil
}

func (ir intRule) getMaxRule(controlValue string) (IntRule, error) {
	controlMax, err := strconv.Atoi(controlValue)
	if err != nil {
		return nil, err
	}

	return func(value int) error {
		if value > controlMax {
			return InvalidIntMax
		}
		return nil
	}, nil
}

func intSetFromSlice(sl []string) (map[int]struct{}, error) {
	result := make(map[int]struct{}, len(sl))
	for _, value := range sl {
		if intValue, err := strconv.Atoi(value); err != nil {
			return nil, err
		} else {
			result[intValue] = struct{}{}
		}
	}
	return result, nil
}

func (ir intRule) getInRule(controlValue string) (IntRule, error) {
	controlValues, err := intSetFromSlice(strings.Split(controlValue, ","))
	if err != nil {
		return nil, err
	}

	return func(value int) error {
		if _, ok := controlValues[value]; !ok {
			return InvalidIntIn
		}
		return nil
	}, nil
}
