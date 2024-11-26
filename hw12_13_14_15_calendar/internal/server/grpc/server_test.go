package grpc_test

import (
	"context"
	"testing"
	"time"

	api "github.com/Dendyator/1/hw12_13_14_15_calendar/api/pb"           //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"      //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/server/grpc" //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage"     //nolint
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateEvent(event storage.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockStorage) UpdateEvent(id string, newEvent storage.Event) error {
	args := m.Called(id, newEvent)
	return args.Error(0)
}

func (m *MockStorage) DeleteEvent(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStorage) GetEvent(id string) (storage.Event, error) {
	args := m.Called(id)
	return args.Get(0).(storage.Event), args.Error(1)
}

func (m *MockStorage) ListEvents() ([]storage.Event, error) {
	args := m.Called()
	return args.Get(0).([]storage.Event), args.Error(1)
}

func TestCreateEvent(t *testing.T) {
	mockStorage := new(MockStorage)
	mockLogger := logger.New("info")
	server := grpc.NewGRPCServer(mockStorage, mockLogger)

	event := &api.Event{
		Id:          "1",
		Title:       "Test Event",
		Description: "This is a test",
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Add(2 * time.Hour).Unix(),
	}

	mockStorage.On("CreateEvent", mock.Anything).Return(nil)

	req := &api.CreateEventRequest{Event: event}
	res, err := server.CreateEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockStorage.AssertExpectations(t)
}

func TestUpdateEvent(t *testing.T) {
	mockStorage := new(MockStorage)
	mockLogger := logger.New("info")
	server := grpc.NewGRPCServer(mockStorage, mockLogger)

	event := &api.Event{
		Id:          "1",
		Title:       "Updated Event",
		Description: "This event is updated",
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Add(2 * time.Hour).Unix(),
	}

	mockStorage.On("UpdateEvent", "1", mock.Anything).Return(nil)

	req := &api.UpdateEventRequest{Id: "1", Event: event}
	res, err := server.UpdateEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockStorage.AssertExpectations(t)
}

func TestDeleteEvent(t *testing.T) {
	mockStorage := new(MockStorage)
	mockLogger := logger.New("info")
	server := grpc.NewGRPCServer(mockStorage, mockLogger)

	mockStorage.On("DeleteEvent", "1").Return(nil)

	req := &api.DeleteEventRequest{Id: "1"}
	res, err := server.DeleteEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockStorage.AssertExpectations(t)
}

func TestGetEvent(t *testing.T) {
	mockStorage := new(MockStorage)
	mockLogger := logger.New("info")
	server := grpc.NewGRPCServer(mockStorage, mockLogger)

	storedEvent := storage.Event{
		ID:          "1",
		Title:       "Test Event",
		Description: "This is a test",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(2 * time.Hour),
	}

	mockStorage.On("GetEvent", "1").Return(storedEvent, nil)

	req := &api.GetEventRequest{Id: "1"}
	res, err := server.GetEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, storedEvent.ID, res.Event.Id)
	mockStorage.AssertExpectations(t)
}

func TestListEvents(t *testing.T) {
	mockStorage := new(MockStorage)
	mockLogger := logger.New("info")
	server := grpc.NewGRPCServer(mockStorage, mockLogger)

	storedEvents := []storage.Event{
		{ID: "1", Title: "Event 1"},
		{ID: "2", Title: "Event 2"},
	}

	mockStorage.On("ListEvents").Return(storedEvents, nil)

	req := &api.ListEventsRequest{}
	res, err := server.ListEvents(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, len(storedEvents), len(res.Events))
	mockStorage.AssertExpectations(t)
}
