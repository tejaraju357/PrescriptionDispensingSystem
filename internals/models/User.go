package models

// User represents an application user
// @Description User can be an admin, doctor, or pharmacist
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // admin, doctor, pharmacist
	Password  string    `json:"password"`
}
