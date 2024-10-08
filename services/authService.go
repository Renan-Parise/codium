package services

import (
	"github.com/Renan-Parise/codium-auth/entities"
	"github.com/Renan-Parise/codium-auth/errors"
	"github.com/Renan-Parise/codium-auth/repositories"
	"github.com/Renan-Parise/codium-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, error)
	Register(user entities.User) error
	Update(ID int, user entities.User) error
	DeactivateAccount(ID int) error
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.NewServiceError("authentication failed because user does not exist")
	}

	if !user.Active {
		return "", errors.NewServiceError("authentication failed because account is deactivated")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.NewServiceError("authentication failed because password is incorrect")
	}

	return utils.GenerateToken(user.ID)
}

func (s *authService) Register(user entities.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	_, err := s.userRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.NewServiceError("user already exists. please login or use another email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewServiceError("failed to hash password. please try again")
	}

	user.Password = string(hashedPassword)

	err = s.userRepo.Create(user)
	if err != nil {
		return errors.NewServiceError("failed to register user. please try again")
	}

	return nil
}

func (s *authService) Update(ID int, user entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewServiceError("failed to hash password. please try again")
	}

	user.Password = string(hashedPassword)

	err = s.userRepo.Update(ID, user)
	if err != nil {
		return errors.NewServiceError("failed to update user. please try again")
	}

	return nil
}

func (s *authService) DeactivateAccount(ID int) error {
	err := s.userRepo.DeactivateUser(ID)
	if err != nil {
		return errors.NewServiceError("failed to deactivate account")
	}
	return nil
}
