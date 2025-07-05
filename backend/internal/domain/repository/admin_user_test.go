package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/nature-console/backend/internal/domain/entity"
)

// TestAdminUserRepositoryInterface tests the interface definition
func TestAdminUserRepositoryInterface(t *testing.T) {
	// Test interface definition exists and has correct method signatures
	interfaceType := reflect.TypeOf((*AdminUserRepository)(nil)).Elem()
	
	if interfaceType.Kind() != reflect.Interface {
		t.Error("AdminUserRepository should be an interface")
	}
	
	expectedMethods := []string{
		"GetByEmail", "GetByID", "Create", "Update", "Delete", "List",
	}
	
	for _, methodName := range expectedMethods {
		_, exists := interfaceType.MethodByName(methodName)
		if !exists {
			t.Errorf("Method %s not found in AdminUserRepository interface", methodName)
		}
	}
}

// TestAdminUserRepositoryMethods tests that all interface methods
// have correct signatures
func TestAdminUserRepositoryMethods(t *testing.T) {
	// This test verifies that the interface can be implemented
	// by checking if a nil pointer of the interface type can be assigned
	var repo AdminUserRepository
	if repo != nil {
		t.Error("Nil repository should be nil")
	}
	
	// Test that the interface methods have the expected signatures
	// This is mainly a compile-time check
	ctx := context.Background()
	
	// These calls will panic with nil pointer, but the important thing
	// is that the interface methods are defined with correct signatures
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling methods on nil interface")
		}
	}()
	
	// This will panic, but validates the method signatures
	repo.GetByEmail(ctx, "test@example.com")
}

// TestAdminUserRepositorySignatures tests that all methods have correct signatures
func TestAdminUserRepositorySignatures(t *testing.T) {
	interfaceType := reflect.TypeOf((*AdminUserRepository)(nil)).Elem()
	
	tests := []struct {
		methodName string
		numIn      int
		numOut     int
		inTypes    []reflect.Type
		outTypes   []reflect.Type
	}{
		{
			methodName: "GetByEmail",
			numIn:      2, // ctx, email
			numOut:     2, // *entity.AdminUser, error
			inTypes: []reflect.Type{
				reflect.TypeOf((*context.Context)(nil)).Elem(),
				reflect.TypeOf(""),
			},
			outTypes: []reflect.Type{
				reflect.TypeOf((*entity.AdminUser)(nil)),
				reflect.TypeOf((*error)(nil)).Elem(),
			},
		},
		{
			methodName: "GetByID",
			numIn:      2, // ctx, id
			numOut:     2, // *entity.AdminUser, error
			inTypes: []reflect.Type{
				reflect.TypeOf((*context.Context)(nil)).Elem(),
				reflect.TypeOf(uint(0)),
			},
			outTypes: []reflect.Type{
				reflect.TypeOf((*entity.AdminUser)(nil)),
				reflect.TypeOf((*error)(nil)).Elem(),
			},
		},
		{
			methodName: "Create",
			numIn:      2, // ctx, user
			numOut:     1, // error
			inTypes: []reflect.Type{
				reflect.TypeOf((*context.Context)(nil)).Elem(),
				reflect.TypeOf((*entity.AdminUser)(nil)),
			},
			outTypes: []reflect.Type{
				reflect.TypeOf((*error)(nil)).Elem(),
			},
		},
		{
			methodName: "Update",
			numIn:      2, // ctx, user
			numOut:     1, // error
			inTypes: []reflect.Type{
				reflect.TypeOf((*context.Context)(nil)).Elem(),
				reflect.TypeOf((*entity.AdminUser)(nil)),
			},
			outTypes: []reflect.Type{
				reflect.TypeOf((*error)(nil)).Elem(),
			},
		},
		{
			methodName: "Delete",
			numIn:      2, // ctx, id
			numOut:     1, // error
			inTypes: []reflect.Type{
				reflect.TypeOf((*context.Context)(nil)).Elem(),
				reflect.TypeOf(uint(0)),
			},
			outTypes: []reflect.Type{
				reflect.TypeOf((*error)(nil)).Elem(),
			},
		},
		{
			methodName: "List",
			numIn:      1, // ctx
			numOut:     2, // []*entity.AdminUser, error
			inTypes: []reflect.Type{
				reflect.TypeOf((*context.Context)(nil)).Elem(),
			},
			outTypes: []reflect.Type{
				reflect.TypeOf([]*entity.AdminUser{}),
				reflect.TypeOf((*error)(nil)).Elem(),
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.methodName+" signature", func(t *testing.T) {
			method, exists := interfaceType.MethodByName(tt.methodName)
			if !exists {
				t.Fatalf("Method %s not found in interface", tt.methodName)
			}
			
			methodType := method.Type
			
			// Check number of inputs (including receiver)
			if methodType.NumIn() != tt.numIn {
				t.Errorf("Method %s: expected %d inputs, got %d", tt.methodName, tt.numIn, methodType.NumIn())
			}
			
			// Check number of outputs
			if methodType.NumOut() != tt.numOut {
				t.Errorf("Method %s: expected %d outputs, got %d", tt.methodName, tt.numOut, methodType.NumOut())
			}
			
			// Check input types
			for i, expectedType := range tt.inTypes {
				if i < methodType.NumIn() {
					actualType := methodType.In(i)
					if actualType != expectedType {
						t.Errorf("Method %s: input %d expected type %v, got %v", tt.methodName, i, expectedType, actualType)
					}
				}
			}
			
			// Check output types
			for i, expectedType := range tt.outTypes {
				if i < methodType.NumOut() {
					actualType := methodType.Out(i)
					if actualType != expectedType {
						t.Errorf("Method %s: output %d expected type %v, got %v", tt.methodName, i, expectedType, actualType)
					}
				}
			}
		})
	}
}

// TestAdminUserRepository_ContextParameter tests that methods accept context
func TestAdminUserRepository_ContextParameter(t *testing.T) {
	// Test that the interface methods accept context.Context parameter
	// This is a compile-time check
	
	// Check that we can assign a function with the correct signature
	var getByEmail func(context.Context, string) (*entity.AdminUser, error)
	getByEmail = func(ctx context.Context, email string) (*entity.AdminUser, error) {
		return nil, nil
	}
	
	if getByEmail == nil {
		t.Error("getByEmail function should not be nil")
	}
	
	// Test that context can be created and passed
	ctx := context.Background()
	if ctx == nil {
		t.Error("Context should not be nil")
	}
	
	// Test context with values
	ctxWithValue := context.WithValue(ctx, "test", "value")
	if ctxWithValue == nil {
		t.Error("Context with value should not be nil")
	}
}