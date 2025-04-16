package util

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	nullStr = "NULL"
)

type ConvertOpertion func(v *reflect.Value, fieldName string, valueStr string) error

func ConvertPercent(v *reflect.Value, fieldName string, valueStr string) error {
	field := v.FieldByName(fieldName)
	floatValue, err := strconv.ParseFloat(valueStr[:len(valueStr)-1], 32)

	if err == nil {
		field.SetFloat(floatValue / 100)
	} else {
		err = errors.New(fieldName + "is not percent! value is " + valueStr)
	}

	return err
}

func ConvertStorageSize(v *reflect.Value, fieldName string, valueStr string) error {
	field := v.FieldByName(fieldName)
	floatValue, err := strconv.ParseFloat(valueStr[:len(valueStr)-1], 32)

	if err == nil {
		if strings.HasSuffix(valueStr, "P") {
			floatValue *= 1024
		}

		field.SetFloat(floatValue)
	} else {
		err = errors.New(fieldName + "is not storage size! value is " + valueStr)
	}

	return err
}

func Convert(v *reflect.Value, fieldName string, valueStr string) error {
	var err error = nil
	var floatValue float64 = 0
	var intValue int64 = 0

	field := v.FieldByName(fieldName)
	if field.Type().Kind() >= reflect.Int && field.Type().Kind() <= reflect.Float64 {
		valueStr = strings.ReplaceAll(valueStr, ",", "")
	}

	switch field.Type().Kind() {
	case reflect.Float32:
		if nullStr == valueStr || valueStr == "" {
			floatValue = 0
		} else {
			floatValue, err = strconv.ParseFloat(valueStr, 32)
		}
	case reflect.Int32:
		intValue, err = strconv.ParseInt(valueStr, 10, 32)
	case reflect.Int64:
		intValue, err = strconv.ParseInt(valueStr, 10, 64)
	case reflect.String:
	default:
		err = errors.New("Don't support type, type : " + field.Type().Kind().String() + "; fieldName : " + fieldName)
	}

	if err == nil {
		switch field.Type().Kind() {
		case reflect.Float32:
			field.SetFloat(floatValue)
		case reflect.Int32:
			field.SetInt(intValue)
		case reflect.Int64:
			field.SetInt(intValue)
		case reflect.String:
			field.SetString(valueStr)
		default:
		}
	}

	return err
}
