package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ExtractAmount(amount string) (float64, error) {
	if amount == "" {
		return 0, fmt.Errorf("input cannot be empty")
	}

	parsedAmount, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		return 0, fmt.Errorf("invalid input: %s", amount)
	}

	if parsedAmount <= 0 {
		return 0, fmt.Errorf("input %s must be greater than zero", amount)
	}

	return parsedAmount, nil
}

func ExtractInputs(input string) (string, float64, error) {
	result := strings.SplitN(input, "=", 2)

	if len(result) > 1 {
		code := strings.ToUpper(strings.TrimSpace(result[0]))
		amount, err := ExtractAmount(strings.TrimSpace(result[1]))
		return code, amount, err
	}

	return strings.ToUpper(strings.TrimSpace(result[0])), 0, nil

}
