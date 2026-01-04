package user

type User struct {
	ID           int64   `json:"-"`
	FirstName    string  `json:"firstName,omitempty"`
	LastName     string  `json:"lastName,omitempty"`
	Email        string  `json:"email,omitempty"`
	ProfileImage *string `json:"profileImage,omitempty"`
}
