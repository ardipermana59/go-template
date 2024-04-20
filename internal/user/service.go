package user

import (
	"github.com/ardipermana59/go-template/internal/auth"
	"github.com/ardipermana59/go-template/internal/common/apperror"
	"gorm.io/gorm"
)

type Service interface {
	Register(dto RegisterDTO) (*UserResponse, apperror.AppErrors)
	Login(dto LoginDTO) (*LoginResponse, apperror.AppErrors)
	GetAllUsers() ([]UserResponse, apperror.AppErrors)
	GetUserByID(id uint) (*UserResponse, apperror.AppErrors)
	GetProfile(id uint) (*UserResponse, apperror.AppErrors)
	UpdateUser(id uint, dto UpdateUserDTO) (*UserResponse, apperror.AppErrors)
	ChangePassword(id uint, dto ChangePasswordDTO) apperror.AppErrors
	DeleteUser(id uint) apperror.AppErrors
}

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) Service {
	return &service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *service) Register(dto RegisterDTO) (*UserResponse, apperror.AppErrors) {
	if s.repo.EmailExists(dto.Email) {
		return nil, apperror.EmailAlreadyExists()
	}

	user := &User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Role:     "user",
	}

	if err := user.HashPassword(); err != nil {
		return nil, apperror.DatabaseError(err)
	}

	if err := s.repo.Create(user); err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return user.ToResponse(), nil
}

func (s *service) Login(dto LoginDTO) (*LoginResponse, apperror.AppErrors) {
	user, err := s.repo.FindByEmail(dto.Email)
	if err != nil {
		return nil, apperror.InvalidCredentials()
	}

	if !user.CheckPassword(dto.Password) {
		return nil, apperror.InvalidCredentials()
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return &LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

func (s *service) GetAllUsers() ([]UserResponse, apperror.AppErrors) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	var responses []UserResponse
	for _, user := range users {
		responses = append(responses, *user.ToResponse())
	}

	return responses, nil
}

func (s *service) GetUserByID(id uint) (*UserResponse, apperror.AppErrors) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.UserNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}
	return user.ToResponse(), nil
}

func (s *service) GetProfile(id uint) (*UserResponse, apperror.AppErrors) {
	return s.GetUserByID(id)
}

func (s *service) UpdateUser(id uint, dto UpdateUserDTO) (*UserResponse, apperror.AppErrors) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.UserNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.Email != "" {
		existingUser, _ := s.repo.FindByEmail(dto.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, apperror.EmailAlreadyExists()
		}
		user.Email = dto.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return user.ToResponse(), nil
}

func (s *service) ChangePassword(id uint, dto ChangePasswordDTO) apperror.AppErrors {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.UserNotFound()
		}
		return apperror.DatabaseError(err)
	}

	if !user.CheckPassword(dto.OldPassword) {
		return apperror.OldPasswordIncorrect()
	}

	user.Password = dto.NewPassword
	if err := user.HashPassword(); err != nil {
		return apperror.DatabaseError(err)
	}

	if err := s.repo.Update(user); err != nil {
		return apperror.DatabaseError(err)
	}

	return nil
}

func (s *service) DeleteUser(id uint) apperror.AppErrors {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.UserNotFound()
		}
		return apperror.DatabaseError(err)
	}

	if err := s.repo.Delete(user.ID); err != nil {
		return apperror.DatabaseError(err)
	}

	return nil
}
