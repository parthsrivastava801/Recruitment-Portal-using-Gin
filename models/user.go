package models

import (
	"errors"
	"sync"
)

// Define roles.
type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleRecruiter  Role = "recruiter"
	RoleApplicant  Role = "applicant"
)

// User represents a user in the system.
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
	// Approved is relevant for recruiter accounts;
	// for a new user (default Applicant) it is ignored.
	Approved bool `json:"approved"`
}

// --- In-Memory Store for demonstration purposes ---
var (
	users          = make(map[string]*User) // key by email
	lastUserID int = 0
	mu         sync.Mutex
)

// CreateUser registers a user using their OAuth information.
// If the user already exists, it returns the existing user.
func CreateUser(email, name string) (*User, error) {
	mu.Lock()
	defer mu.Unlock()

	if user, exists := users[email]; exists {
		return user, nil // already registered
	}

	lastUserID++
	user := &User{
		ID:       lastUserID,
		Email:    email,
		Name:     name,
		Role:     RoleApplicant, // default role is Applicant
		Approved: false,         // recruiter approval is pending if role changes later
	}

	users[email] = user
	return user, nil
}

// GetUser retrieves a user by email.
func GetUser(email string) (*User, error) {
	mu.Lock()
	defer mu.Unlock()

	user, exists := users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GrantRecruiterStatus grants a user recruiter status by changing their role
// and marking them as approved. Only a Super Admin should call this.
func GrantRecruiterStatus(email string) (*User, error) {
	mu.Lock()
	defer mu.Unlock()

	user, exists := users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	user.Role = RoleRecruiter
	user.Approved = true
	return user, nil
}
