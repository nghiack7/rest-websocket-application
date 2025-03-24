package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/personal/task-management/internal/delivery/rest/dtos"
	"github.com/personal/task-management/internal/domain/user"
	"github.com/personal/task-management/internal/mocks"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	userRepo    *mocks.MockUserRepository
	hasher      *mocks.MockHasher
	jwtService  *mocks.MockJWTTokenServicer
	userService UserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.userRepo = mocks.NewMockUserRepository(suite.ctrl)
	suite.hasher = mocks.NewMockHasher(suite.ctrl)
	suite.jwtService = mocks.NewMockJWTTokenServicer(suite.ctrl)
	suite.userService = NewUserService(suite.userRepo, suite.hasher, suite.jwtService)
}

func (suite *UserServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *UserServiceTestSuite) TestRegisterUser_Success() {
	// Test data
	input := dtos.RegisterUserInput{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "employee",
	}
	hashedPassword := "hashed_password"
	userID := uuid.New()

	// Set up expectations
	suite.hasher.EXPECT().
		HashPassword(input.Password).
		Return(hashedPassword, nil)

	suite.userRepo.EXPECT().
		GetByEmail(gomock.Any(), input.Email).
		Return(nil, user.ErrUserNotFound)

	suite.userRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *user.User) error {
			suite.Equal(input.Email, u.Email)
			suite.Equal(input.Name, u.Name)
			suite.Equal(hashedPassword, u.Password)
			u.ID = userID
			return nil
		})

	// Call the service method
	result, err := suite.userService.RegisterUser(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(userID, result.ID)
	suite.Equal(input.Email, result.Email)
	suite.Equal(input.Name, result.Name)
	suite.Equal("employee", result.Role)
}

func (suite *UserServiceTestSuite) TestRegisterUser_EmailExists() {
	// Test data
	input := dtos.RegisterUserInput{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     "employee",
	}

	existingUser := &user.User{
		ID:    uuid.New(),
		Email: input.Email,
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByEmail(gomock.Any(), input.Email).
		Return(existingUser, nil)

	// Call the service method
	result, err := suite.userService.RegisterUser(context.Background(), input)

	// Assertions
	suite.Error(err)
	suite.Equal(user.ErrEmailExists, err)
	suite.Nil(result)
}

func (suite *UserServiceTestSuite) TestLogin_Success() {
	// Test data
	input := dtos.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	storedUser := &user.User{
		ID:       uuid.New(),
		Email:    input.Email,
		Password: "hashed_password",
		Name:     "Test User",
		Role:     user.Employee,
	}

	expectedToken := "jwt_token"

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByEmail(gomock.Any(), input.Email).
		Return(storedUser, nil)

	suite.hasher.EXPECT().
		ComparePasswords(storedUser.Password, input.Password).
		Return(true)

	suite.jwtService.EXPECT().
		GenerateToken(storedUser.ID, storedUser.Email, storedUser.Role.String()).
		Return(expectedToken, nil)

	// Call the service method
	result, err := suite.userService.Login(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedToken, result.AuthToken)
	suite.Equal(storedUser.Email, result.User.Email)
	suite.Equal(storedUser.Name, result.User.Name)
}

func (suite *UserServiceTestSuite) TestLogin_InvalidCredentials() {
	// Test data
	input := dtos.LoginInput{
		Email:    "test@example.com",
		Password: "wrong_password",
	}

	storedUser := &user.User{
		ID:       uuid.New(),
		Email:    input.Email,
		Password: "hashed_password",
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByEmail(gomock.Any(), input.Email).
		Return(storedUser, nil)

	suite.hasher.EXPECT().
		ComparePasswords(storedUser.Password, input.Password).
		Return(false)

	// Call the service method
	result, err := suite.userService.Login(context.Background(), input)

	// Assertions
	suite.Error(err)
	suite.Equal(ErrInvalidCredentials, err)
	suite.Nil(result)
}

func (suite *UserServiceTestSuite) TestGetUser_Success() {
	// Test data
	userID := uuid.New()
	input := dtos.GetUserInput{
		ID: &userID,
	}

	expectedUser := &user.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(expectedUser, nil)

	// Call the service method
	result, err := suite.userService.GetUser(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedUser.ID, result.ID)
	suite.Equal(expectedUser.Email, result.Email)
	suite.Equal(expectedUser.Name, result.Name)
}

func (suite *UserServiceTestSuite) TestUpdateUser_Success() {
	// Test data
	userID := uuid.New()
	newName := "Updated Name"
	newPassword := "new_password"
	hashedPassword := "hashed_new_password"

	input := dtos.UpdateUserInput{
		ID:       userID,
		Name:     &newName,
		Password: &newPassword,
	}

	existingUser := &user.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Old Name",
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(existingUser, nil)

	suite.hasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.userRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *user.User) error {
			suite.Equal(newName, u.Name)
			suite.Equal(hashedPassword, u.Password)
			return nil
		})

	// Call the service method
	result, err := suite.userService.UpdateUser(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(newName, result.Name)
	suite.Equal(hashedPassword, result.Password)
}

func (suite *UserServiceTestSuite) TestListUsers_Success() {
	// Test data
	input := dtos.ListUsersInput{
		Offset: 0,
		Limit:  10,
	}

	expectedUsers := []*user.User{
		{
			ID:    uuid.New(),
			Email: "user1@example.com",
			Name:  "User 1",
		},
		{
			ID:    uuid.New(),
			Email: "user2@example.com",
			Name:  "User 2",
		},
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		List(gomock.Any(), input.Offset, input.Limit).
		Return(expectedUsers, nil)

	// Call the service method
	result, err := suite.userService.ListUsers(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Len(result, 2)
	suite.Equal(expectedUsers[0].ID, result[0].ID)
	suite.Equal(expectedUsers[1].ID, result[1].ID)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
