package calculator_test

import (
	"errors"
	"math"
	"testing"

	"github.com/tekmen0/simple-calculator/backend/internal/calculator"
)

func TestSupports(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "add", want: true},
		{name: "subtract", want: true},
		{name: "multiply", want: true},
		{name: "divide", want: true},
		{name: "power", want: false},
		{name: "Add", want: false},
	}
	for _, tt := range tests {
		if got := calculator.Supports(tt.name); got != tt.want {
			t.Errorf("Supports(%q) = %t, want %t", tt.name, got, tt.want)
		}
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		left, right float64
		want        float64
	}{
		{name: "addition", operation: "add", left: -2, right: 3.5, want: 1.5},
		{name: "floating point rounding", operation: "add", left: 0.1, right: 0.2, want: 0.30000000000000004},
		{name: "subtraction order", operation: "subtract", left: 2, right: 7, want: -5},
		{name: "multiplication", operation: "multiply", left: -1.5, right: 3, want: -4.5},
		{name: "division order", operation: "divide", left: 7, right: 2, want: 3.5},
		{name: "zero dividend", operation: "divide", left: 0, right: 2, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Calculate(tt.operation, tt.left, tt.right)
			if err != nil || got != tt.want {
				t.Errorf("Calculate(%q, %g, %g) = (%g, %v), want (%g, nil)",
					tt.operation, tt.left, tt.right, got, err, tt.want)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		left, right float64
		want        error
	}{
		{name: "unknown operation", operation: "power", want: calculator.ErrUnknownOperation},
		{name: "zero divisor", operation: "divide", left: 7, right: 0, want: calculator.ErrDivisionByZero},
		{
			name: "negative zero divisor", operation: "divide", left: 7,
			right: math.Copysign(0, -1), want: calculator.ErrDivisionByZero,
		},
		{name: "NaN operand", operation: "add", left: math.NaN(), right: 2, want: calculator.ErrInvalidOperand},
		{name: "infinite operand", operation: "subtract", left: 1, right: math.Inf(-1), want: calculator.ErrInvalidOperand},
		{
			name: "multiplication overflow", operation: "multiply", left: math.MaxFloat64,
			right: 2, want: calculator.ErrResultOutOfRange,
		},
		{
			name: "division overflow", operation: "divide", left: math.MaxFloat64,
			right: 0.5, want: calculator.ErrResultOutOfRange,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calculator.Calculate(tt.operation, tt.left, tt.right)
			if !errors.Is(err, tt.want) {
				t.Errorf("Calculate(%q, %g, %g) error = %v, want %v", tt.operation, tt.left, tt.right, err, tt.want)
			}
		})
	}
}
