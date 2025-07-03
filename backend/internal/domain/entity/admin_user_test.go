package entity

import (
	"encoding/json"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestAdminUser_TableName(t *testing.T) {
	adminUser := AdminUser{}
	expected := "admin_users"
	
	if tableName := adminUser.TableName(); tableName != expected {
		t.Errorf("Expected table name '%s', got '%s'", expected, tableName)
	}
}

func TestAdminUser_JSONSerialization(t *testing.T) {
	now := time.Now()
	adminUser := AdminUser{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword123",
		Name:         "Test User",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Test JSON marshaling
	jsonBytes, err := json.Marshal(adminUser)
	if err != nil {
		t.Fatalf("Failed to marshal AdminUser to JSON: %v", err)
	}

	// Verify password hash is not included in JSON
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check that password hash is excluded (json:"-" tag)
	if _, exists := jsonMap["password_hash"]; exists {
		t.Error("Password hash should not be included in JSON serialization")
	}

	// Check that other fields are included
	expectedFields := []string{"id", "email", "name", "created_at", "updated_at"}
	for _, field := range expectedFields {
		if _, exists := jsonMap[field]; !exists {
			t.Errorf("Expected field '%s' to be present in JSON", field)
		}
	}

	// Test JSON unmarshaling
	var unmarshaled AdminUser
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to AdminUser: %v", err)
	}

	// Verify fields (password hash should be empty since it's not in JSON)
	if unmarshaled.ID != adminUser.ID {
		t.Errorf("Expected ID %d, got %d", adminUser.ID, unmarshaled.ID)
	}
	if unmarshaled.Email != adminUser.Email {
		t.Errorf("Expected Email %s, got %s", adminUser.Email, unmarshaled.Email)
	}
	if unmarshaled.Name != adminUser.Name {
		t.Errorf("Expected Name %s, got %s", adminUser.Name, unmarshaled.Name)
	}
	if unmarshaled.PasswordHash != "" {
		t.Error("Password hash should be empty after JSON unmarshaling")
	}
}

func TestAdminUser_StructTags(t *testing.T) {
	// Test that the struct has the expected GORM tags through reflection
	// This is more of a compile-time check, but we can verify the struct is properly defined
	
	tests := []struct {
		name        string
		create      func() AdminUser
		expectValid bool
	}{
		{
			name: "Valid admin user",
			create: func() AdminUser {
				return AdminUser{
					Email:        "test@example.com",
					PasswordHash: "hashedpassword",
					Name:         "Test User",
				}
			},
			expectValid: true,
		},
		{
			name: "Empty admin user",
			create: func() AdminUser {
				return AdminUser{}
			},
			expectValid: true, // Struct creation should always succeed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := tt.create()
			
			// Basic struct validation - ensure we can create and access fields
			if user.TableName() != "admin_users" {
				t.Error("TableName method should always return 'admin_users'")
			}
			
			// Test that we can set and get all fields
			user.ID = 1
			user.Email = "updated@example.com"
			user.PasswordHash = "newhash"
			user.Name = "Updated Name"
			user.CreatedAt = time.Now()
			user.UpdatedAt = time.Now()
			user.DeletedAt = gorm.DeletedAt{}
			
			if user.ID != 1 {
				t.Error("Failed to set/get ID field")
			}
			if user.Email != "updated@example.com" {
				t.Error("Failed to set/get Email field")
			}
			if user.PasswordHash != "newhash" {
				t.Error("Failed to set/get PasswordHash field")
			}
			if user.Name != "Updated Name" {
				t.Error("Failed to set/get Name field")
			}
		})
	}
}

func TestAdminUser_DefaultValues(t *testing.T) {
	adminUser := AdminUser{}
	
	// Test default values for a new AdminUser
	if adminUser.ID != 0 {
		t.Errorf("Expected default ID to be 0, got %d", adminUser.ID)
	}
	if adminUser.Email != "" {
		t.Errorf("Expected default Email to be empty, got '%s'", adminUser.Email)
	}
	if adminUser.PasswordHash != "" {
		t.Errorf("Expected default PasswordHash to be empty, got '%s'", adminUser.PasswordHash)
	}
	if adminUser.Name != "" {
		t.Errorf("Expected default Name to be empty, got '%s'", adminUser.Name)
	}
	if !adminUser.CreatedAt.IsZero() {
		t.Error("Expected default CreatedAt to be zero time")
	}
	if !adminUser.UpdatedAt.IsZero() {
		t.Error("Expected default UpdatedAt to be zero time")
	}
}

func TestAdminUser_TimeFields(t *testing.T) {
	now := time.Now()
	adminUser := AdminUser{
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	// Test that time fields work correctly
	if adminUser.CreatedAt != now {
		t.Error("CreatedAt field not set correctly")
	}
	if adminUser.UpdatedAt != now {
		t.Error("UpdatedAt field not set correctly")
	}
	
	// Test that we can update timestamps
	later := now.Add(time.Hour)
	adminUser.UpdatedAt = later
	
	if adminUser.UpdatedAt != later {
		t.Error("Failed to update UpdatedAt field")
	}
	if adminUser.CreatedAt != now {
		t.Error("CreatedAt should not change when UpdatedAt is modified")
	}
}

func TestAdminUser_DeletedAt(t *testing.T) {
	adminUser := AdminUser{}
	
	// Test default DeletedAt (should be zero value)
	if adminUser.DeletedAt.Valid {
		t.Error("DeletedAt should not be valid by default")
	}
	
	// Test setting DeletedAt
	now := time.Now()
	adminUser.DeletedAt = gorm.DeletedAt{
		Time:  now,
		Valid: true,
	}
	
	if !adminUser.DeletedAt.Valid {
		t.Error("DeletedAt should be valid after setting")
	}
	if adminUser.DeletedAt.Time != now {
		t.Error("DeletedAt time not set correctly")
	}
}

func TestAdminUser_PointerMethods(t *testing.T) {
	// Test that TableName method works with both value and pointer receivers
	adminUser := AdminUser{}
	adminUserPtr := &AdminUser{}
	
	if adminUser.TableName() != "admin_users" {
		t.Error("TableName should work with value receiver")
	}
	if adminUserPtr.TableName() != "admin_users" {
		t.Error("TableName should work with pointer receiver")
	}
}

func TestAdminUser_FieldTypes(t *testing.T) {
	adminUser := AdminUser{
		ID:           123,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	// Test that fields can handle their expected types
	if adminUser.ID != 123 {
		t.Error("ID field should handle uint type")
	}
	
	// Test that string fields can handle various inputs
	testStrings := []string{
		"",
		"simple",
		"email@domain.com",
		"a very long name with spaces and special characters!@#$%",
		"unicode: 日本語 ñoël français",
	}
	
	for _, str := range testStrings {
		adminUser.Email = str
		adminUser.Name = str
		adminUser.PasswordHash = str
		
		if adminUser.Email != str {
			t.Errorf("Email field failed to store string: %s", str)
		}
		if adminUser.Name != str {
			t.Errorf("Name field failed to store string: %s", str)
		}
		if adminUser.PasswordHash != str {
			t.Errorf("PasswordHash field failed to store string: %s", str)
		}
	}
}