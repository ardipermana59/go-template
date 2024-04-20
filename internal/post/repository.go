package post

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(post *Post) error
	FindAll() ([]Post, error)
	FindByID(id uint) (*Post, error)
	FindByUserID(userID uint) ([]Post, error)
	Update(post *Post) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(post *Post) error {
	return r.db.Create(post).Error
}

func (r *repository) FindAll() ([]Post, error) {
	var posts []Post
	err := r.db.Preload("User").Find(&posts).Error
	return posts, err
}

func (r *repository) FindByID(id uint) (*Post, error) {
	var post Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) FindByUserID(userID uint) ([]Post, error) {
	var posts []Post
	err := r.db.Preload("User").Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}

func (r *repository) Update(post *Post) error {
	return r.db.Save(post).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Post{}, id).Error
}
