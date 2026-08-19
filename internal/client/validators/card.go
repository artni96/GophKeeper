package validators

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// PANValidator checks if the PAN is valid according to Luhn algorithm.
func PANValidator(number uint64) error {
	strNumber := strconv.FormatUint(number, 10)
	if len(strNumber) != 16 {
		return fmt.Errorf("Invalid PAN: number length should be 16 digits")
	}
	sum := 0
	strLength := len(strNumber)
	parity := strLength % 2

	for i := 0; i < strLength-1; i++ {
		digit := int(strNumber[i] - '0')
		if parity == (i+1)%2 {
			sum += digit
		} else if digit > 4 {
			sum += 2*digit - 9
		} else {
			sum += 2 * digit
		}
	}
	controlSum := (10 - (sum % 10)) % 10
	ok := int(strNumber[strLength-1]-'0') == controlSum
	if !ok {
		return fmt.Errorf("Invalid PAN")
	}
	return nil
}

// CVVValidator validates the CVV value.
func CVVValidator(cvv string) error {
	var errs []string
	if cvv != "" {
		intCVV, err := strconv.Atoi(cvv)
		if err != nil {
			errs = append(errs, fmt.Sprintf("invalid cvv value: %d", intCVV))
		}
		if (len(cvv) < 3 || len(cvv) > 5) && intCVV < 0 {
			errs = append(errs, fmt.Sprintf("invalid cvv value: %d", intCVV))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// PINValidator validates the PIN value.
func PINValidator(pin string) error {
	var errs []string
	if pin != "" {
		intPIN, err := strconv.Atoi(pin)
		if err != nil {
			errs = append(errs, fmt.Sprintf("invalid pin value: %d", intPIN))
		}
		if len(pin) != 4 && (intPIN < 0 || intPIN > 10000) {
			errs = append(errs, fmt.Sprintf("invalid pin value: %d", intPIN))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// ExpiryDateValidator validates the Expiry date value.
func ExpiryDateValidator(expiryDate string) error {
	var errs []string
	if expiryDate != "" {
		strMonth := expiryDate[:2]
		intMonth, err := strconv.Atoi(strMonth)
		if err != nil {
			errs = append(errs, fmt.Sprintf("invalid expiry date value: %s", expiryDate))
		}
		if intMonth < 0 || intMonth > 12 {
			errs = append(errs, fmt.Sprintf("invalid expiry date value: %s", expiryDate))
		}
		strYear := expiryDate[3:]
		intYear, err := strconv.Atoi(strYear)
		if err != nil {
			errs = append(errs, fmt.Sprintf("invalid expiry date value: %s", expiryDate))
		}
		if intYear < 0 || intYear > 99 {
			errs = append(errs, fmt.Sprintf("invalid expiry date value: %s", expiryDate))
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}
