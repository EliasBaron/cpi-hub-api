package apperror

import (
	"fmt"
	"log"
)

// ExampleUsage demonstrates how to use the apperror package
func ExampleUsage() {
	// Example 1: Creating a NotFound error
	notFoundErr := NewNotFound("User not found", "user_id: 123", "user_service.go:GetUser")
	fmt.Printf("NotFound Error: %s\n", notFoundErr.Error())

	// Example 2: Creating an InvalidData error
	invalidDataErr := NewInvalidData("Invalid email format", "invalid_email@", "user_handler.go:Create")
	fmt.Printf("InvalidData Error: %s\n", invalidDataErr.Error())

	// Example 3: Creating an InternalServer error
	internalErr := NewInternalServer("Database connection failed", "connection timeout", "db_service.go:Connect")
	fmt.Printf("InternalServer Error: %s\n", internalErr.Error())

	// Example 4: Getting status code and message
	statusCode, message := StatusCodeAndMessage(notFoundErr)
	fmt.Printf("Status Code: %d, Message: %s\n", statusCode, message)

	// Example 5: Checking error type
	if Is(notFoundErr, NotFound) {
		log.Println("This is a NotFound error")
	}

	if Is(invalidDataErr, InvalidData) {
		log.Println("This is an InvalidData error")
	}
}
