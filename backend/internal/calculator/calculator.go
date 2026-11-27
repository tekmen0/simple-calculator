package calculator

import (
	"errors"
	"math"
)

var (
	ErrUnknownOperation = errors.New("unknown operation")
	ErrInvalidOperand   = errors.New("operand must be finite")
	ErrDivisionByZero   = errors.New("division by zero")
	ErrResultOutOfRange = errors.New("result is out of range")
)

func Supports(name string) bool {
	return lookup(name) != nil
}

func Calculate(name string, left, right float64) (float64, error) {
	apply := lookup(name)
	if apply == nil {
		return 0, ErrUnknownOperation
	}
	if !isFinite(left) || !isFinite(right) {
		return 0, ErrInvalidOperand
	}

	result, err := apply(left, right)
	if err != nil {
		return 0, err
	}
	if !isFinite(result) {
		return 0, ErrResultOutOfRange
	}
	return result, nil
}

func lookup(name string) func(float64, float64) (float64, error) {
	switch name {
	case "add":
		return add
	case "subtract":
		return subtract
	case "multiply":
		return multiply
	case "divide":
		return divide
	default:
		return nil
	}
}

func add(left, right float64) (float64, error) {
	return left + right, nil
}

func subtract(left, right float64) (float64, error) {
	return left - right, nil
}

func multiply(left, right float64) (float64, error) {
	return left * right, nil
}

func divide(left, right float64) (float64, error) {
	if right == 0 {
		return 0, ErrDivisionByZero
	}
	return left / right, nil
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
