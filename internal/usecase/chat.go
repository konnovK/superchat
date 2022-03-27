package usecase

import (
	"github.com/konnovK/superchat/internal/entity"
)

type ChatDTO interface {
	GetAllChats() (entity.GetActiveChatsResponse, error)
}
type ChatContent struct {
}
