package professional

type UpdateRequest struct {
	Email string `json:"email" validate:"email,required,max=50"`
}
