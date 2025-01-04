package encryptions

import (
	"fmt"
	"reflect"
	"strings"

	models "github.com/roto17/zeus/lib/models/products"
	"github.com/roto17/zeus/lib/utils"
	encryptions "github.com/roto17/zeus/lib/utils"
)

func EncryptObjectID(data interface{}) interface{} {
	originalValue := reflect.ValueOf(data)
	originalType := originalValue.Type()

	// Ensure that the input is a struct
	if originalValue.Kind() != reflect.Struct {
		panic("input must be a struct")
	}

	var fields []reflect.StructField

	// Iterate over fields of the struct
	for i := 0; i < originalType.NumField(); i++ {
		field := originalType.Field(i)
		fieldType := field.Type

		// Handle nested structs
		if fieldType.Kind() == reflect.Struct {
			if fieldType.String() == "time.Time" {
				// Keep time fields as is
				fields = append(fields, field)
			} else {
				// Recursively handle nested structs
				subFields := EncryptObjectID(reflect.New(fieldType).Elem().Interface())
				fields = append(fields, reflect.StructField{
					Name: field.Name,
					Type: reflect.TypeOf(subFields),
					Tag:  field.Tag,
				})
			}
		} else if strings.HasSuffix(field.Name, "ID") && (fieldType.Kind() == reflect.Int || fieldType.Kind() == reflect.Uint) {
			// Encrypt fields that end with 'ID' and are of type int or uint
			fields = append(fields, reflect.StructField{
				Name: field.Name,
				Type: reflect.TypeOf(""), // Convert to string
				Tag:  field.Tag,
			})
		} else {
			// Preserve non-ID fields
			fields = append(fields, field)
		}
	}

	// Create the new struct with modified fields
	newStructType := reflect.StructOf(fields)
	newStruct := reflect.New(newStructType).Elem()

	// Iterate over original fields and set values
	for i := 0; i < originalValue.NumField(); i++ {
		originalField := originalValue.Field(i)
		newField := newStruct.Field(i)

		// Encrypt ID fields (and convert them to string)
		if strings.HasSuffix(fields[i].Name, "ID") && (fields[i].Type.Kind() == reflect.String) {
			encryptedID := utils.EncryptID(uint(originalField.Uint())) // Encrypt as uint
			encryptedIDStr := fmt.Sprintf("%v", encryptedID)           // Convert to string
			newField.SetString(encryptedIDStr)
		} else if fields[i].Type.Kind() == reflect.Struct {
			// Recursively handle nested structs (like orderproducts)
			originalFieldValue := originalField.Interface()

			if reflect.TypeOf(originalFieldValue).String() != "time.Time" {
				// Recursively encrypt nested structs
				encryptedValue := EncryptObjectID(originalFieldValue)

				// Ensure that the types match
				if reflect.TypeOf(encryptedValue) == newField.Type() {
					newField.Set(reflect.ValueOf(encryptedValue))
				} else {
					// Handle type mismatch during encryption
					fmt.Println("Encrypted type does not match the expected type")
				}
			} else {
				// For time fields, copy as is
				newField.Set(originalField)
			}
		} else if fields[i].Type.Kind() == reflect.Slice {
			// Handle slices like 'orderproducts'
			originalSlice := originalField.Interface()
			originalSliceValue := reflect.ValueOf(originalSlice)

			// Create a slice to hold encrypted elements (this is where we collect the encrypted elements first)
			var encryptedElements []interface{}

			// Iterate over the original slice and encrypt elements before adding them to the list
			for j := 0; j < originalSliceValue.Len(); j++ {
				originalElement := originalSliceValue.Index(j)

				// Ensure that the element is a struct
				if originalElement.Kind() == reflect.Struct {
					// Encrypt the struct element
					encryptedElement := EncryptObjectID(originalElement.Interface())

					// Ensure that the encrypted element is of the correct type
					if encryptedElement != nil {
						// Collect the encrypted element
						encryptedElements = append(encryptedElements, encryptedElement)
					} else {
						// Log and handle type mismatch during encryption
						fmt.Printf("Encrypted element is nil for index %d\n", j)
					}
				} else {
					// If not a struct, just append the original element
					encryptedElements = append(encryptedElements, originalElement.Interface())
				}
			}

			// Now that we have all encrypted elements, create the slice using the element type of the encrypted elements
			if len(encryptedElements) > 0 {
				// Get the type of the encrypted elements (this will be the slice element type)
				// elementType := reflect.TypeOf(encryptedElements[0])

				// Create a new slice of type []models.OrderProduct
				// This is where we ensure the slice type matches the expected type
				encryptedSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(models.OrderProductEncrypted{})), len(encryptedElements), len(encryptedElements))

				// Populate the slice with the encrypted elements
				for j, encryptedElement := range encryptedElements {
					// Set the value of the encrypted element
					encryptedSlice.Index(j).Set(reflect.ValueOf(encryptedElement))
				}

				// Set the newly created slice in the new struct field
				newField.Set(encryptedSlice)
			} else {
				// If no elements, just copy the original slice
				newField.Set(originalField)
			}
		} else {
			// Copy non-ID fields as is
			newField.Set(originalField)
		}
	}

	// Return the new struct as an interface{}
	return newStruct.Interface()
}

// func EncryptObjectID(data interface{}) interface{} {
// 	// Get the type and value of the original data
// 	originalValue := reflect.ValueOf(data)
// 	originalType := originalValue.Type()

// 	// Ensure the input is a struct
// 	if originalValue.Kind() != reflect.Struct {
// 		panic("input must be a struct")
// 	}

// 	// Create a slice to store the fields of the new struct
// 	var fields []reflect.StructField

// 	// Iterate through the fields of the original struct
// 	for i := 0; i < originalType.NumField(); i++ {

// 		field := originalType.Field(i)
// 		fieldType := field.Type

// 		// If the field is a nested struct, handle it recursively
// 		if fieldType.Kind() == reflect.Struct {

// 			if fieldType.String() != "time.Time" {
// 				// Recursively modify nested structs

// 				subFields := EncryptObjectID(reflect.New(fieldType).Elem().Interface())

// 				// Append the field with the updated type of the nested struct
// 				fields = append(fields, reflect.StructField{
// 					Name: field.Name,
// 					Type: reflect.TypeOf(subFields), // Set the type of the nested struct
// 					Tag:  field.Tag,
// 				})

// 			} else {

// 				// Copy the time.Time field as is
// 				fields = append(fields, field)
// 			}

// 		} else if strings.HasSuffix(field.Name, "ID") && fieldType.Kind() == reflect.Uint {
// 			// Modify fields ending with "ID" of type uint to string

// 			fields = append(fields, reflect.StructField{
// 				Name: field.Name,
// 				Type: reflect.TypeOf(""), // Change type to string
// 				Tag:  field.Tag,
// 			})
// 		} else {
// 			// Otherwise, keep the original field type
// 			fields = append(fields, field)
// 		}
// 	}

// 	// Create a new struct type with the modified fields
// 	newStructType := reflect.StructOf(fields)
// 	newStruct := reflect.New(newStructType).Elem()

// 	// Copy field values from the original struct to the new struct
// 	for i := 0; i < originalValue.NumField(); i++ {
// 		originalField := originalValue.Field(i)
// 		newField := newStruct.Field(i)

// 		// Handle "ID" fields: Encrypt them
// 		if strings.HasSuffix(fields[i].Name, "ID") && fields[i].Type.Kind() == reflect.String {
// 			// Encrypt the ID field and set it
// 			newField.SetString(utils.EncryptID(uint(originalField.Uint())))
// 		} else if fields[i].Type.Kind() == reflect.Struct {
// 			// Recursively handle nested structs
// 			originalFieldValue := originalField.Interface()
// 			if reflect.TypeOf(originalFieldValue).String() != "time.Time" {
// 				encryptedValue := EncryptObjectID(originalFieldValue)

// 				// Ensure types match before setting
// 				if reflect.TypeOf(encryptedValue) == newField.Type() {
// 					newField.Set(reflect.ValueOf(encryptedValue))
// 				}
// 			} else {
// 				// Copy time.Time fields as is
// 				newField.Set(originalField)
// 			}
// 		} else {
// 			// Copy other fields as is
// 			newField.Set(originalField)
// 		}
// 	}

// 	// Return the new struct as an interface{}
// 	return newStruct.Interface()
// }

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
