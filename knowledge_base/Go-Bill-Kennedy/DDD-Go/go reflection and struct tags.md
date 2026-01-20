#go #reflection #struct-tags

## Parser
```go
package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// expands each field if root struct has nested struct
// example:
/*
	type T struct {
		A string
		B struct {
			C string
		}
	}

	Field will be computed for each of A and C,
	Field for A: Field{Name:"A", Key: []string{"A"}}
	Field for C: Field{Name:"C", Key: []string{"B", "C"}} -> can be referenced as B_C
*/
func ExtractFields(inp interface{}) ([]Field, error) {
	val := reflect.ValueOf(inp)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("t must be of type struct")
	}
	val = val.Elem()
	typ := val.Type()

	fields := []Field{}
	for i := range val.NumField() {
		v := val.Field(i)
		t := typ.Field(i)
		field := Field{
			Name: t.Name,
			Key:  []string{},
		}

		field.Key = append(field.Key, t.Name)

		if t.Type.Kind() == reflect.Struct {
			nestedVal := v.Addr().Interface()
			nestedFields, err := ExtractFields(nestedVal)
			if err != nil {
				return nil, err
			}
			field.Name = nestedFields[len(nestedFields)-1].Name
			for _, v := range nestedFields {
				field.Key = append(field.Key, v.Key...)
			}
		}

		if t.Type.Kind() != reflect.Struct {
			confTag := t.Tag.Get("conf")
			if strings.Trim(confTag, " ") != "" {
				confs := strings.SplitN(confTag, ":", 2)
				if len(confs) == 2 {
					switch confs[0] {
					case "default":
						fmt.Printf("default for %v:%v\n", field.Name, confs[1])
						switch t.Type.Kind() {
						case reflect.String:
							v.SetString(confs[1])
						case reflect.Int:
							in, err := strconv.ParseInt(confs[1], 0, t.Type.Bits())
							if err != nil {
								return nil, fmt.Errorf("invalid default value: %w", err)
							}
							v.SetInt(in)
						}
					}
				}
			}
		}

		fields = append(fields, field)
	}

	return fields, nil
}

// represents each field of root struct
type Field struct {
	Name string
	Key  []string
	Tag  string
}
```
### Tests
```go
package parser

import "testing"

func TestExtractFields(t *testing.T) {
	cfg := struct {
		Version string
		Build   string
		Web     struct {
			APIHost string `conf:"default:test-host"`
			Port    string `conf:"default:test-port"`
		}
		DB struct {
			ConnStr string `conf:"default:test-conn-str"`
			Nested  struct {
				TestField    string `conf:"default:test-field"`
				TestIntField int    `conf:"default:10"`
			}
		}
	}{}

	fields, err := ExtractFields(&cfg)
	if err != nil {
		t.Errorf("should not have failed. error: %v", err)
	}
	for _, v := range fields {
		t.Log(v.Name)
		t.Log(v.Key)
		t.Log("---")
	}
	t.Log(cfg)
}
```
```sh
go test -v .
=== RUN   TestExtractFields
default for APIHost:test-host
default for Port:test-port
default for ConnStr:test-conn-str
default for TestField:test-field
default for TestIntField:10
    parser_test.go:27: Version
    parser_test.go:28: [Version]
    parser_test.go:29: ---
    parser_test.go:27: Build
    parser_test.go:28: [Build]
    parser_test.go:29: ---
    parser_test.go:27: Port
    parser_test.go:28: [Web APIHost Port]
    parser_test.go:29: ---
    parser_test.go:27: TestIntField
    parser_test.go:28: [DB ConnStr Nested TestField TestIntField]
    parser_test.go:29: ---
    parser_test.go:31: {  {test-host test-port} {test-conn-str {test-field 10}}}
--- PASS: TestExtractFields (0.00s)
PASS
ok      github.com/vishal2098govind/service/cmd/learn-reflect/parser    0.485s
```