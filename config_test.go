// Unit tests for  gocfg
package gocfg

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNonExistingConfig(t *testing.T) {
	config, err := LoadConfiguration("dont_exist.cfg")
	if err == nil {
		t.Fatalf("Non existing configuration should return error")
	}
	if len(config.properties) != 0 {
		t.Fatalf("There should be no properties on invalid config file, but it was %d", len(config.properties))
	}
}

func TestConfig(t *testing.T) {
	contents :=
		`key1=value1
		# This is a comment
key2  =   value2
just ignore this
	key 3 = value  3`
	fullPath := createConfigFile(t, "TestConfig.cfg", contents)
	config, err := LoadConfiguration(fullPath)
	if err != nil {
		t.Fatal(err)
	}

	value1, ok := config.properties["key1"]
	if !ok {
		t.Fatal(err)
	}
	if value1 != "value1" {
		t.Fatalf("Expected value1 but was %s", value1)
	}

	value2, ok := config.properties["key2"]
	if !ok {
		t.Fatal(err)
	}
	if value2 != "value2" {
		t.Fatalf("Expected value2 but was %s", value1)
	}

	value3, ok := config.properties["key 3"]
	if !ok {
		t.Fatal(err)
	}
	if value3 != "value  3" {
		t.Fatalf("Expected 'value  3' but was %s", value1)
	}
}

func TestHasKey(t *testing.T) {
	contents :=
		`key1=874
key2=invalid int
key3=3.14`
	fullPath := createConfigFile(t, "TestGetValues.cfg", contents)
	config, err := LoadConfiguration(fullPath)
	if err != nil {
		t.Fatal(err)
	}

	// Test HasKey
	if config.HasKey("key2") == false {
		t.Fatalf("key2 shall exist")
	}

	if config.HasKey("dontexist") == true {
		t.Fatalf("dontexist shall not exist")
	}

}

func TestGetString(t *testing.T) {
	contents :=
		`key1=874
key2=invalid int
key3=3.14`
	fullPath := createConfigFile(t, "TestGetString.cfg", contents)
	config, err := LoadConfiguration(fullPath)
	if err != nil {
		t.Fatal(err)
	}
	// Test GetString
	valueStr := config.GetString("key1", "default")
	if valueStr != "874" {
		t.Fatalf("Expected 874 but was %s", valueStr)
	}

	valueStr = config.GetString("dontexist", "hello world")
	if valueStr != "hello world" {
		t.Fatalf("Expected default 'hello world' but was %s", valueStr)
	}

}

func TestGetInt(t *testing.T) {
	contents :=
		`key1=874
key2=invalid int
key3=3.14`
	fullPath := createConfigFile(t, "TestGetInt.cfg", contents)
	config, err := LoadConfiguration(fullPath)
	if err != nil {
		t.Fatal(err)
	}

	// Test GetInt
	valueInt, err := config.GetInt("key1", 3)
	if err != nil {
		t.Fatal(err)
	}
	if valueInt != 874 {
		t.Fatalf("Expected 874 but was %d", valueInt)
	}

	valueInt, err = config.GetInt("key2", 3)
	if err == nil {
		t.Fatal("A string as an integer should give errors")
	}
	if valueInt != 3 {
		t.Fatalf("Expected default 3 but was %d", valueInt)
	}

	valueInt, err = config.GetInt("key3", 99)
	if err == nil {
		t.Fatal("A float as an integer should give errors")
	}
	if valueInt != 99 {
		t.Fatalf("Expected default 99 but was %d", valueInt)
	}

	valueInt, err = config.GetInt("dontexist", 32)
	if err != nil {
		t.Fatal("A missing key shall not give errors")
	}
	if valueInt != 32 {
		t.Fatalf("Expected default 32 but was %d", valueInt)
	}

}

func TestGetFloat(t *testing.T) {
	contents :=
		`key1=874
key2=invalid int
key3=3.14`
	fullPath := createConfigFile(t, "TestGetFloat.cfg", contents)
	config, err := LoadConfiguration(fullPath)
	if err != nil {
		t.Fatal(err)
	}

	// Test GetFloat
	valueFloat, err := config.GetFloat("key1", 3)
	if err != nil {
		t.Fatal(err)
	}
	if valueFloat != 874 {
		t.Fatalf("Expected 874 but was %f", valueFloat)
	}

	valueFloat, err = config.GetFloat("key2", 9.4323223)
	if err == nil {
		t.Fatal("A string as a float should give errors")
	}
	if valueFloat != 9.4323223 {
		t.Fatalf("Expected default 9.4323223 but was %f", valueFloat)
	}

	valueFloat, err = config.GetFloat("key3", 99)
	if err != nil {
		t.Fatal(err)
	}
	if valueFloat != 3.14 {
		t.Fatalf("Expected default 3.14 but was %f", valueFloat)
	}

	valueFloat, err = config.GetFloat("dontexist", 23423.67234)
	if err != nil {
		t.Fatal("A missing key shall not give errors")
	}
	if valueFloat != 23423.67234 {
		t.Fatalf("Expected default 23423.67234 but was %f", valueFloat)
	}

}

func TestGetBool(t *testing.T) {
	contents :=
		`key1=on
key2=OFF
key3=true
key4=False
key5=fake`
	fullPath := createConfigFile(t, "TestGetBool.cfg", contents)
	config, err := LoadConfiguration(fullPath)
	if err != nil {
		t.Fatal(err)
	}

	// Test GetBool
	valueBool, err := config.GetBool("key1", false)
	if err != nil {
		t.Fatal(err)
	}
	if valueBool != true {
		t.Fatalf("Expected true but was %t", valueBool)
	}

	valueBool, err = config.GetBool("key2", true)
	if err != nil {
		t.Fatal(err)
	}
	if valueBool != false {
		t.Fatalf("Expected false but was %t", valueBool)
	}

	valueBool, err = config.GetBool("key3", false)
	if err != nil {
		t.Fatal(err)
	}
	if valueBool != true {
		t.Fatalf("Expected true but was %t", valueBool)
	}

	valueBool, err = config.GetBool("key4", true)
	if err != nil {
		t.Fatal(err)
	}
	if valueBool != false {
		t.Fatalf("Expected false but was %t", valueBool)
	}

	valueBool, err = config.GetBool("key5", true)
	if err == nil {
		t.Fatal("A non valid bool string should give errors")
	}
	if valueBool != true {
		t.Fatalf("Expected default true but was %t", valueBool)
	}

	valueBool, err = config.GetBool("dontexist", true)
	if err != nil {
		t.Fatal(err)
	}
	if valueBool != true {
		t.Fatalf("Expected default true but was %t", valueBool)
	}

}

// createConfigFile creates a configuration file. Returns the full path to it.
func createConfigFile(t *testing.T, name, contents string) string {
	os.MkdirAll("tmpout", os.ModeDir)
	fullName := "tmpout/" + name
	os.Remove(fullName) // Remove old if it exist
	err := ioutil.WriteFile(fullName, []byte(contents), 0644)
	if err != nil {
		t.Fatalf("Unable to create configuration file. Reason: %s", err)
	}
	return fullName
}

/*
func defaultConfig(t *testing.T) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		t.Fatal("Environment GOPATH needs to be set for this test")
	}
	return filepath.Join(gopath, "src", "github.com", "midstar", "plm", "plm.config")
}

func TestConfig(t *testing.T) {
	config := LoadConfiguration(defaultConfig(t))
	assertEqualsInt(t, "config.Port", 12124, config.Port)
	assertEqualsInt(t, "config.FastLogTimeMs", 6000, config.FastLogTimeMs)
	assertEqualsInt(t, "config.SlowLogFactor", 10, config.SlowLogFactor)
	assertEqualsInt(t, "config.FastLogSize", 600, config.FastLogSize)
	assertEqualsInt(t, "config.SlowLogSize", 1440, config.SlowLogSize)
}

func TestConfigInvalidFile(t *testing.T) {
	config := LoadConfiguration("dont_exist.properties")
	assertEqualsInt(t, "config.Port", 12124, config.Port)
	assertEqualsInt(t, "config.FastLogTimeMs", 3000, config.FastLogTimeMs)
	assertEqualsInt(t, "config.SlowLogFactor", 20, config.SlowLogFactor)
	assertEqualsInt(t, "config.FastLogSize", 1200, config.FastLogSize)
	assertEqualsInt(t, "config.SlowLogSize", 1440, config.SlowLogSize)
}

func TestLoadPropertyInt(t *testing.T) {
	properties := make(map[string]string)
	properties["mkey"] = "invalid"
	value := getPropertyInt(properties, "mkey", 3)
	assertEqualsInt(t, "mKey value", 3, value)
}

func TestLoadProperties(t *testing.T) {
	properties, errprop := LoadPropertyFile(defaultConfig(t))
	if errprop != nil {
		t.Fatal(errprop)
	}
	assertEqualsInt(t, "Size of properties", 5, len(properties))
	assertEqualsStr(t, "Value of property port", "12124", properties["port"])
	assertEqualsStr(t, "Value of property fastLogTimeMs", "6000", properties["fastLogTimeMs"])
	assertEqualsStr(t, "Value of property slowLogFactor", "10", properties["slowLogFactor"])
	assertEqualsStr(t, "Value of property fastLogSize", "600", properties["fastLogSize"])
	assertEqualsStr(t, "Value of property slowLogSize", "1440", properties["slowLogSize"])
}

func TestLoadPropertiesInvalidFile(t *testing.T) {
	properties, err := oadPropertyFile("dont_exist.properties")
	if err == nil {
		t.Fatal("Expected an error when loading properties")
	}
	assertEqualsInt(t, "Size of properties", 0, len(properties))
}
*/
