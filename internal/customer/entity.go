package customer

import "time"

// EditCustomerRequest for handling request from client
type EditCustomerRequest struct {
	ID                int
	Name              string `json:"name"`
	ProfilePicture    string `json:"profile_picture"`
	DateOfBirth       time.Time
	DateOfBirthString string `json:"date_of_birth"`
	Gender            int    `json:"gender"`
}

// Profile consists of data shown in customer's profile page
type Profile struct {
	PhoneNumber    string    `json:"phone_number" db:"phone_number"`
	Name           string    `json:"name" db:"name"`
	Gender         int       `json:"gender" db:"gender"`
	DateOfBirth    time.Time `json:"date_of_birth" db:"date_of_birth"`
	ProfilePicture string    `json:"image" db:"image"`
}
