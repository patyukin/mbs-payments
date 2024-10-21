package usecase

type Registry interface {
}

type UseCase struct {
	registry Registry
}

func New(registry Registry) *UseCase {
	return &UseCase{
		registry: registry,
	}
}
