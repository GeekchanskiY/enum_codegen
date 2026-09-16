package generator

import (
	"io"
	"text/template"

	"github.com/GeekchanskiY/enum_codegen/pkg/enum"
)

const Template = `// Code generated via enum_codegen DO NOT EDIT.
package {{ .PackageName }}

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Force interface implementation
var (
	_ sql.Scanner   = (*{{ .EnumName }})(nil)
	_ driver.Valuer = (*{{ .EnumName }})(nil)
	_ fmt.Stringer  = (*{{ .EnumName }})(nil)
	_ json.Marshaler = (*{{ .EnumName }})(nil)
)

var {{ .EnumName }}Tags = map[{{ .EnumName }}]string{
	{{- range .Enums }} 
		{{ .Name }}: {{ printf "%q" .SnakeName }},
	{{- end }}
}

var {{ .EnumName }}Types = map[string]{{ .EnumName }}{
	{{- range .Enums }} 
		{{ printf "%q" .SnakeName }}: {{ .Name }},
	{{- end }}
}

var {{ .EnumName }}Translations = map[{{ .EnumName }}]string{
	{{- range .Enums }} 
		{{ .Name }}: {{ printf "%q" .Translate }},
	{{- end }}
}

func (t *{{ .EnumName }}) Scan(src any) error {
	var value string
	switch v := src.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	case nil:
		{{- if .HasUndefined }}
		*t = Undefined
		return nil
		{{- else }}
		return fmt.Errorf("src is nil")
		{{- end }}
	default:
		return fmt.Errorf("src is %T, not string or []byte", src)
	}

	if v, ok := {{ .EnumName }}Types[value]; ok {
		*t = v
		return nil
	}

	{{- if .HasUndefined }}
	*t = Undefined
	return nil
	{{- else }}
	return fmt.Errorf("invalid {{ .EnumName }}: %s", value)
	{{- end }}
}

func (t {{ .EnumName }}) Value() (driver.Value, error) {
	return {{ .EnumName }}Tags[t], nil
}

func (t {{ .EnumName }}) String() string {
	return {{ .EnumName }}Tags[t]
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

	if v, ok := {{ .EnumName }}Types[s]; ok {
		*t = v

		return nil
	}

	return fmt.Errorf("invalid {{ .EnumName }}: %s", s)
}
`

func CompileTemplate(wr io.Writer, packageName, enumName string, data enum.Enum) error {
	return compileTemplate(wr, Template, packageName, enumName, data)
}

func compileTemplate(wr io.Writer, templateText, packageName, enumName string, data enum.Enum) error {
	tmpl, err := template.New("enum_code").Parse(templateText)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, map[string]any{
		"PackageName":  packageName,
		"Enums":        data,
		"EnumName":     enumName,
		"HasUndefined": data.HasUndefined(),
	})
}
