package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
)

func ReadConfig(path string) (*Config, error) {
	var conf Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	err = toml.Unmarshal(data, &conf)
	if err != nil {
		return nil, fmt.Errorf("umarshalling data: %w", err)
	}
	return &conf, nil
}

// ReadEnv reads struct recursively and set value from Environment if struct tag 'env' exists
func ReadEnv(s interface{}) error {
	reflectType := reflect.TypeOf(s).Elem()
	reflectValue := reflect.ValueOf(s).Elem()

	for i := 0; i < reflectType.NumField(); i++ {
		typeField := reflectType.Field(i)

		value := reflectValue.Field(i)
		kind := value.Kind()

		if kind == reflect.Struct {
			err := ReadEnv(reflectValue.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
		}

		env, ok := typeField.Tag.Lookup("env")
		if !ok || env == "" {
			continue
		}

		v := os.Getenv(env)
		if v == "" {
			continue
		}

		switch kind {

		case reflect.String:
			value.SetString(v)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf(
					"can't parse %s with type %s: %w",
					typeField.Name,
					typeField.Type,
					err,
				)
			}

			value.SetInt(num)
		case reflect.Bool:
			x, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf(
					"can't parse %s with type %s: %w",
					typeField.Name,
					typeField.Type,
					err,
				)
			}

			value.SetBool(x)
		default:
			return fmt.Errorf("can't set %s", reflectType)
		}

	}

	return nil
}
