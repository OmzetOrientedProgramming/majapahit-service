package customer

import "time"

// Profile consists of data shown in customer's profile page
type Profile struct {
	PhoneNumber 	string 		`json:"phone_number" db:"phone_number"`
	Name        	string 		`json:"name" db:"name"`
	Gender 			int 		`json:"gender" db:"gender"`
	DateOfBirth 	time.Time 	`json:"date_of_birth" db:"date_of_birth"`
	ProfilePicture 	string 		`json:"image" db:"image"`
}