package encryptions

import (
	"reflect"
	"strings"

	"github.com/roto17/zeus/lib/utils"
	encryptions "github.com/roto17/zeus/lib/utils"
)

func EncryptObjectID(data interface{}) interface{} {
	// Get the type and value of the original data
	originalValue := reflect.ValueOf(data)
	originalType := originalValue.Type()

	// Ensure the input is a struct
	if originalValue.Kind() != reflect.Struct {
		panic("input must be a struct")
	}

	// Create a slice to store the fields of the new struct
	var fields []reflect.StructField

	// Iterate through the fields of the original struct
	for i := 0; i < originalType.NumField(); i++ {
		field := originalType.Field(i)
		fieldType := field.Type

		if fieldType.Kind() == reflect.Struct && fieldType.String() != "time.Time" {
			// If the field is a struct (not time.Time), recursively modify its fields
			subFields := EncryptObjectID(reflect.New(fieldType).Elem().Interface())

			fields = append(fields, reflect.StructField{
				Name: field.Name,
				Type: reflect.TypeOf(subFields), // Set the type of the nested struct
				Tag:  field.Tag,
			})
		} else if strings.HasSuffix(field.Name, "ID") && fieldType.Kind() == reflect.Uint {
			// Modify fields ending with "ID" of type uint to string
			fields = append(fields, reflect.StructField{
				Name: field.Name,
				Type: reflect.TypeOf(""), // Change type to string
				Tag:  field.Tag,
			})
		} else {
			// Otherwise, keep the original field type
			fields = append(fields, field)
		}
	}

	// Create a new struct type with the modified fields
	newStructType := reflect.StructOf(fields)
	newStruct := reflect.New(newStructType).Elem()

	// Copy field values from the original struct to the new struct
	for i := 0; i < originalValue.NumField(); i++ {
		originalField := originalValue.Field(i)
		newField := newStruct.Field(i)

		// Handle "ID" fields
		if strings.HasSuffix(fields[i].Name, "ID") && fields[i].Type.Kind() == reflect.String {
			newField.SetString(utils.EncryptID(uint(originalField.Uint())))
		} else if fields[i].Type.Kind() == reflect.Struct {
			// Recursively handle nested structs
			originalFieldValue := originalField.Interface()
			if reflect.TypeOf(originalFieldValue).String() != "time.Time" {
				encryptedValue := EncryptObjectID(originalFieldValue)

				// Ensure types match before setting
				if reflect.TypeOf(encryptedValue) == newField.Type() {
					newField.Set(reflect.ValueOf(encryptedValue))
				}
			} else {
				newField.Set(originalField) // Copy time.Time fields as is
			}
		} else {
			// Copy other fields as is
			newField.Set(originalField)
		}
	}

	// Return the new struct as an interface{}
	return newStruct.Interface()
}

func DecryptObjectID(data interface{}, target interface{}) interface{} {
	// Get the value and type of the original data
	originalValue := reflect.ValueOf(data)
	originalType := originalValue.Type()

	// Ensure the input is a struct
	if originalValue.Kind() != reflect.Struct {
		panic("input must be a struct")
	}

	// Get the type of the target interface
	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()

	// Ensure the target is a struct
	if targetType.Kind() != reflect.Struct {
		panic("target must be a struct type")
	}

	// Create a new instance of the target type
	newStruct := reflect.New(targetType).Elem()

	// Iterate through the fields of the original struct
	for i := 0; i < originalType.NumField(); i++ {
		originalField := originalValue.Field(i)
		originalFieldType := originalType.Field(i)

		// Find the corresponding field in the target struct
		newField := newStruct.FieldByName(originalFieldType.Name)
		if !newField.IsValid() || !newField.CanSet() {
			continue
		}

		// Handle nested structs
		if originalField.Kind() == reflect.Struct && originalField.Type().String() != "time.Time" {
			// Recursively decrypt nested structs
			decryptedValue := DecryptObjectID(originalField.Interface(), reflect.New(newField.Type()).Interface())

			// Ensure types match before setting
			if reflect.TypeOf(decryptedValue) == newField.Type() {
				newField.Set(reflect.ValueOf(decryptedValue))
			}
		} else if strings.HasSuffix(originalFieldType.Name, "ID") && originalFieldType.Type.Kind() == reflect.String {
			// Decrypt fields ending with "ID"
			if newField.Kind() == reflect.Uint && originalField.IsValid() && originalField.Kind() == reflect.String && originalField.String() != "" {
				newField.SetUint(uint64(encryptions.DecryptID(originalField.String())))
			}
		} else {
			// Copy other fields as is
			newField.Set(originalField)
		}
	}

	// Return the new struct as an interface{}
	return newStruct.Interface()
}
