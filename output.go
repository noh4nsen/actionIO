package actionIO

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

type Output struct{}

func (o Output) Write(config interface{}) error {
	envVars, err := o.extractStructValues(config)
	if err != nil {
		return err
	}

	return o.writeToEnvFile(envVars)
}

func (o Output) writeToEnvFile(envVars map[string]string) error {
	outputPath := os.Getenv("GITHUB_OUTPUT")
	if outputPath == "" {
		return fmt.Errorf("GITHUB_OUTPUT environment variable is not set")
	}

	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	for key, value := range envVars {
		if _, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	return nil
}

func (o Output) extractStructValues(config interface{}) (map[string]string, error) {
	envVars := make(map[string]string)
	configValue := reflect.ValueOf(config)

	if configValue.Kind() != reflect.Ptr || configValue.IsNil() {
		return nil, fmt.Errorf("config must be a non-nil pointer")
	}

	configValue = configValue.Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		fieldValue := configValue.Field(i)
		tag := configType.Field(i).Tag.Get("action")
		if tag != "" {
			if fieldValue.Kind() == reflect.String {
				envVars["INPUT_"+tag] = fieldValue.String()
			} else {
				log.Printf("Field %s is not a string, skipping", configType.Field(i).Name)
			}
		}
	}

	return envVars, nil
}
