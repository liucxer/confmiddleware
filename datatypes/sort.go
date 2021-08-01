package datatypes

import (
	"bytes"
	"go/ast"
	"reflect"
	"strings"

	"github.com/go-courier/sqlx/v2"
	"github.com/go-courier/sqlx/v2/builder"
)

// openapi:strfmt sort-rule
type Sort struct {
	By  string
	Asc bool
}

func (sort Sort) IsZero() bool {
	return sort.By == ""
}

func (sort *Sort) OrderFor(db sqlx.DBExecutor, model builder.Model) *builder.Order {
	if sort.By == "" {
		return nil
	}

	fieldName := ""

	eachStructField(reflect.Indirect(reflect.ValueOf(model)).Type(), func(structField reflect.StructField, jsonName string) bool {
		if sort.By == jsonName {
			fieldName = structField.Name
			return false
		}
		return true
	})

	f := db.T(model).F(fieldName)

	if sort.Asc {
		return builder.AscOrder(f)
	}

	return builder.DescOrder(f)
}

func eachStructField(structType reflect.Type, fn func(structField reflect.StructField, jsonName string) bool) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		if field.Type.Kind() == reflect.Interface {
			continue
		}

		if ast.IsExported(field.Name) {
			tagDBValue, tagDBExists := field.Tag.Lookup("db")
			if tagDBExists {
				if tagDBValue != "-" {
					tagJSONValue, tagJSONExists := field.Tag.Lookup("json")
					if tagJSONExists {
						if tagJSONValue != "-" {
							if !fn(field, outputName(field.Name, tagJSONValue)) {
								break
							}
						}
					}
				}
			} else if field.Anonymous {
				eachStructField(field.Type, fn)
				continue
			}
		}
	}
}

func outputName(fieldName, tagValue string) string {
	n := strings.Split(tagValue, ",")[0]
	if n == "" {
		return fieldName
	}
	return n
}

func (sort Sort) MarshalText() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if sort.By != "" {
		buf.WriteString(sort.By)

		if sort.Asc {
			buf.WriteString("!asc")
		}
	}

	return buf.Bytes(), nil
}

func (sort *Sort) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	parts := bytes.Split(data, []byte("!"))

	s := Sort{}

	if len(parts) > 0 {
		s.By = string(parts[0])
	}

	if len(parts) > 1 {
		if string(parts[1]) == "asc" {
			s.Asc = true
		}
	}

	*sort = s

	return nil
}
