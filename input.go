package actionIO

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type Input struct{}

func (i Input) Load(config interface{}) error {
	envVarsMap, err := i.getEnvVars()
	if err != nil {
		return err
	}

	configValue := reflect.ValueOf(config)
	if configValue.Kind() != reflect.Ptr || configValue.IsNil() {
		return fmt.Errorf("config must be a non-nil pointer")
	}

	configValue = configValue.Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		fieldValue := configValue.Field(i)

		envTag := configType.Field(i).Tag.Get("action")
		if envTag != "" {
			if envValue, exists := envVarsMap["INPUT_"+envTag]; exists {
				if fieldValue.CanSet() && fieldValue.Kind() == reflect.String {
					fieldValue.SetString(envValue)
				} else {
					log.Printf("Cannot set field %s", configType.Field(i).Name)
				}
			}
		}
	}
	return nil
}

func (i Input) getEnvVars() (map[string]string, error) {
	envVarsMap := make(map[string]string)
	envVarsList := os.Environ()

	for _, envVar := range envVarsList {
		keyValuePair := strings.SplitN(envVar, "=", 2)
		if len(keyValuePair) == 2 {
			envVarsMap[keyValuePair[0]] = keyValuePair[1]
		} else {
			return nil, fmt.Errorf("invalid environment variable format: %s", envVar)
		}
	}
	return envVarsMap, nil
}
