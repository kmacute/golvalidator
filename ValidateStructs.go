package golvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ValidateStructs(s interface{}) map[string][]string {
	var errors = make(map[string][]string)
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		config := v.Type().Field(i).Tag.Get("validate")
		fieldName := v.Type().Field(i).Tag.Get("json")
		fieldValue := v.Field(i).Interface()

		fieldError := validateField(config, fieldValue, fieldName, v)
		if len(fieldError) > 0 {
			errors[fieldName] = fieldError
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return errors
}

func validateField(config string, fieldValue interface{}, fieldName string, refl reflect.Value) []string {
	var errors []string
	configArray := strings.Split(config, "|")
	v := reflect.ValueOf(fieldValue)
	dataType := ""
	for i := 0; i < len(configArray); i++ {
		validation := strings.Split(configArray[i], ":")
		switch validation[0] {
		case "nullable":
			if isEmptyValue(v) {
				return errors
			}
		default:
			newDataType, validatedField := validate(validation, v, fieldValue, fieldName, dataType, refl)
			if newDataType != "" {
				dataType = newDataType
			}
			if validatedField != "" {
				errors = append(errors, validatedField)
			}
		}
	}

	return errors
}

func validate(validation []string, v reflect.Value, fieldValue interface{}, fieldName string, dataType string, refl reflect.Value) (string, string) {
	localType := dataType
	switch validation[0] {
	case "required":
		if isEmptyValue(v) {
			return localType, fmt.Sprintf("The %s field is required.", RemoveUnderscore(fieldName))
		}

	case "alpha", "string":
		localType = "string"
		if !IsAlpha(v.String()) {
			return localType, fmt.Sprintf("The %s must only contain letters.", RemoveUnderscore(fieldName))
		}

	case "numeric":
		localType = "string"
		if IsFloat(v.String()) {
			return localType, fmt.Sprintf("The %s must be a number.", RemoveUnderscore(fieldName))
		}

	case "alpha_num":
		if !IsAlphanumeric(v.String()) {
			return localType, fmt.Sprintf("The %s must only contain letters and numbers.", RemoveUnderscore(fieldName))
		}

	case "alpha_space":
		localType = "string"
		if !IsAlphaSpace(v.String()) {
			return localType, fmt.Sprintf("The %s must only contain letters, numbers, dashes and underscores.", RemoveUnderscore(fieldName))
		}

	case "alpha_dash":
		localType = "string"
		if !IsAlphaDash(v.String()) {
			return localType, fmt.Sprintf("The %s must only contain letters, numbers, dashes and underscores.", RemoveUnderscore(fieldName))
		}

	case "date":
		if !IsDate(v.String()) {
			return localType, fmt.Sprintf("The %s is not a valid date.", RemoveUnderscore(fieldName))
		}

	case "email":
		if !IsEmail(v.String()) {
			return localType, fmt.Sprintf("The %s must be a valid email address.", RemoveUnderscore(fieldName))
		}

	case "same":
		params := strings.Split(validation[1], ",")
		value2nd := refl.FieldByName(ToCamel(params[0])).Interface()
		if !IsSame(v.String(), value2nd.(string)) {
			return localType, fmt.Sprintf("The %s and %s must match.", RemoveUnderscore(fieldName), RemoveUnderscore(params[0]))
		}

	case "min":
		switch dataType {
		case "string":
			minCharacterLength, _ := strconv.Atoi(validation[1])
			textLength := len(v.String())

			if textLength < minCharacterLength {
				return localType, fmt.Sprintf("The %s must be at least %s characters.", RemoveUnderscore(fieldName), validation[1])
			}

		case "numeric":
			maxValue, _ := strconv.ParseFloat(validation[1], 64)
			floatValue, _ := ToFloat(fieldValue)

			if floatValue < maxValue {
				return localType, fmt.Sprintf("The %s must be at least %s.", RemoveUnderscore(fieldName), validation[1])
			}
		}

	case "max":
		switch dataType {
		case "string":
			maxCharacterLength, _ := strconv.Atoi(validation[1])
			textLength := len(v.String())

			if textLength > maxCharacterLength {
				return localType, fmt.Sprintf("The %s not be greater than %s characters.", RemoveUnderscore(fieldName), validation[1])
			}

		case "numeric":
			maxValue, _ := strconv.ParseFloat(validation[1], 64)
			floatValue, _ := ToFloat(fieldValue)

			if floatValue > maxValue {
				return localType, fmt.Sprintf("The %s not be greater than %s.", RemoveUnderscore(fieldName), validation[1])
			}
		}

	case "between":
		params := strings.Split(validation[1], ",")
		switch dataType {
		case "string":
			min, _ := strconv.Atoi(params[0])
			max, _ := strconv.Atoi(params[1])
			textLength := len(v.String())

			if !IsBetweenInt(textLength, min, max) {
				return localType, fmt.Sprintf("The %s must be between %s and %s characters.", RemoveUnderscore(fieldName), params[0], params[1])
			}

		case "numeric":
			min, _ := strconv.ParseFloat(params[0], 64)
			max, _ := strconv.ParseFloat(params[1], 64)
			floatValue, _ := ToFloat(fieldValue)

			if !IsBetween(floatValue, min, max) {
				return localType, fmt.Sprintf("The %s must be between %s and %s.", RemoveUnderscore(fieldName), params[0], params[1])
			}
		}

	case "digits":
		textLength := len(v.String())
		digits, _ := strconv.Atoi(validation[1])
		if textLength != digits {
			return localType, fmt.Sprintf("The %s must be %s digits.", RemoveUnderscore(fieldName), validation[1])
		}

	case "digits_between":
		params := strings.Split(validation[1], ",")
		min, _ := strconv.Atoi(params[0])
		max, _ := strconv.Atoi(params[1])
		textLength := len(v.String())

		if !IsBetweenInt(textLength, min, max) {
			return localType, fmt.Sprintf("The %s must be between %s and %s digits.", RemoveUnderscore(fieldName), params[0], params[1])
		}

	case "lt":
		switch dataType {
		case "string":
			lessThanCharacterLength, _ := strconv.Atoi(validation[1])
			textLength := len(v.String())

			if !(textLength < lessThanCharacterLength) {
				return localType, fmt.Sprintf("The %s must be less than %s characters.", RemoveUnderscore(fieldName), validation[1])
			}

		case "numeric":
			lessThanValue, _ := strconv.ParseFloat(validation[1], 64)
			floatValue, _ := ToFloat(fieldValue)

			if !(floatValue < lessThanValue) {
				return localType, fmt.Sprintf("The %s must be less than %s.", RemoveUnderscore(fieldName), validation[1])
			}
		}

	case "gt":
		switch dataType {
		case "string":
			greaterThanCharacterLength, _ := strconv.Atoi(validation[1])
			textLength := len(v.String())

			if !(textLength > greaterThanCharacterLength) {
				return localType, fmt.Sprintf("The %s must be greater than %s characters.", RemoveUnderscore(fieldName), validation[1])
			}

		case "numeric":
			greaterThanValue, _ := strconv.ParseFloat(validation[1], 64)
			floatValue, _ := ToFloat(fieldValue)

			if !(floatValue > greaterThanValue) {
				return localType, fmt.Sprintf("The %s must be greater than %s.", RemoveUnderscore(fieldName), validation[1])
			}
		}

	case "lte":
		switch dataType {
		case "string":
			lessThanCharacterLength, _ := strconv.Atoi(validation[1])
			textLength := len(v.String())

			if !(textLength <= lessThanCharacterLength) {
				return localType, fmt.Sprintf("The %s must be less than %s characters.", RemoveUnderscore(fieldName), validation[1])
			}

		case "numeric":
			lessThanValue, _ := strconv.ParseFloat(validation[1], 64)
			floatValue, _ := ToFloat(fieldValue)

			if !(floatValue <= lessThanValue) {
				return localType, fmt.Sprintf("The %s must be less than %s.", RemoveUnderscore(fieldName), validation[1])
			}
		}

	case "gte":
		switch dataType {
		case "string":
			greaterThanCharacterLength, _ := strconv.Atoi(validation[1])
			textLength := len(v.String())

			if !(textLength >= greaterThanCharacterLength) {
				return localType, fmt.Sprintf("The %s must be greater than %s characters.", RemoveUnderscore(fieldName), validation[1])
			}

		case "numeric":
			greaterThanValue, _ := strconv.ParseFloat(validation[1], 64)
			floatValue, _ := ToFloat(fieldValue)

			if !(floatValue >= greaterThanValue) {
				return localType, fmt.Sprintf("The %s must be greater than %s.", RemoveUnderscore(fieldName), validation[1])
			}
		}

	case "required_if":
		params := strings.Split(validation[1], ",")
		secondFieldValue := refl.FieldByName(ToCamel(params[0])).Interface()
		if v.String() != "" {
			return localType, ""
		}

		if secondFieldValue.(string) == params[1] {
			return localType, fmt.Sprintf("The %s field is required when %s is %s.", RemoveUnderscore(fieldName), RemoveUnderscore(params[0]), params[1])
		}

	case "required_with":
		if v.String() != "" {
			return localType, ""
		}

		secondFieldValue := refl.FieldByName(ToCamel(validation[1])).Interface()
		if secondFieldValue != "" {
			return localType, fmt.Sprintf("The %s field is required when %s is present.", RemoveUnderscore(fieldName), RemoveUnderscore(validation[1]))
		}

	case "ip":
		if !IsIP(v.String()) {
			return localType, fmt.Sprintf("The %s must be a valid IP address.", RemoveUnderscore(fieldName))
		}

	case "ipv4":
		if !IsIPv4(v.String()) {
			return localType, fmt.Sprintf("The %s must be a valid IPv4 address.", RemoveUnderscore(fieldName))
		}

	case "ipv6":
		if !IsIPv6(v.String()) {
			return localType, fmt.Sprintf("The %s must be a valid IPv6 address.", RemoveUnderscore(fieldName))
		}

	case "url":
		if !IsURL(v.String()) {
			return localType, fmt.Sprintf("The %s format is invalid.", RemoveUnderscore(fieldName))
		}

	case "credit_card":
		if !IsCreditCard(v.String()) {
			return localType, fmt.Sprintf("The %s must have a valid credit card number.", RemoveUnderscore(fieldName))
		}
	}

	return localType, ""
}
