# Enum_codegen

Fast, useful, easy-to-use enum code generator.

## Links
 - [Overview](#overview)
 - [Requirements](#requirements)
 - [Quick Start](#quick-start)
 - [Generated Code Example](#generated-code-example)

## Overview

Enum_codegen is a library for quick enum creation

Creates maps `var Tags = map[Enum]string`, `var Types = map[string]Enum`, `var Translations = map[Enum]string`
Automatically implements `sql.Scanner`, `driver.Valuer`, `fmt.Stringer`, `json.Marshaler` interfaces

Forces Enum to have default Undefined value.

Automatically generates Enum string value with format: `EnumValue1 -> enum_value_1`

Using comments allows to specify string value and translation value

## Requirements

- go 1.23.0 or later

## Quick Start

1. Install enum_codegen:

   ```shell
   go install github.com/GeekchanskiY/enum_codegen/cmd/enum_codegen
   ```
    
    Before other steps make sure that go/bin is in your $PATH
    

2. Create a Go file with an enum declaration:

   ```go
    package main

    //go:generate enum_codegen
    type Enum int
    
    const (
    // Undefined Translate="Undefined value"
    Undefined Enum = iota
    // EnumValue1 Value="SuperbValue" Translate="Enum 1"
    EnumValue1
    // EnumValue2 Your documentation here if need Translate="Enum 2" Value="MegaValue"
    EnumValue2
    // EnumValue3 Translate="Enum 3"
    EnumValue3
    // Translate="Some value Enum 4"
    EnumValue4
    EnumValue5
    )
   ```

3. Run go generate:
    
    you can specify the file to generate from explicitly
    
   ```shell
   go generate package/file.go
   ```
   
    or generate from every file in project. Result will not change

    ```shell
   go generate ./...
    ```

## Generated code example
```go
// Code generated via enum_codegen DO NOT EDIT.
package test_package

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	_ sql.Scanner   = (*Enum)(nil)
	_ driver.Valuer = (*Enum)(nil)
	_ fmt.Stringer  = (*Enum)(nil)
	_ json.Marshaler = (*Enum)(nil)
)

var Tags = map[Enum]string{ 
		Undefined: "Undefined", 
		EnumValue1: "Sample value", 
		EnumValue2: "enum_value_2", 
		EnumValue3: "enum_value_3", 
		EnumValue4: "enum_value_4", 
		EnumValue5: "enum_value_5",
}

var Types = map[string]Enum{ 
		"Undefined": Undefined, 
		"Sample value": EnumValue1, 
		"enum_value_2": EnumValue2, 
		"enum_value_3": EnumValue3, 
		"enum_value_4": EnumValue4, 
		"enum_value_5": EnumValue5,
}

var Translations = map[Enum]string{ 
		Undefined: "Enum Value is undefined", 
		EnumValue1: "enum_value_1", 
		EnumValue2: "enum_value_2", 
		EnumValue3: "Enum ultra value", 
		EnumValue4: "Enum last value", 
		EnumValue5: "enum_value_5",
}

func (t *Enum) Scan(src any) error {
	value, ok := src.(string)
	
	if !ok {
		return errors.New("src is not string")
	}

	*t = Undefined
	if v, ok := Types[value]; ok {
		*t = v
	}

	return nil
}

func (t Enum) Value() (driver.Value, error) {
	return Tags[t], nil
}

func (t Enum) String() string {
	return Tags[t]
}

func (t Enum) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *Enum) UnmarshalJSON(data []byte) error {
	var (
		s string
		err error
	)

	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}

	if v, ok := Types[s]; ok {
		*t = v

		return nil
	}

	return fmt.Errorf("invalid status: %s", s)
}
```