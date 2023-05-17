package user

type Service interface {
	IsAuthenticated(auth string) bool
}

type serviceImpl struct {
}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) IsAuthenticated(auth string) bool {
	return true
}
