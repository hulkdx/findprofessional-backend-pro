package professional

type ProCreateRequest struct {
	Email           string  `json:"email" validate:"email,required,max=50"`
	Password        string  `json:"password" validate:"max=50"`
	FirstName       string  `json:"firstName" validate:"max=50"`
	LastName        string  `json:"lastName" validate:"max=50"`
	SessionPlatform *string `json:"sessionPlatform" validate:"omitempty,max=50"`
	SessionLink     *string `json:"sessionLink" validate:"omitempty,max=2048"`
	AboutMe         string  `json:"aboutMe" validate:"max=500"`
	CoachType       string  `json:"coachType" validate:"max=50"`
	Price           int64   `json:"price"`
	PriceCurrency   string  `json:"priceCurrency" validate:"max=50"`
}
