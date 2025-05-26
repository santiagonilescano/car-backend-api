package get_car_by_id

import (
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services" // To define the interface for the mock
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm" // For gorm.ErrRecordNotFound
)

// MockCarService is a mock implementation of services.CarService for testing.
type MockCarService struct {
	GetCarByIDFunc func(ctx context.Context, id uuid.UUID) (*entities.Car, error)
	// Unused methods for this specific test, but part of the interface
	CreateCarFunc func(ctx context.Context, car *entities.Car) (*entities.Car, error)
	GetCarsFunc   func(ctx context.Context) ([]*entities.Car, error)
	UpdateCarFunc func(ctx context.Context, car *entities.Car) (*entities.Car, error)
}

// Ensure MockCarService implements services.CarService
var _ services.CarService = (*MockCarService)(nil)

func (m *MockCarService) CreateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	if m.CreateCarFunc != nil {
		return m.CreateCarFunc(ctx, car)
	}
	panic("CreateCarFunc not implemented in mock")
}

func (m *MockCarService) GetCars(ctx context.Context) ([]*entities.Car, error) {
	if m.GetCarsFunc != nil {
		return m.GetCarsFunc(ctx)
	}
	panic("GetCarsFunc not implemented in mock")
}

func (m *MockCarService) GetCarByID(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
	if m.GetCarByIDFunc != nil {
		return m.GetCarByIDFunc(ctx, id)
	}
	panic("GetCarByIDFunc not implemented in mock for this test run")
}
func (m *MockCarService) UpdateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	if m.UpdateCarFunc != nil {
		return m.UpdateCarFunc(ctx, car)
	}
	panic("UpdateCarFunc not implemented in mock")
}

func TestGetCarByIDQueryHandler_Execute(t *testing.T) {
	ctx := context.Background()
	targetCarID := uuid.New()
	
	mockCarEntity := &entities.Car{
		ID:        targetCarID,
		ModelID:   uuid.New(),
		OwnerID:   uuid.New(),
		Year:      2023,
		Color:     "Blue",
		VIN:       "TESTVINGETBYID00",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Define the test case structure directly
	tests := []struct {
		name          string
		request       GetCarByIDRequest
		setupMock     func(mock *MockCarService, currentTestRequest GetCarByIDRequest)
		expectedCar   *entities.Car
		expectedError error
	}{
		{
			name:    "Success Case: Car Found",
			request: GetCarByIDRequest{ID: targetCarID},
			setupMock: func(mock *MockCarService, currentTestRequest GetCarByIDRequest) {
				mock.GetCarByIDFunc = func(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
					assert.Equal(t, currentTestRequest.ID, id, "Mock: GetCarByID received incorrect ID")
					return mockCarEntity, nil
				}
			},
			expectedCar:   mockCarEntity,
			expectedError: nil,
		},
		{
			name:    "Not Found Case: Car does not exist",
			request: GetCarByIDRequest{ID: targetCarID},
			setupMock: func(mock *MockCarService, currentTestRequest GetCarByIDRequest) {
				mock.GetCarByIDFunc = func(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
					assert.Equal(t, currentTestRequest.ID, id, "Mock: GetCarByID received incorrect ID")
					return nil, gorm.ErrRecordNotFound // Simulate GORM's not found error
				}
			},
			expectedCar:   nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:    "Service Error Case: Generic error from service",
			request: GetCarByIDRequest{ID: targetCarID},
			setupMock: func(mock *MockCarService, currentTestRequest GetCarByIDRequest) {
				mock.GetCarByIDFunc = func(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
					assert.Equal(t, currentTestRequest.ID, id, "Mock: GetCarByID received incorrect ID")
					return nil, errors.New("database connection error")
				}
			},
			expectedCar:   nil,
			expectedError: errors.New("database connection error"),
		},
	}
	
	for _, tt := range tests {
		currentTest := tt // Capture range variable by value
		t.Run(currentTest.name, func(t *testing.T) {
			t.Parallel() // Mark test for parallel execution
			mockService := &MockCarService{}
			if currentTest.setupMock != nil {
				// Pass currentTest.request to setupMock
				currentTest.setupMock(mockService, currentTest.request)
			}

			handler := NewGetCarByIDQueryHandler(mockService)
			resultCar, err := handler.Execute(currentTest.request, ctx)

			assert.Equal(t, currentTest.expectedCar, resultCar)
			if currentTest.expectedError != nil {
				assert.EqualError(t, err, currentTest.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
