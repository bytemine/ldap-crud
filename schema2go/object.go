package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

// template to generate code for an object
const objectTemp = `{{define "OB"}} // {{ .Name }}{{if .Desc}}: {{ .Desc }}{{end}}
type {{ .Name }} struct {
dn	string
{{ range $k, $v := .ObjectClasses }}{{ $k | title | replace }} bool
{{ end }}

{{ if .Must }}//MUST attributes
{{ template "ATDECL" .Must }}{{ end }}
{{if .May}}// MAY attributes
{{ template "ATDECL" .May }}{{ end }}
}

func New{{ .Name }}(dn string) *{{.Name}} {
o := new({{ .Name }})
o.dn = dn
return o
}

func (o *{{ .Name }})FilterObjectClass() string {
return "{{ .FilterObjectClass }}"
}

func (o *{{ .Name }})Copy() crud.Item {
c := New{{ .Name }}(o.dn)
{{ template "ATCOPY" .Must }}
{{ template "ATCOPY" .May }}
return c
}

func (o *{{ .Name }})Dn() string {
return o.dn
}

{{ if .DNFormat }}func (o *{{ .Name }})FormatDn() {
o.dn = fmt.Sprintf("{{ .DNFormat }}", []string{ {{range $v := .DNAttributes}}o.{{$v | title | replace }}, {{end}} }...)
}{{end}}

func (o *{{ .Name }})MarshalLDAP() (*ldap.Entry, error) {
e := ldap.NewEntry(o.dn)
{{ range $k, $v := .ObjectClasses }}
if o.{{ $k | title | replace }} { 
e.AddAttributeValue("objectClass", "{{$k}}")
{{range $k, $v := $v.Must}}if len(o.{{ $k | title | replace }}) == 0 {
return nil, errors.New(fmt.Sprintf("Marshalling %v: Attribute %v is empty", "{{ $.Name }}", "{{ $k | title | replace }}"))
}
{{ end }}
{{ template "ATMARSHAL" .Must }}
{{ template "ATMARSHAL" .May }}
} {{ end }}
return e, nil
}

func (o *{{ .Name }})UnmarshalLDAP(e *ldap.Entry) error {
o.dn = e.DN

for _, v := range e.GetAttributeValues("objectClass") {
	switch strings.ToLower(v) {
		{{range $k, $v := .ObjectClasses}}case "{{ $k | lower }}":
			o.{{ $k | title | replace }} = true
		{{end}}
	}
}

{{ template "ATUNMARSHAL" .Must }}
{{ template "ATUNMARSHAL" .May }}
return nil
}{{end}}
{{define "ATDECL"}}{{ range $k, $v := . }} // {{ $v.Desc }}
{{$k | title | replace }} {{ if $v.SingleValue}}string{{else}}[]string{{end}} ` + "`json:\",omitempty\"`" + `
{{end}}{{end}}
{{define "ATCOPY" }}{{ range $k, $v := . }}
{{ if $v.SingleValue }}c.{{ $k | title | replace }} = o.{{ $k | title | replace }}{{else}}c.{{ $k | title | replace }} = make([]string, len(o.{{ $k | title | replace}}))
copy(c.{{ $k | title | replace }}, o.{{ $k | title | replace }}){{end}}{{end}}{{end}}
{{define "ATMARSHAL" }}{{range $k, $v := .}}
{{ if $v.SingleValue }}e.AddAttributeValue("{{ $k }}", o.{{ $k | title | replace }}){{else}}e.AddAttributeValues("{{ $k }}", o.{{ $k | title | replace }}){{end}}{{end}}{{end}}
{{define "ATUNMARSHAL" }}{{range $k, $v := .}}
{{ if $v.SingleValue }}o.{{ $k | title | replace }} = e.GetAttributeValue("{{ $k }}"){{else}}o.{{ $k | title | replace }} = e.GetAttributeValues("{{ $k }}"){{end}}{{end}}{{end}}
{{ template "OB" . }}
`

// An object is a description of which objectclasses an object found in the directory could
// consist. Out of this description the code necessary for marshaling to and unmarshaling from
// such objects is generated with the Code method.
type Object struct {
	Name              string
	Desc              string
	ObjectClasses     []string
	FilterObjectClass string
	DNFormat          string
	DNAttributes      []string
}

// Generates Go-code for itself. The struct and it's methods implement the Item interface of
// package crud.
func (o Object) Code() (string, error) {
	funcMap := template.FuncMap{
		"title":   strings.Title,
		"replace": nameReplacer.Replace,
		"lower":   strings.ToLower,
		"join": func(l []string) string {
			return strings.Join(l, ", ")
		},
	}

	t := template.Must(template.New("object").Funcs(funcMap).Parse(objectTemp))

	data := struct {
		Name          string
		Desc          string
		ObjectClasses map[string]struct {
			Must map[string]*attributetype
			May  map[string]*attributetype
		}
		FilterObjectClass string
		DNFormat          string
		DNAttributes      []string
		Must              map[string]*attributetype
		May               map[string]*attributetype
	}{Name: o.Name, Desc: o.Desc, FilterObjectClass: o.FilterObjectClass}

	data.ObjectClasses = make(map[string]struct {
		Must map[string]*attributetype
		May  map[string]*attributetype
	})
	data.Must = make(map[string]*attributetype)
	data.May = make(map[string]*attributetype)

	for _, ocName := range o.ObjectClasses {
		oc, ok := objectclassdefs[strings.ToLower(ocName)]
		if !ok {
			return "", errors.New(fmt.Sprintf("undefined object class: %v\n", ocName))
		}

		must := make(map[string]*attributetype)

		for _, v := range oc.Must {
			var attr *attributetype
			if _, ok := data.Must[v]; ok {
				continue
			}

			if a, ok := attributetypedefs[strings.ToLower(v)]; ok {
				attr = a
			} else {
				attr = &attributetype{Name: []string{v}, Desc: "attribute definition missing"}
			}

			must[v] = attr
			data.Must[v] = attr
		}

		may := make(map[string]*attributetype)

		for _, v := range oc.May {
			var attr *attributetype
			if _, ok := data.May[v]; ok {
				continue
			}

			// ignore attributes present as must
			if _, ok := data.Must[v]; ok {
				continue
			}

			if a, ok := attributetypedefs[strings.ToLower(v)]; ok {
				attr = a
			} else {
				attr = &attributetype{Name: []string{v}, Desc: "attribute definition missing"}
			}

			may[v] = attr
			data.May[v] = attr
		}

		tmp := struct {
			Must map[string]*attributetype
			May  map[string]*attributetype
		}{Must: must, May: may}

		data.ObjectClasses[ocName] = tmp
	}

	if o.DNFormat != "" && len(o.DNAttributes) == 0 {
		return "", errors.New("DNFormat without DNAttributes")
	}

	for _, v := range o.DNAttributes {
		if _, ok := data.Must[v]; ok {
			continue
		}

		if _, ok := data.May[v]; ok {
			continue
		}

		return "", errors.New(fmt.Sprintf("Undefined attribute %v in DNAttributes", v))
	}

	data.DNFormat = o.DNFormat
	data.DNAttributes = o.DNAttributes

	var w bytes.Buffer
	err := t.Execute(&w, data)
	return w.String(), err

}

// unmarshal the json object definitions
func parseObjects(r []byte) ([]Object, error) {
	objects := make([]Object, 0)
	err := json.Unmarshal(r, &objects)

	return objects, err
}
