# enum_codegen

Generate polished Go enum helpers from plain `const` blocks.

`enum_codegen` turns integer-backed enums into a small, predictable API for string conversion, JSON, SQL, and display metadata. Keep your source enum readable, add one `//go:generate` directive, and let the generator write the repetitive code.

## What It Generates

For an enum named `Status`, the generator creates:

- `StatusTags`: `map[Status]string` for enum-to-string lookup.
- `StatusTypes`: `map[string]Status` for string-to-enum lookup.
- `StatusTranslations`: `map[Status]string` for display text.
- `String() string`.
- `MarshalJSON() ([]byte, error)`.
- `UnmarshalJSON([]byte) error`.
- `Scan(any) error` for `database/sql`.
- `Value() (driver.Value, error)` for database writes.

The generated names are enum-scoped, so multiple enums can live in the same package without `Tags`, `Types`, or `Translations` collisions.

## Requirements

- Go `1.23.0` or newer.
- Integer-backed enum types.
- `go generate` for code generation.

## Install

```shell
go install github.com/GeekchanskiY/enum_codegen/cmd/enum_codegen@latest
```

Make sure your Go binary directory is available in `PATH`:

```shell
export PATH="$(go env GOPATH)/bin:$PATH"
```

## Quick Start

Create an enum and place `//go:generate enum_codegen` directly above the type declaration:

```go
package orders

//go:generate enum_codegen
type Status int

const (
	// Undefined Value="undefined" Translate="Unknown status"
	Undefined Status = iota
	// StatusDraft Translate="Draft order"
	StatusDraft
	// StatusPaid Value="paid" Translate="Paid order"
	StatusPaid
	// StatusCancelled Translate="Cancelled order"
	StatusCancelled
)
```

Run generation for one file:

```shell
go generate ./status.go
```

Or generate everything in the module:

```shell
go generate ./...
```

The output file will be named like:

```text
status_Status__gen.go
```

## Comment Metadata

By default, enum names are converted from PascalCase or camelCase to snake_case:

```text
StatusCancelled -> status_cancelled
EnumValue1      -> enum_value_1
```

You can override the generated string value and display text with comment metadata:

```go
// StatusPaid Value="paid" Translate="Paid order"
StatusPaid
```

Supported metadata:

- `Value="..."`: custom string representation used by `String`, JSON, SQL, and lookup maps.
- `Translate="..."`: display text stored in `<EnumName>Translations`.

## Undefined Behavior

If your enum contains a constant named exactly `Undefined`, generated SQL scanning uses it as a safe fallback for `nil` or unknown database values.

```go
const (
	Undefined Status = iota
	StatusDraft
)
```

If there is no `Undefined` constant, unknown or `nil` scanned values return an error instead of generating code that fails to compile.

To require an `Undefined` value during generation, use:

```go
//go:generate enum_codegen -f
```

or:

```go
//go:generate enum_codegen --force-undefined
```

## Generated API Example

```go
var status Status

_ = status.Scan([]byte("paid"))
fmt.Println(status == StatusPaid) // true

data, _ := json.Marshal(StatusCancelled)
fmt.Println(string(data)) // "status_cancelled"

fmt.Println(StatusTranslations[StatusDraft]) // Draft order
```

## Examples

- [Default generation](examples/default_generation/file.go)
- [Generated output](examples/default_generation/file_Enum__gen.go)
- [Forced undefined validation](examples/force_undefined/file.go)

## Development

Run the full test suite:

```shell
go test ./...
```

Run static checks:

```shell
go vet ./...
```

Run coverage for the implementation packages and checked-in generated example:

```shell
go test ./cmd/enum_codegen ./pkg/enum ./pkg/generator ./pkg/parser ./examples/default_generation -cover
```

## Roadmap

- Multi-language translation support.
- Configurable output map names.
- Additional string value generation strategies.

## License

This project is available under the [MIT License](LICENSE).