package post

import (
	"github.com/ardipermana59/go-template/internal/common/apperror"
	"gorm.io/gorm"
)

type Service interface {
	CreatePost(userID uint, dto CreatePostDTO) (*PostResponse, apperror.AppErrors)
	GetAllPosts() ([]PostResponse, apperror.AppErrors)
	GetPostByID(id uint) (*PostResponse, apperror.AppErrors)
	GetPostsByUserID(userID uint) ([]PostResponse, apperror.AppErrors)
	GetMyPosts(userID uint) ([]PostResponse, apperror.AppErrors)
	UpdatePost(id, userID uint, dto UpdatePostDTO) (*PostResponse, apperror.AppErrors)
	DeletePost(id, userID uint) apperror.AppErrors
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePost(userID uint, dto CreatePostDTO) (*PostResponse, apperror.AppErrors) {
	post := &Post{
		Title:   dto.Title,
		Content: dto.Content,
		UserID:  userID,
	}

	if err := s.repo.Create(post); err != nil {
		return nil, apperror.DatabaseError(err)
	}

	createdPost, err := s.repo.FindByID(post.ID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return createdPost.ToResponse(), nil
}

func (s *service) GetAllPosts() ([]PostResponse, apperror.AppErrors) {
	posts, err := s.repo.FindAll()
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	var responses []PostResponse
	for _, post := range posts {
		responses = append(responses, *post.ToResponse())
	}

	return responses, nil
}

func (s *service) GetPostByID(id uint) (*PostResponse, apperror.AppErrors) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.PostNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}
	return post.ToResponse(), nil
}

func (s *service) GetPostsByUserID(userID uint) ([]PostResponse, apperror.AppErrors) {
	posts, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	var responses []PostResponse
	for _, post := range posts {
		responses = append(responses, *post.ToResponse())
	}

	return responses, nil
}

func (s *service) GetMyPosts(userID uint) ([]PostResponse, apperror.AppErrors) {
	return s.GetPostsByUserID(userID)
}

func (s *service) UpdatePost(id, userID uint, dto UpdatePostDTO) (*PostResponse, apperror.AppErrors) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperror.PostNotFound()
		}
		return nil, apperror.DatabaseError(err)
	}

	if post.UserID != userID {
		return nil, apperror.OwnershipRequired()
	}

	if dto.Title != "" {
		post.Title = dto.Title
	}
	if dto.Content != "" {
		post.Content = dto.Content
	}

	if err := s.repo.Update(post); err != nil {
		return nil, apperror.DatabaseError(err)
	}

	updatedPost, err := s.repo.FindByID(id)
	if err != nil {
		return nil, apperror.DatabaseError(err)
	}

	return updatedPost.ToResponse(), nil
}

func (s *service) DeletePost(id, userID uint) apperror.AppErrors {
	post, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperror.PostNotFound()
		}
		return apperror.DatabaseError(err)
	}

	if post.UserID != userID {
		return apperror.OwnershipRequired()
	}

	if err := s.repo.Delete(id); err != nil {
		return apperror.DatabaseError(err)
	}

	return nil
}
