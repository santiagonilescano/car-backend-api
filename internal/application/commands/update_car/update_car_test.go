package update_car

import (
	"car-service/cmd/api/mediator"
	"car-service/internal/domain/entities"
	"car-service/internal/domain/services" // To define the interface for the mock
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm" // For gorm.ErrRecordNotFound
)

// MockCarService is a mock implementation of services.CarService for testing.
type MockCarService struct {
	CreateCarFunc  func(ctx context.Context, car *entities.Car) (*entities.Car, error)
	GetCarsFunc    func(ctx context.Context) ([]*entities.Car, error)
	GetCarByIDFunc func(ctx context.Context, id uuid.UUID) (*entities.Car, error)
	UpdateCarFunc  func(ctx context.Context, car *entities.Car) (*entities.Car, error)
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
	panic("GetCarByIDFunc not implemented in mock")
}

func (m *MockCarService) UpdateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	if m.UpdateCarFunc != nil {
		return m.UpdateCarFunc(ctx, car)
	}
	panic("UpdateCarFunc not implemented in mock")
}

// Helper functions for creating pointers to values
func strPtr(s string) *string { return &s }
func intPtr(i int) *int    { return &i }
func boolPtr(b bool) *bool   { return &b }
func uuidPtr(u uuid.UUID) *uuid.UUID { return &u }


func TestUpdateCarCommandHandler_Validate(t *testing.T) {
	carID := uuid.New()
	ginCtx, _ := gin.CreateTestContext(nil) 
	cmdCtx := &mediator.CommandContext{}    

	// Define test case structure for Validate
	type ValidateTestCase struct {
		name             string
		request          UpdateCarRequest
		expectedErrCount int
		expectedFields   map[string]string 
	}

	tests := []ValidateTestCase{
		{
			name: "Valid request - update color",
			request: UpdateCarRequest{
				ID:    carID,
				Color: strPtr("Blue"),
			},
			expectedErrCount: 0,
		},
		{
			name: "Invalid VIN - too short",
			request: UpdateCarRequest{
				ID:  carID,
				VIN: strPtr("SHORTVIN"),
			},
			expectedErrCount: 1,
			expectedFields:  map[string]string{"vin": "El VIN debe tener 17 caracteres"},
		},
		{
			name: "Invalid VIN - empty",
			request: UpdateCarRequest{
				ID:  carID,
				VIN: strPtr(""),
			},
			expectedErrCount: 1,
			expectedFields:  map[string]string{"vin": "El VIN no puede estar vacío si se proporciona"},
		},
		{
			name: "Invalid Year - too old",
			request: UpdateCarRequest{
				ID:   carID,
				Year: intPtr(1800),
			},
			expectedErrCount: 1,
			expectedFields:  map[string]string{"year": fmt.Sprintf("El año debe estar entre 1900 y %d", time.Now().Year()+1)},
		},
		{
			name: "Invalid Year - too new",
			request: UpdateCarRequest{
				ID:   carID,
				Year: intPtr(time.Now().Year() + 2),
			},
			expectedErrCount: 1,
			expectedFields:  map[string]string{"year": fmt.Sprintf("El año debe estar entre 1900 y %d", time.Now().Year()+1)},
		},
		{
			name: "No update fields provided in body",
			request: UpdateCarRequest{
				ID: carID, 
			},
			expectedErrCount: 1,
			expectedFields:  map[string]string{"requestBody": "Al menos un campo debe ser proporcionado para la actualización"},
		},
		{
			name: "Missing ID in request struct (path param)", 
			request: UpdateCarRequest{
				Color: strPtr("Blue"), 
			},
			expectedErrCount: 1,
			expectedFields:  map[string]string{"id": "El ID del auto es requerido en la URL"},
		},
	}

	for _, tt := range tests {
		currentTest := tt 
		t.Run(currentTest.name, func(t *testing.T) {
			t.Parallel()
			mockService := &MockCarService{} 
			handler := NewUpdateCarCommandHandler(mockService)

			validationErrors := handler.Validate(currentTest.request, ginCtx, cmdCtx)
			assert.Len(t, validationErrors, currentTest.expectedErrCount)

			if currentTest.expectedErrCount > 0 {
				foundErrors := make(map[string]string)
				for _, err := range validationErrors {
					foundErrors[err.Field] = err.Message
				}
				for field, msgSubstring := range currentTest.expectedFields {
					assert.Contains(t, foundErrors[field], msgSubstring, "Expected error for field '%s' not found or message mismatch", field)
				}
			}
		})
	}
}


func TestUpdateCarCommandHandler_Execute(t *testing.T) {
	baseCtx := context.Background() 
	
	baseTime := time.Now().Add(-24 * time.Hour) 
	carID := uuid.New()
	originalModelID := uuid.New()
	originalOwnerID := uuid.New()
	newModelID := uuid.New() 

	baseCarEntity := &entities.Car{
		ID:        carID,
		ModelID:   originalModelID,
		OwnerID:   originalOwnerID,
		Year:      2020,
		Color:     "Red",
		VIN:       "ORIGINALVIN123456",
		Active:    true,
		CreatedAt: baseTime,
		UpdatedAt: baseTime,
	}

	// Define test case structure for Execute
	type ExecuteTestCase struct {
		name           string
		request        UpdateCarRequest
		setupMock      func(mock *MockCarService, currentRequest UpdateCarRequest) 
		expectedCar    *entities.Car 
		expectedErrStr string        
		verify         func(t *testing.T, resultCar *entities.Car, originalCar *entities.Car, request UpdateCarRequest)
	}

	tests := []ExecuteTestCase{
		{
			name: "Success Case - Partial Update (Color and Year)",
			request: UpdateCarRequest{
				ID:    carID,
				Color: strPtr("Blue"),
				Year:  intPtr(2022),
			},
			setupMock: func(mock *MockCarService, currentRequest UpdateCarRequest) {
				mock.GetCarByIDFunc = func(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
					assert.Equal(t, currentRequest.ID, id)
					carCopy := *baseCarEntity 
					return &carCopy, nil
				}
				mock.UpdateCarFunc = func(ctx context.Context, carToUpdate *entities.Car) (*entities.Car, error) {
					assert.Equal(t, currentRequest.ID, carToUpdate.ID)
					assert.Equal(t, *currentRequest.Color, carToUpdate.Color)
					assert.Equal(t, *currentRequest.Year, carToUpdate.Year)
					assert.Equal(t, baseCarEntity.ModelID, carToUpdate.ModelID) 
					assert.Equal(t, baseCarEntity.VIN, carToUpdate.VIN)         
					carToUpdate.UpdatedAt = time.Now() 
					return carToUpdate, nil
				}
			},
			expectedCar: &entities.Car{ 
				ID:      carID,
				ModelID: originalModelID,
				OwnerID: originalOwnerID,
				Year:    2022, 
				Color:   "Blue", 
				VIN:     "ORIGINALVIN123456",
				Active:  true,
			},
			verify: func(t *testing.T, resultCar *entities.Car, originalCar *entities.Car, request UpdateCarRequest) {
				assert.NotNil(t, resultCar)
				assert.True(t, resultCar.UpdatedAt.After(originalCar.UpdatedAt), "UpdatedAt should be more recent")
			},
		},
		{
			name: "Success Case - Update ModelID",
			request: UpdateCarRequest{
				ID:      carID,
				ModelID: uuidPtr(newModelID),
			},
			setupMock: func(mock *MockCarService, currentRequest UpdateCarRequest) {
				mock.GetCarByIDFunc = func(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
					carCopy := *baseCarEntity
					return &carCopy, nil
				}
				mock.UpdateCarFunc = func(ctx context.Context, carToUpdate *entities.Car) (*entities.Car, error) {
					assert.Equal(t, *currentRequest.ModelID, carToUpdate.ModelID)
					carToUpdate.UpdatedAt = time.Now()
					return carToUpdate, nil
				}
			},
			expectedCar: &entities.Car{
				ID:      carID,
				ModelID: newModelID, 
				OwnerID: originalOwnerID,
				Year:    2020,
				Color:   "Red",
				VIN:     "ORIGINALVIN123456",
				Active:  true,
			},
		},
		{
			name:    "Car Not Found Case during GetCarByID",
			request: UpdateCarRequest{ID: carID, Color: strPtr("Green")},
			setupMock: func(mock *MockCarService, currentRequest UpdateCarRequest) {
				mock.GetCarByIDFunc = func(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
					return nil, gorm.ErrRecordNotFound
				}
			},
			expectedErrStr: gorm.ErrRecordNotFound.Error(),
		},
		{
			name:    "Update Fails - GetCarByID returns generic error",
			request: UpdateCarRequest{ID: carID, Color: strPtr("Green")},
			setupMock: func(mock *MockCarService, currentRequest UpdateCarRequest) {
				mock.GetCarByIDFunc = func(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
					return nil, errors.New("db connection error")
				}
			},
			expectedErrStr: "db connection error",
		},
		{
			name:    "Update Fails - Service UpdateCar returns error",
			request: UpdateCarRequest{ID: carID, Color: strPtr("Green")},
			setupMock: func(mock *MockCarService, currentRequest UpdateCarRequest) {
				mock.GetCarByIDFunc = func(ctx context.Context, id uuid.UUID) (*entities.Car, error) {
					carCopy := *baseCarEntity
					return &carCopy, nil
				}
				mock.UpdateCarFunc = func(ctx context.Context, carToUpdate *entities.Car) (*entities.Car, error) {
					return nil, errors.New("failed to update in db")
				}
			},
			expectedErrStr: "failed to update in db",
		},
	}


	for _, tt := range tests {
		currentTest := tt 
		t.Run(currentTest.name, func(t *testing.T) {
			t.Parallel()
			mockService := &MockCarService{}
			if currentTest.setupMock != nil {
				currentTest.setupMock(mockService, currentTest.request)
			}

			handler := NewUpdateCarCommandHandler(mockService)
			cmdCtx := context.WithValue(baseCtx, "test_name", currentTest.name) 

			resultCar, err := handler.Execute(currentTest.request, &cmdCtx)

			if currentTest.expectedErrStr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), currentTest.expectedErrStr)
				assert.Nil(t, resultCar)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resultCar)
				assert.Equal(t, currentTest.expectedCar.ID, resultCar.ID)
				assert.Equal(t, currentTest.expectedCar.ModelID, resultCar.ModelID)
				assert.Equal(t, currentTest.expectedCar.OwnerID, resultCar.OwnerID)
				assert.Equal(t, currentTest.expectedCar.Year, resultCar.Year)
				assert.Equal(t, currentTest.expectedCar.Color, resultCar.Color)
				assert.Equal(t, currentTest.expectedCar.VIN, resultCar.VIN)
				assert.Equal(t, currentTest.expectedCar.Active, resultCar.Active)
				
				if currentTest.verify != nil {
					currentTest.verify(t, resultCar, baseCarEntity, currentTest.request)
				}
			}
		})
	}
}
