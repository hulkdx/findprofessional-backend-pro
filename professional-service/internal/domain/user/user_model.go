package user

type User struct {
	// TODO: convert it to int64
	ID           int     `json:"-"`
	FirstName    string  `json:"firstName,omitempty"`
	LastName     string  `json:"lastName,omitempty"`
	Email        string  `json:"email,omitempty"`
	ProfileImage *string `json:"profileImage,omitempty"`
}
