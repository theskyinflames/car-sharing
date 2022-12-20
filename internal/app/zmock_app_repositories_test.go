// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app_test

import (
	"context"
	"github.com/google/uuid"
	"sync"
	"theskyinflames/car-sharing/internal/app"
	"theskyinflames/car-sharing/internal/domain"
)

// Ensure, that GroupsRepositoryMock does implement app.GroupsRepository.
// If this is not the case, regenerate this file with moq.
var _ app.GroupsRepository = &GroupsRepositoryMock{}

// GroupsRepositoryMock is a mock implementation of app.GroupsRepository.
//
//	func TestSomethingThatUsesGroupsRepository(t *testing.T) {
//
//		// make and configure a mocked app.GroupsRepository
//		mockedGroupsRepository := &GroupsRepositoryMock{
//			AddFunc: func(ctx context.Context, g domain.Group) error {
//				panic("mock out the Add method")
//			},
//			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (domain.Group, error) {
//				panic("mock out the FindByID method")
//			},
//			FindGroupsWithoutCarFunc: func(ctx context.Context) ([]domain.Group, error) {
//				panic("mock out the FindGroupsWithoutCar method")
//			},
//			RemoveAllFunc: func(ctx context.Context) error {
//				panic("mock out the RemoveAll method")
//			},
//			RemoveByIDFunc: func(ctx context.Context, ID uuid.UUID) error {
//				panic("mock out the RemoveByID method")
//			},
//			UpdateFunc: func(ctx context.Context, g domain.Group) error {
//				panic("mock out the Update method")
//			},
//		}
//
//		// use mockedGroupsRepository in code that requires app.GroupsRepository
//		// and then make assertions.
//
//	}
type GroupsRepositoryMock struct {
	// AddFunc mocks the Add method.
	AddFunc func(ctx context.Context, g domain.Group) error

	// FindByIDFunc mocks the FindByID method.
	FindByIDFunc func(ctx context.Context, ID uuid.UUID) (domain.Group, error)

	// FindGroupsWithoutCarFunc mocks the FindGroupsWithoutCar method.
	FindGroupsWithoutCarFunc func(ctx context.Context) ([]domain.Group, error)

	// RemoveAllFunc mocks the RemoveAll method.
	RemoveAllFunc func(ctx context.Context) error

	// RemoveByIDFunc mocks the RemoveByID method.
	RemoveByIDFunc func(ctx context.Context, ID uuid.UUID) error

	// UpdateFunc mocks the Update method.
	UpdateFunc func(ctx context.Context, g domain.Group) error

	// calls tracks calls to the methods.
	calls struct {
		// Add holds details about calls to the Add method.
		Add []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// G is the g argument value.
			G domain.Group
		}
		// FindByID holds details about calls to the FindByID method.
		FindByID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the ID argument value.
			ID uuid.UUID
		}
		// FindGroupsWithoutCar holds details about calls to the FindGroupsWithoutCar method.
		FindGroupsWithoutCar []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// RemoveAll holds details about calls to the RemoveAll method.
		RemoveAll []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// RemoveByID holds details about calls to the RemoveByID method.
		RemoveByID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the ID argument value.
			ID uuid.UUID
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// G is the g argument value.
			G domain.Group
		}
	}
	lockAdd                  sync.RWMutex
	lockFindByID             sync.RWMutex
	lockFindGroupsWithoutCar sync.RWMutex
	lockRemoveAll            sync.RWMutex
	lockRemoveByID           sync.RWMutex
	lockUpdate               sync.RWMutex
}

// Add calls AddFunc.
func (mock *GroupsRepositoryMock) Add(ctx context.Context, g domain.Group) error {
	callInfo := struct {
		Ctx context.Context
		G   domain.Group
	}{
		Ctx: ctx,
		G:   g,
	}
	mock.lockAdd.Lock()
	mock.calls.Add = append(mock.calls.Add, callInfo)
	mock.lockAdd.Unlock()
	if mock.AddFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.AddFunc(ctx, g)
}

// AddCalls gets all the calls that were made to Add.
// Check the length with:
//
//	len(mockedGroupsRepository.AddCalls())
func (mock *GroupsRepositoryMock) AddCalls() []struct {
	Ctx context.Context
	G   domain.Group
} {
	var calls []struct {
		Ctx context.Context
		G   domain.Group
	}
	mock.lockAdd.RLock()
	calls = mock.calls.Add
	mock.lockAdd.RUnlock()
	return calls
}

// FindByID calls FindByIDFunc.
func (mock *GroupsRepositoryMock) FindByID(ctx context.Context, ID uuid.UUID) (domain.Group, error) {
	callInfo := struct {
		Ctx context.Context
		ID  uuid.UUID
	}{
		Ctx: ctx,
		ID:  ID,
	}
	mock.lockFindByID.Lock()
	mock.calls.FindByID = append(mock.calls.FindByID, callInfo)
	mock.lockFindByID.Unlock()
	if mock.FindByIDFunc == nil {
		var (
			groupOut domain.Group
			errOut   error
		)
		return groupOut, errOut
	}
	return mock.FindByIDFunc(ctx, ID)
}

// FindByIDCalls gets all the calls that were made to FindByID.
// Check the length with:
//
//	len(mockedGroupsRepository.FindByIDCalls())
func (mock *GroupsRepositoryMock) FindByIDCalls() []struct {
	Ctx context.Context
	ID  uuid.UUID
} {
	var calls []struct {
		Ctx context.Context
		ID  uuid.UUID
	}
	mock.lockFindByID.RLock()
	calls = mock.calls.FindByID
	mock.lockFindByID.RUnlock()
	return calls
}

// FindGroupsWithoutCar calls FindGroupsWithoutCarFunc.
func (mock *GroupsRepositoryMock) FindGroupsWithoutCar(ctx context.Context) ([]domain.Group, error) {
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockFindGroupsWithoutCar.Lock()
	mock.calls.FindGroupsWithoutCar = append(mock.calls.FindGroupsWithoutCar, callInfo)
	mock.lockFindGroupsWithoutCar.Unlock()
	if mock.FindGroupsWithoutCarFunc == nil {
		var (
			groupsOut []domain.Group
			errOut    error
		)
		return groupsOut, errOut
	}
	return mock.FindGroupsWithoutCarFunc(ctx)
}

// FindGroupsWithoutCarCalls gets all the calls that were made to FindGroupsWithoutCar.
// Check the length with:
//
//	len(mockedGroupsRepository.FindGroupsWithoutCarCalls())
func (mock *GroupsRepositoryMock) FindGroupsWithoutCarCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockFindGroupsWithoutCar.RLock()
	calls = mock.calls.FindGroupsWithoutCar
	mock.lockFindGroupsWithoutCar.RUnlock()
	return calls
}

// RemoveAll calls RemoveAllFunc.
func (mock *GroupsRepositoryMock) RemoveAll(ctx context.Context) error {
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockRemoveAll.Lock()
	mock.calls.RemoveAll = append(mock.calls.RemoveAll, callInfo)
	mock.lockRemoveAll.Unlock()
	if mock.RemoveAllFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.RemoveAllFunc(ctx)
}

// RemoveAllCalls gets all the calls that were made to RemoveAll.
// Check the length with:
//
//	len(mockedGroupsRepository.RemoveAllCalls())
func (mock *GroupsRepositoryMock) RemoveAllCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockRemoveAll.RLock()
	calls = mock.calls.RemoveAll
	mock.lockRemoveAll.RUnlock()
	return calls
}

// RemoveByID calls RemoveByIDFunc.
func (mock *GroupsRepositoryMock) RemoveByID(ctx context.Context, ID uuid.UUID) error {
	callInfo := struct {
		Ctx context.Context
		ID  uuid.UUID
	}{
		Ctx: ctx,
		ID:  ID,
	}
	mock.lockRemoveByID.Lock()
	mock.calls.RemoveByID = append(mock.calls.RemoveByID, callInfo)
	mock.lockRemoveByID.Unlock()
	if mock.RemoveByIDFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.RemoveByIDFunc(ctx, ID)
}

// RemoveByIDCalls gets all the calls that were made to RemoveByID.
// Check the length with:
//
//	len(mockedGroupsRepository.RemoveByIDCalls())
func (mock *GroupsRepositoryMock) RemoveByIDCalls() []struct {
	Ctx context.Context
	ID  uuid.UUID
} {
	var calls []struct {
		Ctx context.Context
		ID  uuid.UUID
	}
	mock.lockRemoveByID.RLock()
	calls = mock.calls.RemoveByID
	mock.lockRemoveByID.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *GroupsRepositoryMock) Update(ctx context.Context, g domain.Group) error {
	callInfo := struct {
		Ctx context.Context
		G   domain.Group
	}{
		Ctx: ctx,
		G:   g,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	if mock.UpdateFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.UpdateFunc(ctx, g)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//
//	len(mockedGroupsRepository.UpdateCalls())
func (mock *GroupsRepositoryMock) UpdateCalls() []struct {
	Ctx context.Context
	G   domain.Group
} {
	var calls []struct {
		Ctx context.Context
		G   domain.Group
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}

// Ensure, that CarsRepositoryMock does implement app.CarsRepository.
// If this is not the case, regenerate this file with moq.
var _ app.CarsRepository = &CarsRepositoryMock{}

// CarsRepositoryMock is a mock implementation of app.CarsRepository.
//
//	func TestSomethingThatUsesCarsRepository(t *testing.T) {
//
//		// make and configure a mocked app.CarsRepository
//		mockedCarsRepository := &CarsRepositoryMock{
//			AddAllFunc: func(ctx context.Context, cars []domain.Car) error {
//				panic("mock out the AddAll method")
//			},
//			FindAllFunc: func(ctx context.Context) ([]domain.Car, error) {
//				panic("mock out the FindAll method")
//			},
//			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (domain.Car, error) {
//				panic("mock out the FindByID method")
//			},
//			RemoveAllFunc: func(ctx context.Context) error {
//				panic("mock out the RemoveAll method")
//			},
//			UpdateFunc: func(ctx context.Context, car domain.Car) error {
//				panic("mock out the Update method")
//			},
//		}
//
//		// use mockedCarsRepository in code that requires app.CarsRepository
//		// and then make assertions.
//
//	}
type CarsRepositoryMock struct {
	// AddAllFunc mocks the AddAll method.
	AddAllFunc func(ctx context.Context, cars []domain.Car) error

	// FindAllFunc mocks the FindAll method.
	FindAllFunc func(ctx context.Context) ([]domain.Car, error)

	// FindByIDFunc mocks the FindByID method.
	FindByIDFunc func(ctx context.Context, ID uuid.UUID) (domain.Car, error)

	// RemoveAllFunc mocks the RemoveAll method.
	RemoveAllFunc func(ctx context.Context) error

	// UpdateFunc mocks the Update method.
	UpdateFunc func(ctx context.Context, car domain.Car) error

	// calls tracks calls to the methods.
	calls struct {
		// AddAll holds details about calls to the AddAll method.
		AddAll []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cars is the cars argument value.
			Cars []domain.Car
		}
		// FindAll holds details about calls to the FindAll method.
		FindAll []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// FindByID holds details about calls to the FindByID method.
		FindByID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the ID argument value.
			ID uuid.UUID
		}
		// RemoveAll holds details about calls to the RemoveAll method.
		RemoveAll []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Car is the car argument value.
			Car domain.Car
		}
	}
	lockAddAll    sync.RWMutex
	lockFindAll   sync.RWMutex
	lockFindByID  sync.RWMutex
	lockRemoveAll sync.RWMutex
	lockUpdate    sync.RWMutex
}

// AddAll calls AddAllFunc.
func (mock *CarsRepositoryMock) AddAll(ctx context.Context, cars []domain.Car) error {
	callInfo := struct {
		Ctx  context.Context
		Cars []domain.Car
	}{
		Ctx:  ctx,
		Cars: cars,
	}
	mock.lockAddAll.Lock()
	mock.calls.AddAll = append(mock.calls.AddAll, callInfo)
	mock.lockAddAll.Unlock()
	if mock.AddAllFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.AddAllFunc(ctx, cars)
}

// AddAllCalls gets all the calls that were made to AddAll.
// Check the length with:
//
//	len(mockedCarsRepository.AddAllCalls())
func (mock *CarsRepositoryMock) AddAllCalls() []struct {
	Ctx  context.Context
	Cars []domain.Car
} {
	var calls []struct {
		Ctx  context.Context
		Cars []domain.Car
	}
	mock.lockAddAll.RLock()
	calls = mock.calls.AddAll
	mock.lockAddAll.RUnlock()
	return calls
}

// FindAll calls FindAllFunc.
func (mock *CarsRepositoryMock) FindAll(ctx context.Context) ([]domain.Car, error) {
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockFindAll.Lock()
	mock.calls.FindAll = append(mock.calls.FindAll, callInfo)
	mock.lockFindAll.Unlock()
	if mock.FindAllFunc == nil {
		var (
			carsOut []domain.Car
			errOut  error
		)
		return carsOut, errOut
	}
	return mock.FindAllFunc(ctx)
}

// FindAllCalls gets all the calls that were made to FindAll.
// Check the length with:
//
//	len(mockedCarsRepository.FindAllCalls())
func (mock *CarsRepositoryMock) FindAllCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockFindAll.RLock()
	calls = mock.calls.FindAll
	mock.lockFindAll.RUnlock()
	return calls
}

// FindByID calls FindByIDFunc.
func (mock *CarsRepositoryMock) FindByID(ctx context.Context, ID uuid.UUID) (domain.Car, error) {
	callInfo := struct {
		Ctx context.Context
		ID  uuid.UUID
	}{
		Ctx: ctx,
		ID:  ID,
	}
	mock.lockFindByID.Lock()
	mock.calls.FindByID = append(mock.calls.FindByID, callInfo)
	mock.lockFindByID.Unlock()
	if mock.FindByIDFunc == nil {
		var (
			carOut domain.Car
			errOut error
		)
		return carOut, errOut
	}
	return mock.FindByIDFunc(ctx, ID)
}

// FindByIDCalls gets all the calls that were made to FindByID.
// Check the length with:
//
//	len(mockedCarsRepository.FindByIDCalls())
func (mock *CarsRepositoryMock) FindByIDCalls() []struct {
	Ctx context.Context
	ID  uuid.UUID
} {
	var calls []struct {
		Ctx context.Context
		ID  uuid.UUID
	}
	mock.lockFindByID.RLock()
	calls = mock.calls.FindByID
	mock.lockFindByID.RUnlock()
	return calls
}

// RemoveAll calls RemoveAllFunc.
func (mock *CarsRepositoryMock) RemoveAll(ctx context.Context) error {
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockRemoveAll.Lock()
	mock.calls.RemoveAll = append(mock.calls.RemoveAll, callInfo)
	mock.lockRemoveAll.Unlock()
	if mock.RemoveAllFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.RemoveAllFunc(ctx)
}

// RemoveAllCalls gets all the calls that were made to RemoveAll.
// Check the length with:
//
//	len(mockedCarsRepository.RemoveAllCalls())
func (mock *CarsRepositoryMock) RemoveAllCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockRemoveAll.RLock()
	calls = mock.calls.RemoveAll
	mock.lockRemoveAll.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *CarsRepositoryMock) Update(ctx context.Context, car domain.Car) error {
	callInfo := struct {
		Ctx context.Context
		Car domain.Car
	}{
		Ctx: ctx,
		Car: car,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	if mock.UpdateFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.UpdateFunc(ctx, car)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//
//	len(mockedCarsRepository.UpdateCalls())
func (mock *CarsRepositoryMock) UpdateCalls() []struct {
	Ctx context.Context
	Car domain.Car
} {
	var calls []struct {
		Ctx context.Context
		Car domain.Car
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}
