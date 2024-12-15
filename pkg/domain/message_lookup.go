package domain

import (
	"slash10k/pkg/models"
	"slash10k/pkg/utils"
)

type MessageLookup interface {
	IsRegistrationMessage(messageId string) bool
	AddSetup(botSetup models.BotSetup)
}

type messageLookup struct {
	set utils.SyncSet[string]
}

var _ MessageLookup = (*messageLookup)(nil)

func NewMessageLookup(
	botSetups []models.BotSetup,
) *messageLookup {
	m := messageLookup{set: utils.NewSyncSet[string]()}
	for _, botSetup := range botSetups {
		m.set.Add(botSetup.RegistrationMessageId)
	}
	return &m
}

func (m *messageLookup) IsRegistrationMessage(messageId string) bool {
	return m.set.Contains(messageId)
}

func (m *messageLookup) AddSetup(botSetup models.BotSetup) {
	m.set.Add(botSetup.RegistrationMessageId)
}
