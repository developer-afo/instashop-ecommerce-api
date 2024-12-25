package user_service

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	userDto "github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/lib/helper"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	"github.com/developer-afo/instashop-ecommerce-api/repository"
	userRepository "github.com/developer-afo/instashop-ecommerce-api/repository/user"
)

var (
	UserRoleCustomer = "customer"
	UserRoleAdmin    = "admin"
)

type userService struct {
	userRepository userRepository.UserRepositoryInterface
}

type UserServiceInterface interface {
	CreateUser(dto userDto.UserDTO) (userDto.UserDTO, error)
	FindAllUsers(pageable repository.Pageable) ([]userDto.UserDTO, repository.Pagination, error)
	FindUserById(userId string) (userDto.UserDTO, error)
	FindUserByEmail(email string) (userDto.UserDTO, error)
	FindUserByReferralCode(referralCode string) (userDto.UserDTO, error)
	UpdateUser(dto userDto.UserDTO) (userDto.UserDTO, error)
	DeleteUser(uuid uuid.UUID) error
	ConvertToDTO(user models.User) (userDto userDto.UserDTO)
	ConvertToModel(userDto userDto.UserDTO) (user models.User)
}

func NewUserService(userRepository userRepository.UserRepositoryInterface) UserServiceInterface {
	return &userService{userRepository: userRepository}
}

func (service *userService) ConvertToDTO(user models.User) (userDto userDto.UserDTO) {

	userDto.ID = user.ID
	userDto.FirstName = user.FirstName
	userDto.LastName = user.LastName
	userDto.Email = user.Email
	userDto.IsEmailVerified = user.IsEmailVerified
	userDto.Password = user.Password
	userDto.Role = user.Role
	userDto.CreatedAt = user.CreatedAt
	userDto.UpdatedAt = user.UpdatedAt
	userDto.DeletedAt = user.DeletedAt.Time

	return userDto
}

func (service *userService) ConvertToModel(userDto userDto.UserDTO) (user models.User) {

	user.ID = userDto.ID
	user.Email = userDto.Email
	user.FirstName = userDto.FirstName
	user.LastName = userDto.LastName
	user.IsEmailVerified = userDto.IsEmailVerified
	user.Password = userDto.Password
	user.Role = userDto.Role
	user.CreatedAt = userDto.CreatedAt
	user.UpdatedAt = userDto.UpdatedAt
	user.DeletedAt.Time = userDto.DeletedAt

	return user
}

// CreateUser implements UserServiceInterface.
func (service *userService) CreateUser(userDtoArg userDto.UserDTO) (userDto.UserDTO, error) {

	userDtoArg.ReferralCode = helper.GenerateRandomString(10)

	user := service.ConvertToModel(userDtoArg)

	// check if user already exists return error
	_, err := service.userRepository.FindUserByEmail(user.Email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return userDto.UserDTO{}, errors.New("user already exists")
	}

	newRecord, err := service.userRepository.Create(user)

	return service.ConvertToDTO(newRecord), err
}

// FindAllUsers implements UserServiceInterface.
func (service *userService) FindAllUsers(pageable repository.Pageable) ([]userDto.UserDTO, repository.Pagination, error) {
	var userDtos []userDto.UserDTO

	users, pagination, err := service.userRepository.FindAllUsers(pageable)
	for _, user := range users {
		userDtos = append(userDtos, service.ConvertToDTO(user))
	}

	return userDtos, pagination, err
}

// FindUserById implements UserServiceInterface.
func (service *userService) FindUserById(userId string) (userDto.UserDTO, error) {

	_userId, err := uuid.Parse(userId)
	if err != nil {
		return userDto.UserDTO{}, err
	}

	user, err := service.userRepository.FindUserById(_userId)

	return service.ConvertToDTO(user), err
}

// FindUserByEmail implements UserServiceInterface.
func (service *userService) FindUserByEmail(email string) (userDto.UserDTO, error) {

	user, err := service.userRepository.FindUserByEmail(email)

	return service.ConvertToDTO(user), err
}

// FindUserByReferralCode implements UserServiceInterface.
func (service *userService) FindUserByReferralCode(referralCode string) (userDto.UserDTO, error) {

	user, err := service.userRepository.FindUserByReferralCode(referralCode)

	return service.ConvertToDTO(user), err
}

// UpdateUser implements UserServiceInterface.
func (service *userService) UpdateUser(userDto userDto.UserDTO) (userDto.UserDTO, error) {

	user := service.ConvertToModel(userDto)

	updatedRecord, err := service.userRepository.UpdateUser(user)

	return service.ConvertToDTO(updatedRecord), err
}

// DeleteUser implements UserServiceInterface.
func (service *userService) DeleteUser(uuid uuid.UUID) error {
	return service.userRepository.DeleteUser(uuid)
}
