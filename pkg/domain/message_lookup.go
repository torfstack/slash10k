package domain

type MessageLookup interface {
	IsRegistrationMessage(messageId string) bool
	AddSetup(botSetup BotSetup)
}

type messageLookup struct {
	table map[string]any
}

var _ MessageLookup = (*messageLookup)(nil)

func NewMessageLookup(
	botSetups []BotSetup,
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

func (m *messageLookup) AddSetup(botSetup BotSetup) {
	m.table[botSetup.RegistrationMessageId] = struct{}{}
}
