package encryptions

import (
	"reflect"

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

		// If the field is named "ID" and is of type uint, change its type to string
		if field.Name == "ID" && field.Type.Kind() == reflect.Uint {
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

	// Create a new instance of the newly defined struct type
	newStruct := reflect.New(newStructType).Elem()

	// Copy field values from the original struct to the new struct
	for i := 0; i < originalValue.NumField(); i++ {
		originalField := originalValue.Field(i)
		newField := newStruct.Field(i)

		// If the field is named "ID" and is of type string, convert the uint to string
		if fields[i].Name == "ID" && fields[i].Type.Kind() == reflect.String {
			newField.SetString(encryptions.EncryptID(uint(originalField.Uint())))
		} else {
			// Otherwise, copy the value as is
			newField.Set(originalField)
		}
	}

	// Return the new struct as an interface{}
	return newStruct.Interface()
}

func DecryptObjectID(data interface{}) interface{} {
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

		// If the field is named "ID" and is of type uint, change its type to string
		if field.Name == "ID" && field.Type.Kind() == reflect.String {
			fields = append(fields, reflect.StructField{
				Name: field.Name,
				Type: reflect.TypeOf(uint(1)), // Change type to string
				Tag:  field.Tag,
			})
		} else {
			// Otherwise, keep the original field type
			fields = append(fields, field)
		}
	}

	// Create a new struct type with the modified fields
	newStructType := reflect.StructOf(fields)

	// Create a new instance of the newly defined struct type
	newStruct := reflect.New(newStructType).Elem()

	// Copy field values from the original struct to the new struct
	for i := 0; i < originalValue.NumField(); i++ {
		originalField := originalValue.Field(i)
		newField := newStruct.Field(i)

		// If the field is named "ID" and is of type string, convert the uint to string
		if fields[i].Name == "ID" && fields[i].Type.Kind() == reflect.Uint {
			// newField.SetString(encryptions.DecryptID(originalField.String()))
			newField.SetUint(uint64(encryptions.DecryptID(originalField.String())))
		} else {
			// Otherwise, copy the value as is
			newField.Set(originalField)
		}
	}

	// Return the new struct as an interface{}
	return newStruct.Interface()
}
