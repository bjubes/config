# Config [![Go Reference](https://pkg.go.dev/badge/github.com/bjubes/config.svg)](https://pkg.go.dev/github.com/bjubes/config)

A go package for creating default configuration variables that can be overridden by environment variables.

### Install
```
go get github.com/bjubes/config
```

### Usage

1. Create a struct with fields that match your environment variables, and create an instance set to your defaults.
	```go
	type MyConfig struct {
		DB_HOST string
		DB_PORT int
		PROD    bool
	}
	var myConfig config.Configurator = MyConfig{
		DB_HOST: "localhost",
		DB_PORT: 5432,
		PROD:    false,
	}
	```

2. Make your custom struct implement the `Configurator` interface using the following code (just copy and paste)
   ```go
	func (c MyConfig) GetEnvInt(field string) int {
		return config.GetEnvInt(c, field)
	}
	func (c MyConfig) GetEnvString(field string) string {
		return config.GetEnvString(c, field)
	}
	func (c MyConfig) GetEnvBool(field string) bool {
		return config.GetEnvBool(c, field)
	}
   ```

3. Retrieve a value using the methods on your config instance 
    ```go
	host := myConfig.GetEnvString("DB_HOST")
	port := myConfig.GetEnvInt("DB_PORT")
	prod := myConfig.GetEnvBool("PROD")
	```

Values will default to what they are set to in the struct instance, but will be overridden by environment variables if they are set.
Environment variables must match the type specified. For specifics, see [type matching rules](#type-matching-rules) below.

Since the `myConfig` instance has a type of `Configurator`, none of the public fields are accesssable. This forces retreiving values through the `GetEnv` methods, so you never accidentally grab the default value without checking for the environment variable first.


### Type matching rules

**string** - Environment value will be used as long as it is set, even if its an empty string.

**int** - Environment value will be used if the value is an integer, using [`strconv.Atoi`](https://pkg.go.dev/strconv#Atoi).

**bool** - Environment value will be used as long as the environment variable is set. If it is set, it will be true unless the value is one of `false`, `0`, or `off` (case insenstive).
