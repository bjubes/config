# Config

A go package for creating default configuration variables that can be overridden by environment variables.


### Usage

1. Create a struct with fields which match the environment variables you want to set and an instance with defaults:
	```go
	type MyConfig struct {
		DB_HOST string
		DB_PORT int
		PROD    bool
	}
	var myConfig config.Configurator = MyConfig{
		DB_HOST: "host",
		DB_PORT: 1,
		PROD:    true,
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
3. Retrieve a value using the methods on your config object
    ```go
	host := config.GetEnvString("DB_HOST")
	port := config.GetEnvInt("DB_PORT")
	prod := config.GetEnvBool("PROD")
	```

Values will default to what they are set to in the struct, but will be overridden by environment variables if they are set to values, provided they match the type.

since your config instance has a type of Configurator, none of the public fields are accesssable, so you are forced to go through the GetEnv methods. This means you will never grab the default without checking for the environment variable by mistake
