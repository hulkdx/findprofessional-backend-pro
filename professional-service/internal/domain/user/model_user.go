package user

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	Email        string `json:"email,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
}
