package jsonschema

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"
)

// GenerateSchema создаёт JSON Schema для переданного значения v,
// используя рефлексию и тег jsonschema.
// Возвращает схему в виде json.RawMessage.
func GenerateSchema(v any) (json.RawMessage, error) {
	if v == nil {
		return nil, errors.New("input is nil")
	}
	t := reflect.TypeOf(v)
	schema := buildSchema(t)
	raw, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(raw), nil
}

// buildSchema рекурсивно строит JSON-схему для типа t.
func buildSchema(t reflect.Type) map[string]any {
	// Разыменовываем указатели
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// Специальная обработка времени
	if t == reflect.TypeOf(time.Time{}) {
		return map[string]any {
			"type":   "string",
			"format": "date-time",
		}
	}

	schema := make(map[string]any)

	switch t.Kind() {
	case reflect.String:
		schema["type"] = "string"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema["type"] = "integer"

	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"

	case reflect.Bool:
		schema["type"] = "boolean"

	case reflect.Struct:
		schema["type"] = "object"
		props := make(map[string]any)
		required := []string{}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// Пропускаем неэкспортируемые поля
			if field.PkgPath != "" {
				continue
			}

			jsonName := jsonFieldName(field)
			if jsonName == "" {
				continue // помечено json:"-"
			}

			// Читаем тег jsonschema
			isRequired, desc := parseJsonschemaTag(field.Tag.Get("jsonschema"))

			// Схема для типа поля
			fieldSchema := buildSchema(field.Type)
			if desc != "" {
				fieldSchema["description"] = desc
			}
			props[jsonName] = fieldSchema

			if isRequired {
				required = append(required, jsonName)
			}
		}

		schema["properties"] = props
		schema["required"] = required

	case reflect.Slice, reflect.Array:
		schema["type"] = "array"
		schema["items"] = buildSchema(t.Elem())

	case reflect.Map:
		// JSON Schema поддерживает только строковые ключи
		if t.Key().Kind() == reflect.String {
			schema["type"] = "object"
			schema["additionalProperties"] = buildSchema(t.Elem())
		}
		// Для нестроковых ключей оставляем пустую схему (любое значение)

	case reflect.Interface:
		// interface{} — допустим любой тип
		// Оставляем schema пустым (допускает всё)
	}

	return schema
}

// jsonFieldName возвращает имя поля в JSON на основе тега json.
// Для `json:"-"` возвращает пустую строку, для `json:",omitempty"` — имя поля структуры.
func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}
	parts := strings.Split(tag, ",")
	name := parts[0]
	if name == "-" {
		return ""
	}
	if name == "" {
		return field.Name
	}
	return name
}

// parseJsonschemaTag разбирает тег jsonschema и извлекает флаг required и описание.
// Формат: `required,description=Some text`
func parseJsonschemaTag(tag string) (required bool, description string) {
	if tag == "" {
		return false, ""
	}
	parts := strings.Split(tag, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "required" {
			required = true
		} else if strings.HasPrefix(p, "description=") {
			description = strings.TrimPrefix(p, "description=")
		}
	}
	return
}
