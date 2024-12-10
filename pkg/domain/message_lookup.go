package domain

import "slash10k/pkg/models"

type MessageLookup interface {
	IsRegistrationMessage(messageId string) bool
	AddSetup(botSetup models.BotSetup)
}

type messageLookup struct {
	table map[string]any
}

var _ MessageLookup = (*messageLookup)(nil)

func NewMessageLookup(
	botSetups []models.BotSetup,
) *messageLookup {
	m := messageLookup{}
	m.table = make(map[string]any)
	for _, botSetup := range botSetups {
		m.table[botSetup.RegistrationMessageId] = struct{}{}
	}
	return &m
}

func (m *messageLookup) IsRegistrationMessage(messageId string) bool {
	_, ok := m.table[messageId]
	return ok
}

func (m *messageLookup) AddSetup(botSetup models.BotSetup) {
	m.table[botSetup.RegistrationMessageId] = struct{}{}
}
