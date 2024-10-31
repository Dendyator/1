package grpc

import (
	"context"
	"time"

	pb "github.com/Dendyator/1/hw12_13_14_15_calendar/api/pb"        //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/logger"  //nolint
	"github.com/Dendyator/1/hw12_13_14_15_calendar/internal/storage" //nolint
)

type Server struct {
	pb.UnimplementedEventServiceServer
	storage storage.Interface
	logg    *logger.Logger
}

func NewGRPCServer(storage storage.Interface, logg *logger.Logger) *Server {
	return &Server{storage: storage, logg: logg}
}

func (s *Server) CreateEvent(_ context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	s.logg.Info("Creating event: " + req.GetEvent().Title)
	event := storage.Event{
		ID:          req.GetEvent().Id,
		Title:       req.GetEvent().Title,
		Description: req.GetEvent().Description,
		StartTime:   time.Unix(req.GetEvent().StartTime, 0),
		EndTime:     time.Unix(req.GetEvent().EndTime, 0),
	}
	err := s.storage.CreateEvent(event)
	if err != nil {
		s.logg.Error("Failed to create event: " + err.Error())
	}
	return &pb.CreateEventResponse{}, err
}

func (s *Server) UpdateEvent(_ context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	s.logg.Info("Updating event ID: " + req.GetId())
	newEvent := storage.Event{
		ID:          req.GetEvent().Id,
		Title:       req.GetEvent().Title,
		Description: req.GetEvent().Description,
		StartTime:   time.Unix(req.GetEvent().StartTime, 0),
		EndTime:     time.Unix(req.GetEvent().EndTime, 0),
	}
	err := s.storage.UpdateEvent(req.GetId(), newEvent)
	if err != nil {
		s.logg.Error("Failed to update event: " + err.Error())
	}
	return &pb.UpdateEventResponse{}, err
}

func (s *Server) DeleteEvent(_ context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	s.logg.Info("Deleting event ID: " + req.GetId())
	err := s.storage.DeleteEvent(req.GetId())
	if err != nil {
		s.logg.Error("Failed to delete event: " + err.Error())
	}
	return &pb.DeleteEventResponse{}, err
}

func (s *Server) GetEvent(_ context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error) {
	s.logg.Info("Retrieving event ID: " + req.GetId())
	event, err := s.storage.GetEvent(req.GetId())
	if err != nil {
		s.logg.Error("Failed to get event: " + err.Error())
		return nil, err
	}
	return &pb.GetEventResponse{
		Event: &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime.Unix(),
			EndTime:     event.EndTime.Unix(),
		},
	}, nil
}

func (s *Server) ListEvents(_ context.Context, _ *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	s.logg.Info("Listing all events")
	events, err := s.storage.ListEvents()
	if err != nil {
		s.logg.Error("Failed to list events: " + err.Error())
		return nil, err
	}
	pbEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		pbEvents[i] = &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime.Unix(),
			EndTime:     event.EndTime.Unix(),
		}
	}
	return &pb.ListEventsResponse{Events: pbEvents}, nil
}
