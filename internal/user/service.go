package user

import (
	"context"
	"osipovPetRestApi/pkg/logging"
)

type service struct {
	storage Storage
	logger  *logging.Logger
}

func (s *service) Create(ctx context.Context, dto Dto) (u User, err error) {
	//TODO for next one
	return
}
