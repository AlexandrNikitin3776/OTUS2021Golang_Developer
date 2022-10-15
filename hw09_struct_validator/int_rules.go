package hw09structvalidator

import (
	"errors"
	"strconv"
	"strings"
)

var (
	InvalidIntRule = errors.New("int tag has invalid value")
	InvalidIntMin  = errors.New("int value is lower than min value")
	InvalidIntMax  = errors.New("int value is greater than min value")
	InvalidIntIn   = errors.New("int value is not in the values list")
)

type IntRule = func(int) error
type IntRuleGetter = func(string) (IntRule, error)

type intRule struct {
}

func ParseIntRule(rule string) (IntRule, error) {
	funcName, controlValue, found := strings.Cut(rule, ":")
	if !found {
		return nil, InvalidIntRule
	}

	ir := intRule{}
	ruleGetter, found := ir.getRuleGetter(funcName)
	if !found {
		return nil, InvalidIntRule
	}

	checkFunc, err := ruleGetter(controlValue)
	if err != nil {
		return nil, InvalidIntRule
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
		return nil, InvalidIntRule
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
		return nil, InvalidIntRule
	}

	return func(value int) error {
		if value > controlMax {
			return InvalidIntMax
		}
		return nil
	}, nil
}

func (ir intRule) getInRule(controlValue string) (IntRule, error) {
	var err error
	listValues := strings.Split(controlValue, ",")
	values := make([]int, len(listValues))
	for i := range listValues {
		if listValues[i] == "" {
			return nil, InvalidIntRule
		}
		values[i], err = strconv.Atoi(listValues[i])
		if err != nil {
			return nil, InvalidIntRule
		}
	}

	return func(value int) error {
		for _, targetValue := range values {
			if value == targetValue {
				return nil
			}
		}
		return InvalidIntIn
	}, nil
}
