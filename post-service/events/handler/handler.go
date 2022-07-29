package events

import (
	"github.com/project1/post-service/config"
	"github.com/project1/post-service/pkg/logger"
	"github.com/project1/post-service/storage"
	pb "github.com/project1/post-service/genproto"
)

type EventHandler struct {
	config config.Config
	storage storage.IStorage
	log logger.Logger
}

func NewEventHandlerFunc(config config.Config, storage storage.IStorage, log logger.Logger) *EventHandler {
	return &EventHandler{
		config: config,
		storage: storage,
		log: log,
	}
}

func (h *EventHandler) Handler(value []byte) error {
	var user pb.User

	err := user.Unmarshal(value)

	if err != nil {
		return err
	}

	err = h.storage.Post().CreatePostUser(&user)
	if err != nil {
		return err
	}

	return nil
}