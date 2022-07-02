package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// for backwards compatability with go versions <1.18
type any interface{}

type Configurator interface {
	GetEnvInt(field string) int
	GetEnvString(field string) string
	GetEnvBool(field string) bool
}

// GetEnvInt gets an int from the environment, falling back to the same field name in the config struct.
// If it doesn't exist in either, the function will log an error then exit 1
func GetEnvInt(config Configurator, field string) int {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
	}()
	value, err := strconv.Atoi(os.Getenv(field))
	if err != nil {
		value = getInt(config, field)
	}
	return value
}

// GetEnvBool gets a bool from the environment, falling back to the same field name in the config struct.
// If it doesn't exist in either, the function will log an error then exit 1
func GetEnvBool(config Configurator, field string) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
	}()
	value, exists := os.LookupEnv(field)
	if !exists {
		return getBool(config, field)
	}
	v := strings.ToLower(value)
	return !(v == "off" || v == "false" || v == "0")
}

// GetEnvString gets a string from the environment, falling back to the same field name in the config struct.
// If it doesn't exist in either, the function will log an error then exit 1
func GetEnvString(config Configurator, field string) string {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(r)
		}
	}()
	value, exists := os.LookupEnv(field)
	if !exists {
		value = getString(config, field)
	}
	return value
}

// getInt gets the value of an int field in the config struct. Panics if
// the field is not present
func getInt(c any, field string) int {
	defer func() {
		//catch panic when field does not exist
		if r := recover(); r != nil {
			panic(fmt.Sprintf("env var not in config struct: `%v`", field))
		}
	}()
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(field)
	return int(f.Int())
}

// getString gets the value of an string field in the config struct. Panics if
// the field is not present
func getString(c any, field string) string {
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(field)

	//String() doesn't panic. It returns '<T Value>' if the field doesn't exist
	result := f.String()
	re := regexp.MustCompile(`<[a-zA-Z]*\sValue>`)
	if re.MatchString(result) {
		panic(fmt.Sprintf("env var not in config struct: `%v`", field))
	}
	return result
}

// getBool gets the value of a bool field in the config struct. Panics if
// the field is not present
func getBool(c any, field string) bool {
	defer func() {
		//catch panic when field does not exist
		if r := recover(); r != nil {
			panic(fmt.Sprintf("env var not in config struct: `%v`", field))
		}
	}()
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Bool()
}
