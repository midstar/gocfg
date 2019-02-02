// Package gocfg implements a small and simple configuration / INI
// file parser.
//
// The format of the configuration file is:
//
//      # This is a comment
//      key1 = value1
//      key2 = value2
//
// Supported value types are strings, booleans (on/off, true/false or 1/0),
// integers and floats.
package gocfg

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Configuration is the actual configuration type
type Configuration struct {
	properties map[string]string
}

// LoadConfiguration loads configuration from file and returns a
// Configuration. An empty configuraion is returned on errors.
func LoadConfiguration(fileName string) (*Configuration, error) {
	properties, err := loadPropertyFile(fileName)
	return &Configuration{properties: properties}, err
}

// HasKey returns true if key exists
func (c *Configuration) HasKey(key string) bool {
	_, hasKey := c.properties[key]
	return hasKey
}

// GetString get string value of property. If key don't exist
// the defaultValue is returned.
func (c *Configuration) GetString(key string, defaultValue string) string {
	value, hasKey := c.properties[key]
	if !hasKey {
		return defaultValue
	}
	return value
}

// GetInt get integer value of property. If key don't exist
// the defaultValue is returned and error will be nil. If key
// exist but the value is not a valid integer the defaultValue
// is returned and an error description will be provided.
func (c *Configuration) GetInt(key string, defaultValue int) (int, error) {
	value, hasKey := c.properties[key]
	if !hasKey {
		return defaultValue, nil
	}
	intValue, valueerr := strconv.Atoi(value)
	if valueerr != nil {
		return defaultValue, fmt.Errorf("property %s does not have a valid integer value. Using default %d", key, defaultValue)
	}
	return intValue, nil
}

// GetFloat get float value of property. If key don't exist
// the defaultValue is returned and error will be nil. If key
// exist but the value is not a valid float the defaultValue
// is returned and an error description will be provided.
func (c *Configuration) GetFloat(key string, defaultValue float64) (float64, error) {
	value, hasKey := c.properties[key]
	if !hasKey {
		return defaultValue, nil
	}
	floatValue, valueerr := strconv.ParseFloat(value, 64)
	if valueerr != nil {
		return defaultValue, fmt.Errorf("property %s does not have a valid float value. Using default %f", key, defaultValue)
	}
	return floatValue, nil
}

// boolTrueValues texts that represents true. Case insensitive
var boolTrueValues = [...]string{"on", "true", "yes", "enable", "enabled", "1"}

// boolFalseValues texts that represents false. Case insensitive
var boolFalseValues = [...]string{"off", "false", "no", "disable", "disabled", "0"}

// GetBool get bool value of property. If key don't exist
// the defaultValue is returned and error will be nil. If key
// exist but the value is not a valid bool the defaultValue
// is returned and an error description will be provided.
//
// Valid bool values are on/off, true/false, enable/disable,
// enabled/disabled and 1/0. All are case insensitive, for
// example ON, on, On and oN are all valid true values.
func (c *Configuration) GetBool(key string, defaultValue bool) (bool, error) {
	value, hasKey := c.properties[key]
	if !hasKey {
		return defaultValue, nil
	}
	for _, trueValue := range boolTrueValues {
		if strings.EqualFold(value, trueValue) {
			return true, nil
		}
	}
	for _, falseValue := range boolFalseValues {
		if strings.EqualFold(value, falseValue) {
			return false, nil
		}
	}
	return defaultValue, fmt.Errorf("property %s does not have a valid bool value. Using default %t", key, defaultValue)
}

// loadPropertyFile loads property files of the same format as found in Java
// property files and INI files and returns a map of strings.
func loadPropertyFile(fileName string) (map[string]string, error) {
	properties := make(map[string]string)
	b, fileerr := ioutil.ReadFile(fileName)
	if fileerr != nil {
		return properties, fmt.Errorf("unable to load properties from %s. Reason: %s", fileName, fileerr)
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		tLine := strings.TrimSpace(line)
		if len(tLine) > 0 && !strings.HasPrefix(tLine, "#") {
			parts := strings.SplitN(tLine, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if len(key) > 0 {
					properties[key] = value
				}
			}
		}
	}
	return properties, nil
}
