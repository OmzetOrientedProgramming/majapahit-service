package user

import "time"

// Model represent user table on database, which is the parent of CustomerModel
type Model struct {
	ID              int       `db:"id"`
	PhoneNumber     string    `db:"phone_number"`
	Name            string    `db:"name"`
	Status          int       `db:"status"`
	FirebaseLocalID string    `db:"firebase_local_id"`
	Email           string    `db:"email"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
