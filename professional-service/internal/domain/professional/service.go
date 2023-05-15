package professional

type Service interface {
	Repository() Repository
}

type serviceImpl struct {
	repository Repository
}

func (s *serviceImpl) Repository() Repository {
	return s.repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository}
}
