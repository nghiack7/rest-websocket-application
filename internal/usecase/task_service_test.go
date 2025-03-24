package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/personal/task-management/internal/delivery/rest/dtos"
	"github.com/personal/task-management/internal/domain/task"
	"github.com/personal/task-management/internal/domain/user"
	"github.com/personal/task-management/internal/mocks"
)

type TaskServiceTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	taskRepo    *mocks.MockTaskRepository
	userRepo    *mocks.MockUserRepository
	taskService TaskService
}

func (suite *TaskServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.taskRepo = mocks.NewMockTaskRepository(suite.ctrl)
	suite.userRepo = mocks.NewMockUserRepository(suite.ctrl)
	suite.taskService = NewTaskService(suite.taskRepo, suite.userRepo)
}

func (suite *TaskServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *TaskServiceTestSuite) TestCreateTask_Success() {
	// Test data
	creatorID := uuid.New()
	assigneeID := uuid.New()
	input := dtos.CreateTaskInput{
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     time.Now().Add(24 * time.Hour),
		AssigneeID:  assigneeID,
		CreatorID:   creatorID,
	}

	creator := &user.User{
		ID:   creatorID,
		Role: user.Employer,
	}

	assignee := &user.User{
		ID:   assigneeID,
		Role: user.Employee,
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByID(gomock.Any(), creatorID).
		Return(creator, nil)

	suite.userRepo.EXPECT().
		GetByID(gomock.Any(), assigneeID).
		Return(assignee, nil)

	suite.taskRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *task.Task) error {
			suite.Equal(input.Title, t.Title)
			suite.Equal(input.Description, t.Description)
			suite.Equal(input.DueDate, t.DueDate)
			suite.Equal(input.AssigneeID, t.AssigneeID)
			suite.Equal(input.CreatorID, t.CreatorID)
			suite.Equal(task.StatusPending, t.Status)
			return nil
		})

	// Call the service method
	result, err := suite.taskService.CreateTask(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(input.Title, result.Title)
	suite.Equal(input.Description, result.Description)
	suite.Equal(input.DueDate, result.DueDate)
	suite.Equal(input.AssigneeID, result.AssigneeID)
	suite.Equal(input.CreatorID, result.CreatorID)
	suite.Equal(task.StatusPending, result.Status)
}

func (suite *TaskServiceTestSuite) TestGetTask_Success() {
	// Test data
	taskID := uuid.New()
	requesterID := uuid.New()
	input := dtos.GetTaskInput{
		TaskID:      taskID,
		RequesterID: requesterID,
	}

	requester := &user.User{
		ID:   requesterID,
		Role: user.Employee,
	}

	expectedTask := &task.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      task.StatusPending,
		AssigneeID:  requesterID,
		CreatorID:   uuid.New(),
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByID(gomock.Any(), requesterID).
		Return(requester, nil)

	suite.taskRepo.EXPECT().
		GetByID(gomock.Any(), taskID).
		Return(expectedTask, nil)

	// Call the service method
	result, err := suite.taskService.GetTask(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedTask.ID, result.ID)
	suite.Equal(expectedTask.Title, result.Title)
	suite.Equal(expectedTask.Description, result.Description)
	suite.Equal(expectedTask.DueDate, result.DueDate)
	suite.Equal(expectedTask.Status, result.Status)
	suite.Equal(expectedTask.AssigneeID, result.AssigneeID)
	suite.Equal(expectedTask.CreatorID, result.CreatorID)
}

func (suite *TaskServiceTestSuite) TestUpdateTaskStatus_Success() {
	// Test data
	taskID := uuid.New()
	userID := uuid.New()
	input := dtos.UpdateTaskStatusInput{
		TaskID:    taskID,
		UserID:    userID,
		NewStatus: task.StatusInProgress,
	}

	user := &user.User{
		ID:   userID,
		Role: user.Employee,
	}

	existingTask := &task.Task{
		ID:          taskID,
		Title:       "Test Task",
		Description: "Test Description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      task.StatusPending,
		AssigneeID:  userID,
		CreatorID:   uuid.New(),
	}

	// Set up expectations
	suite.taskRepo.EXPECT().
		GetByID(gomock.Any(), taskID).
		Return(existingTask, nil)

	suite.userRepo.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(user, nil)

	suite.taskRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *task.Task) error {
			suite.Equal(task.StatusInProgress, t.Status)
			return nil
		})

	// Call the service method
	result, err := suite.taskService.UpdateTaskStatus(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(task.StatusInProgress, result.Status)
}

func (suite *TaskServiceTestSuite) TestDeleteTask_Success() {
	// Test data
	taskID := uuid.New()
	requesterID := uuid.New()
	input := dtos.DeleteTaskInput{
		TaskID:      taskID,
		RequesterID: requesterID,
	}

	requester := &user.User{
		ID:   requesterID,
		Role: user.Employer,
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByID(gomock.Any(), requesterID).
		Return(requester, nil)

	suite.taskRepo.EXPECT().
		Delete(gomock.Any(), taskID).
		Return(nil)

	// Call the service method
	err := suite.taskService.DeleteTask(context.Background(), input)

	// Assertions
	suite.NoError(err)
}

func (suite *TaskServiceTestSuite) TestGetTasksWithFilter_Success() {
	// Test data
	userID := uuid.New()
	input := dtos.GetTasksWithFilterInput{
		UserID: userID,
		Filter: dtos.TaskFilter{
			Status:    task.StatusPending,
			Limit:     10,
			Offset:    0,
			SortBy:    "created_at",
			SortOrder: "desc",
		},
	}

	user := &user.User{
		ID:   userID,
		Role: user.Employer,
	}

	expectedTasks := []*task.Task{
		{
			ID:          uuid.New(),
			Title:       "Task 1",
			Description: "Description 1",
			DueDate:     time.Now().Add(24 * time.Hour),
			Status:      task.StatusPending,
			AssigneeID:  uuid.New(),
			CreatorID:   uuid.New(),
		},
		{
			ID:          uuid.New(),
			Title:       "Task 2",
			Description: "Description 2",
			DueDate:     time.Now().Add(48 * time.Hour),
			Status:      task.StatusPending,
			AssigneeID:  uuid.New(),
			CreatorID:   uuid.New(),
		},
	}

	// Set up expectations
	suite.userRepo.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(user, nil)

	suite.taskRepo.EXPECT().
		List(gomock.Any(), gomock.Any()).
		Return(expectedTasks, nil)

	// Call the service method
	result, err := suite.taskService.GetTasksWithFilter(context.Background(), input)

	// Assertions
	suite.NoError(err)
	suite.NotNil(result)
	suite.Len(result, 2)
	suite.Equal(expectedTasks[0].ID, result[0].ID)
	suite.Equal(expectedTasks[1].ID, result[1].ID)
}

func TestTaskServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceTestSuite))
}
