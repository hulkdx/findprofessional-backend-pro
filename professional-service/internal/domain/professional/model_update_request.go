package professional

type UpdateRequest struct {
	Email           *string `json:"email" validate:"omitempty,email,max=50"`
	FirstName       string  `json:"firstName" validate:"max=50"`
	LastName        string  `json:"lastName" validate:"max=50"`
	CoachType       string  `json:"coachType" validate:"max=50"`
	Price           *int64  `json:"priceNumber" validate:"omitempty"`
	PriceCurrency   *string `json:"priceCurrency" validate:"omitempty,max=50"`
	ProfileImageUrl *string `json:"profileImageUrl" validate:"omitempty,max=50"`
	Description     *string `json:"description" validate:"omitempty,max=255"`
	SessionPlatform *string `json:"sessionPlatform" validate:"omitempty,max=50"`
	SessionLink     *string `json:"sessionLink" validate:"omitempty,max=2048"`
}
