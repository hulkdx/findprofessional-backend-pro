package professional

type Service interface {
	FindAllProfessional() ([]Professional, error)
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository}
}

func (s *serviceImpl) FindAllProfessional() ([]Professional, error) {
	return s.repository.FindAll("ID", "Email")
}
