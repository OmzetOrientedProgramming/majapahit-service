package customer

import "time"

type EditCustomerRequest struct {
	ID                int
	Name              string `json:"name"`
	ProfilePicture    string `json:"profile_picture"`
	DateOfBirth       time.Time
	DateOfBirthString string `json:"date_of_birth"`
	Gender            int    `json:"gender"`
}
