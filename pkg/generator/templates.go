package generator

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/GeekchanskiY/enum_codegen/pkg/enum"
)

const Template = `// Code generated via enum_codegen DO NOT EDIT.
package {{ .PackageName }}

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Force interface implementation
var (
	_ sql.Scanner   = (*{{ .EnumName }})(nil)
	_ driver.Valuer = (*{{ .EnumName }})(nil)
	_ fmt.Stringer  = (*{{ .EnumName }})(nil)
	_ json.Marshaler = (*{{ .EnumName }})(nil)
)

var Tags = map[{{ .EnumName }}]string{
	{{- range .Enums }} 
		{{ .Name }}: "{{ .SnakeName }}",
	{{- end }}
}

var Types = map[string]{{ .EnumName }}{
	{{- range .Enums }} 
		"{{ .SnakeName }}": {{ .Name }},
	{{- end }}
}

var Translations = map[{{ .EnumName }}]string{
	{{- range .Enums }} 
		{{ .Name }}: "{{ .Translate }}",
	{{- end }}
}

func (t *{{ .EnumName }}) Scan(src any) error {
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

func (t {{ .EnumName }}) Value() (driver.Value, error) {
	return Tags[t], nil
}

func (t {{ .EnumName }}) String() string {
	return Tags[t]
}

func (t {{ .EnumName }}) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *{{ .EnumName }}) UnmarshalJSON(data []byte) error {
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
`

func CompileTemplate(wr io.Writer, packageName, enumName string, data enum.Enum) error {
	tmpl, err := template.New("enum_code").Parse(Template)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse template: %s\n", err)
		os.Exit(1)
	}

	return tmpl.Execute(wr, map[string]any{
		"PackageName": packageName,
		"Enums":       data,
		"EnumName":    enumName,
	})
}
