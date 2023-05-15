package professional

type Service interface {
	Repository() Repository
	GetAllProfessionals() ([]Professional, error)
}

type serviceImpl struct {
	repository Repository
}

func (s *serviceImpl) Repository() Repository {
	return s.repository
}

func (s *serviceImpl) GetAllProfessionals() ([]Professional, error) {
	return s.repository.FindAll()
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository}
}
