package freeipa

import (
	"fmt"
	"unicode"
)

func StringIsNotOnlyDigits(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	for _, chr := range v {
		if !unicode.IsDigit(chr) {
			return nil, nil
		}
	}

	return nil, []error{fmt.Errorf("expected %q to not only contain digits", k)}
}

func StringContainsNoUpperLetter(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	for _, chr := range v {
		if unicode.IsUpper(chr) {
			return nil, []error{fmt.Errorf("expected %q to not contain any upper letter", k)}
		}
	}

	return nil, nil
}
