# Enum_codegen

Pretty enum code generator.

## Links
 - [Overview](#overview)
 - [Requirements](#requirements)
 - [Quick Start](#quick-start)
 - [Generated Code Examples](#generated-code-examples)

## Overview

Enum_codegen is a library for quick enum creation

Generates `var Tags = map[Enum]string`, `var Types = map[string]Enum`, `var Translations = map[Enum]string`.
Implements `sql.Scanner`, `driver.Valuer`, `fmt.Stringer`, `json.Marshaler` interfaces.

Forces Enum to have default Undefined value.

Generates snake_case string Enum value: `EnumValue1 -> enum_value_1`
You can override default string value, and set translation value 

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

## Generated code examples
[Default usage](examples/default_generation/test_file_Enum__gen.go)

## TODO's
- add multi-language support
- add output map naming customization (Tags, Values, Translations)
- add unit tests
- add different string value generators