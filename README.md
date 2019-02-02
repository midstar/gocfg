# gocfg - Golang configuration, property, ini file reader 

[![Documentation](https://godoc.org/github.com/midstar/gocfg?status.svg)](https://godoc.org/github.com/midstar/gocfg)
[![Go Report Card](https://goreportcard.com/badge/github.com/midstar/gocfg)](https://goreportcard.com/report/github.com/midstar/gocfg)
[![AppVeyor](https://ci.appveyor.com/api/projects/status/github/midstar/gocfg?svg=true)](https://ci.appveyor.com/api/projects/status/github/midstar/gocfg)
[![Coverage Status](https://coveralls.io/repos/github/midstar/gocfg/badge.svg?branch=master)](https://coveralls.io/github/midstar/gocfg?branch=master)


Package gocfg is a simple configuration, property and ini file reader.

It supports configurations files with following format

	# This is a comment
	key1 = value1
	key2 = value2

Supported value types are strings, booleans (on/off, true/false, 1/0 ..),
integers and floats.

## Install

```bash
go get github.com/midstar/gocfg
```

## Import

```go
import (
	"github.com/midstar/gocfg"
)
```

## Example Usage

Example of a configuration file, example.cfg:

	# This is a comment
	name = Foo Bar
	age = 34
	height = 1.72
	male = true
	likes bananas = no

Code for reading that file:

```go
package main

import (
	"fmt"
	"github.com/midstar/gocfg"
)

func main() {

	// Load configuration
	config, _ := gocfg.LoadConfiguration("example.cfg")

	valueStr := config.GetString("name", "")
	fmt.Printf("Name is %s\n", valueStr)

	valueInt, _ := config.GetInt("age", 0)
	fmt.Printf("Age is %d\n", valueInt)

	valueFloat, _ := config.GetFloat("height", 0)
	fmt.Printf("Height is %f\n", valueFloat)

	valueBool, _ := config.GetBool("male", false)
	fmt.Printf("Male %t\n", valueBool)

	valueBool, _ = config.GetBool("likes bananas", true)
	fmt.Printf("Likes bananas %t\n", valueBool)

	valueInt, _ = config.GetInt("key_dont_exist", 5)
	fmt.Printf("key_dont_exist default value %d\n", valueInt)

	fmt.Printf("Has key Foot size? %t\n", config.HasKey("Foot size"))

}
```

Output will be:

	Name is Foo Bar
	Age is 34
	Height is 1.720000
	Male true
	Likes bananas false
	key_dont_exist default value 5 
	Has key Foot size? false

## Author and license

This library is written by Joel Midstjärna and is licensed under the MIT License.