package userusecase

import (
	"avitostazhko/internal/domain"
	"avitostazhko/internal/infrastructure/repository"
	"context"
	"fmt"
	"time"
)

type UserUseCase struct {
	userRepository *repository.UserRepository
}

func NewUserUseCase(userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepository: userRepository}
}

func (uc *UserUseCase) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := uc.userRepository.GetUserByID(opCtx, userID)
	if err != nil {
		return nil, fmt.Errorf("userRepository: GetUserByID: %w", err)
	}

	if user == nil {
		return nil, &domain.ErrNotFound{
			Source: "user",
			ID:     userID,
		}
	}

	user.IsActive = isActive
	user, err = uc.userRepository.UpdateUser(opCtx, user)
	if err != nil {
		return nil, fmt.Errorf("userRepository: UpdateUser: %w", err)
	}

	return user, nil
}
