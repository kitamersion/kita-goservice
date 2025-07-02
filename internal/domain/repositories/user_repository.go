package repositories

import (
	"github.com/google/uuid"
	"github.com/kitamersion/go-goservice/internal/domain/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.UserEntity) error
	GetByID(id uuid.UUID) (*entities.UserEntity, error)
	GetByEmail(email string) (*entities.UserEntity, error)
	Update(user *entities.UserEntity) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*entities.UserEntity, error)
	Count() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *entities.UserEntity) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*entities.UserEntity, error) {
	var user entities.UserEntity
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*entities.UserEntity, error) {
	var user entities.UserEntity
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *entities.UserEntity) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.UserEntity{}, "id = ?", id).Error
}

func (r *userRepository) List(limit, offset int) ([]*entities.UserEntity, error) {
	var users []*entities.UserEntity
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entities.UserEntity{}).Count(&count).Error
	return count, err
}
