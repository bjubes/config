package config

import (
	"log"
	"os"
	"reflect"
	"strconv"
)

// for backwards compatability with go versions <1.18
type any interface{}

type Configurator interface {
	GetEnvString(field string) string
	GetEnvBool(field string) bool
	GetEnvInt(field string) int
	GetEnvFloat(field string) float64
}

// GetEnvString gets a string from the environment, falling back to the same field name in the config struct.
// If it doesn't exist in either, the function will log an error and panic
func GetEnvString(config Configurator, field string) string {
	defer panicOnBadField(config, field, "string")
	value, exists := os.LookupEnv(field)
	if !exists {
		value = getString(config, field)
	}
	return value
}

// GetEnvBool gets a bool from the environment, falling back to the same field name in the config struct.
// If it doesn't exist in either, the function will log an error and panic
func GetEnvBool(config Configurator, field string) bool {
	defer panicOnBadField(config, field, "bool")
	value, err := strconv.ParseBool(os.Getenv(field))
	if err != nil {
		value = getBool(config, field)
	}
	return value
}

// GetEnvInt gets an int from the environment, falling back to the same field name in the config struct.
// If it doesn't exist in either, the function will log an error and panic
func GetEnvInt(config Configurator, field string) int {
	defer panicOnBadField(config, field, "int")
	value, err := strconv.Atoi(os.Getenv(field))
	if err != nil {
		value = getInt(config, field)
	}
	return value
}

// GetEnvFloat gets a float from the environment, falling back to the same field name in the config struct.
// If it doesn't exist in either, the function will log an error and panic
func GetEnvFloat(config Configurator, field string) float64 {
	defer panicOnBadField(config, field, "float")
	value, err := strconv.ParseFloat(os.Getenv(field), 64)
	if err != nil {
		value = getFloat(config, field)
	}
	return value
}

// getString gets the value of an string field in the config struct. Panics if
// the field is not present or not a string
func getString(c any, field string) string {
	c_value := reflect.ValueOf(c)
	field_value := reflect.Indirect(c_value).FieldByName(field)

	if !field_value.IsValid() {
		log.Panicf("configurator has no field named '%s'", field)
	}
	if field_value.Kind() != reflect.String {
		log.Panicf("configurator field '%s' is not a string", field)
	}
	return field_value.String()
}

// getBool gets the value of a bool field in the config struct. Panics if
// the field is not present
func getBool(c any, field string) bool {
	return getNonString(c, field).Bool()
}

// getInt gets the value of an int field in the config struct. Panics if
// the field is not present
func getInt(c any, field string) int {
	return int(getNonString(c, field).Int())
}

// getFloat gets the value of a float field in the config struct. Panics if
// the field is not present
func getFloat(c any, field string) float64 {
	return getNonString(c, field).Float()
}

// panicOnBadField is deferred by the GetEnv* functions. It inspects a panic
// that bubbled up from reading a mis-typed or missing field and re-panics with
// a message that says which type was expected.
func panicOnBadField(config Configurator, field, typeName string) {
	if r := recover(); r != nil {
		if reflect.ValueOf(config).FieldByName(field).IsValid() {
			log.Panicf("configurator field '%s' is not a %s", field, typeName)
		}
		log.Panicf("configurator has no %s field named '%s'", typeName, field)
	}
}

func getNonString(c any, field string) reflect.Value {
	defer func() {
		//catch panic when field does not exist
		if r := recover(); r != nil {
			log.Panicf("configurator has no field named '%s'", field)
		}
	}()
	c_value := reflect.ValueOf(c)
	field_value := reflect.Indirect(c_value).FieldByName(field)
	return field_value
}
