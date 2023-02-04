package config

import (
	"os"
	"strings"
	"testing"
)

// Setup dummy config object. Good example of how to do in a real project
type TestConfig struct {
	DB_HOST  string
	DB_PORT  int
	PROD     bool
	COOLDOWN float64
}

var config_struct = &TestConfig{
	DB_HOST:  "host",
	DB_PORT:  1,
	PROD:     true,
	COOLDOWN: 0.5,
}

var config_iface Configurator = config_struct

func (c TestConfig) GetEnvString(field string) string {
	//in real project would be config.GetEnvString(...)
	return GetEnvString(c, field)
}

func (c TestConfig) GetEnvBool(field string) bool {
	//in real project would be config.GetEnvBool(...)
	return GetEnvBool(c, field)
}

func (c TestConfig) GetEnvInt(field string) int {
	//in real project would be config.GetEnvInt(...)
	return GetEnvInt(c, field)
}

func (c TestConfig) GetEnvFloat(field string) float64 {
	//in real project would be config.GetEnvFloat(...)
	return GetEnvFloat(c, field)
}

// End setup

func TestReadValues(t *testing.T) {
	host := getString(config_iface, "DB_HOST")
	port := getInt(config_iface, "DB_PORT")
	prod := getBool(config_iface, "PROD")
	cooldown := getFloat(config_iface, "COOLDOWN")

	if host != config_struct.DB_HOST {
		t.Fatalf("Getting string value failed. Got '%v'. want '%v'", host, config_struct.DB_HOST)
	}

	if port != config_struct.DB_PORT {
		t.Fatalf("Getting int value failed. Got '%v'. want '%v'", port, config_struct.DB_PORT)
	}

	if prod != config_struct.PROD {
		t.Fatalf("Getting bool value failed. Got '%v'. want '%v'", port, config_struct.PROD)
	}

	if cooldown != config_struct.COOLDOWN {
		t.Fatalf("Getting float value failed. Got '%v'. want '%v'", port, config_struct.COOLDOWN)
	}
}

func TestEnvOverridesConfig(t *testing.T) {
	os.Setenv("DB_HOST", "newhost")
	os.Setenv("DB_PORT", "2")
	os.Setenv("PROD", "false")
	os.Setenv("COOLDOWN", "0.6")

	host := config_iface.GetEnvString("DB_HOST")
	port := config_iface.GetEnvInt("DB_PORT")
	prod := config_iface.GetEnvBool("PROD")
	cooldown := config_iface.GetEnvFloat("COOLDOWN")

	if host != "newhost" {
		t.Fatalf("Getting string value failed. Got '%v'. want '%v'", host, "newhost")
	}

	if port != 2 {
		t.Fatalf("Getting int value failed. Got '%v'. want '%v'", port, 2)
	}

	if prod != false {
		t.Fatalf("Getting bool value failed. Got '%v'. want '%v'", prod, false)
	}

	if cooldown != 0.6 {
		t.Fatalf("Getting float value failed. Got '%v'. want '%v'", cooldown, 0.6)
	}

}

func TestConfigFallbackToEnv(t *testing.T) {
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("PROD")
	os.Unsetenv("COOLDOWN")

	host := config_iface.GetEnvString("DB_HOST")
	port := config_iface.GetEnvInt("DB_PORT")
	prod := config_iface.GetEnvBool("PROD")
	cooldown := config_iface.GetEnvFloat("COOLDOWN")

	if host != config_struct.DB_HOST {
		t.Fatalf("Getting string value failed. Got '%v'. want '%v'", host, config_struct.DB_HOST)
	}

	if port != config_struct.DB_PORT {
		t.Fatalf("Getting int value failed. Got '%v'. want '%v'", port, config_struct.DB_PORT)
	}

	if prod != config_struct.PROD {
		t.Fatalf("Getting bool value failed. Got '%v'. want '%v'", port, config_struct.DB_PORT)
	}

	if cooldown != config_struct.COOLDOWN {
		t.Fatalf("Getting float value failed. Got '%v'. want '%v'", port, config_struct.COOLDOWN)
	}

}

func TestWrongIntFieldCausesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Reading non existant field did not cause panic")
		}
	}()

	value := getInt(config_iface, "DOES_NOT_EXIST")
	use(value)
}
func TestWrongStringFieldCausesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Reading non existant field did not cause panic")
		}
	}()

	value := getString(config_iface, "DOES_NOT_EXIST")
	use(value)
}

func TestWrongBoolFieldCausesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Reading non existant field did not cause panic")
		}
	}()

	value := getBool(config_iface, "DOES_NOT_EXIST")
	use(value)
}

func TestWrongFloatFieldCausesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Reading non existant field did not cause panic")
		}
	}()

	value := getFloat(config_iface, "DOES_NOT_EXIST")
	use(value)
}

func TestBoolWithVariousEnvStrings(t *testing.T) {
	falses := []string{
		"0", "f", "F", "FALSE", "false", "False",
	}
	trues := []string{
		"1", "t", "T", "TRUE", "true", "True",
	}

	config_struct.PROD = true
	for _, elem := range falses {
		os.Setenv("PROD", elem)
		if config_iface.GetEnvBool("PROD") {
			t.Fatalf("setting boolean env to '%v' should be interpreted bool as false", elem)
		}
	}
	config_struct.PROD = false
	for _, elem := range trues {
		os.Setenv("PROD", elem)
		if !config_iface.GetEnvBool("PROD") {
			t.Fatalf("setting boolean env to '%v' should be interpreted bool as true", elem)
		}
	}
}

func TestBoolWithBadEnvStrings(t *testing.T) {
	invalid_bool_strings := []string{
		"off", "on", "yes", "no", "y", "n", "nope", "nada", "anything", "",
	}

	config_struct.PROD = true
	for _, elem := range invalid_bool_strings {
		os.Setenv("PROD", elem)
		if !config_iface.GetEnvBool("PROD") {
			t.Fatalf("setting boolean env to '%v' should fallback to default value", elem)
		}
	}
	config_struct.PROD = false
	for _, elem := range invalid_bool_strings {
		os.Setenv("PROD", elem)
		if config_iface.GetEnvBool("PROD") {
			t.Fatalf("setting boolean env to '%v' should fallback to default value", elem)
		}
	}
}

func TestPanicMessageIncludesFieldName(t *testing.T) {
	panicCheck := func(f func()) {
		field_name := "DOES_NOT_EXIST"
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("Reading non existant field did not cause panic")
			} else {
				switch v := r.(type) {
				case string:
					if !strings.Contains(v, field_name) {
						t.Errorf("Panic msg should contain field name '%s'. Got '%v'", field_name, r)
					}
				case error:
					if !strings.Contains(v.Error(), "DOES_NOT_EXIST") {
						t.Errorf("Panic msg should contain field name '%s'. Got '%v'", field_name, r)
					}
				default:
					t.Errorf("Panic's type is unkown. Cannot check error msg. Type is '%T'", r)
				}
			}
		}()
		f()
	}
	panicCheck(func() { config_iface.GetEnvInt("DOES_NOT_EXIST") })
	panicCheck(func() { config_iface.GetEnvFloat("DOES_NOT_EXIST") })
	panicCheck(func() { config_iface.GetEnvString("DOES_NOT_EXIST") })
	panicCheck(func() { config_iface.GetEnvBool("DOES_NOT_EXIST") })
}

func use(vals ...any) {
	for _, val := range vals {
		_ = val
	}
}
