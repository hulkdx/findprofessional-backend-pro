package professional

type CreateRequest struct {
	Email     string `json:"email" validate:"email,required,max=50"`
	Password  string `json:"password" validate:"max=50"`
	FirstName string `json:"firstName" validate:"max=50"`
	LastName  string `json:"lastName" validate:"max=50"`
	SkypeId   string `json:"skypeId" validate:"max=50"`
	AboutMe   string `json:"aboutMe" validate:"max=50"`
}
