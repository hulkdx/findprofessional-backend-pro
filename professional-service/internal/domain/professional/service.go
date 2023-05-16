package professional

type Service interface {
	Repository() Repository
	FindAllProfessional() ([]Professional, error)
}

type serviceImpl struct {
	repository Repository
}

func (s *serviceImpl) Repository() Repository {
	return s.repository
}

func (s *serviceImpl) FindAllProfessional() ([]Professional, error) {
	return s.repository.FindAll("id", "email")
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository}
}
