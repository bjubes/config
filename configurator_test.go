package config

import (
	"os"
	"testing"
)

// Setup dummy config object. Good example of how to do in a real project
type TestConfig struct {
	DB_HOST string
	DB_PORT int
	PROD    bool
}

var config_struct = TestConfig{
	DB_HOST: "host",
	DB_PORT: 1,
	PROD:    true,
}

var config_iface Configurator = config_struct

func (c TestConfig) GetEnvInt(field string) int {
	//in real project would be config.GetEnvInt(...)
	return GetEnvInt(c, field)
}

func (c TestConfig) GetEnvString(field string) string {
	//in real project would be config.GetEnvString(...)
	return GetEnvString(c, field)
}

func (c TestConfig) GetEnvBool(field string) bool {
	//in real project would be config.GetEnvBool(...)
	return GetEnvBool(c, field)
}

// End setup

func TestReadValues(t *testing.T) {
	host := getString(config_iface, "DB_HOST")
	port := getInt(config_iface, "DB_PORT")
	prod := getBool(config_iface, "PROD")

	if host != config_struct.DB_HOST {
		t.Fatalf("Getting string value failed. Got '%v'. want '%v'", host, config_struct.DB_HOST)
	}

	if port != config_struct.DB_PORT {
		t.Fatalf("Getting int value failed. Got '%v'. want '%v'", port, config_struct.DB_PORT)
	}

	if prod != config_struct.PROD {
		t.Fatalf("Getting bool value failed. Got '%v'. want '%v'", port, config_struct.PROD)
	}
}

func TestEnvOverridesConfig(t *testing.T) {
	os.Setenv("DB_HOST", "newhost")
	os.Setenv("DB_PORT", "2")
	os.Setenv("PROD", "false")

	host := config_iface.GetEnvString("DB_HOST")
	port := config_iface.GetEnvInt("DB_PORT")
	prod := config_iface.GetEnvBool("PROD")

	if host != "newhost" {
		t.Fatalf("Getting string value failed. Got '%v'. want '%v'", host, "newhost")
	}

	if port != 2 {
		t.Fatalf("Getting int value failed. Got '%v'. want '%v'", port, 2)
	}

	if prod != false {
		t.Fatalf("Getting bool value failed. Got '%v'. want '%v'", prod, false)
	}

}

func TestConfigFallbackToEnv(t *testing.T) {
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("PROD")

	host := config_iface.GetEnvString("DB_HOST")
	port := config_iface.GetEnvInt("DB_PORT")
	prod := config_iface.GetEnvBool("PROD")

	if host != config_struct.DB_HOST {
		t.Fatalf("Getting string value failed. Got '%v'. want '%v'", host, config_struct.DB_HOST)
	}

	if port != config_struct.DB_PORT {
		t.Fatalf("Getting int value failed. Got '%v'. want '%v'", port, config_struct.DB_PORT)
	}

	if prod != config_struct.PROD {
		t.Fatalf("Getting bool value failed. Got '%v'. want '%v'", port, config_struct.DB_PORT)
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

func TestBoolWithVariousTrues(t *testing.T) {
	falses := []string{
		"false", "off", "0",
	}
	trues := []string{
		"", "anything", "true", "on", "1",
	}

	config_struct.PROD = false
	for _, elem := range falses {
		os.Setenv("PROD", elem)
		if config_iface.GetEnvBool("PROD") {
			t.Fatalf("setting env to '%v' should be interpreted bool as false", elem)
		}
	}
	config_struct.PROD = true
	for _, elem := range trues {
		os.Setenv("PROD", elem)
		if !config_iface.GetEnvBool("PROD") {
			t.Fatalf("setting env to '%v' should be interpreted bool as true", elem)
		}
	}
}

func use(vals ...any) {
	for _, val := range vals {
		_ = val
	}
}
