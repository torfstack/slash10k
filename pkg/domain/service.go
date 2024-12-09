package domain

import (
	"context"
	"errors"
	"fmt"
	"slash10k/pkg/convert/fromdb"
	"slash10k/pkg/db"
	"slash10k/pkg/models"
	sqlc "slash10k/sql/gen"
)

type Service interface {
	AddPlayer(ctx context.Context, discordId string, discordName string, guildId string, nick string) error
	DeletePlayer(ctx context.Context, discordId string, guildId string) error
	GetAllPlayers(ctx context.Context, guildId string) ([]models.Player, error)
	GetPlayer(ctx context.Context, discordId string, guildId string) (*models.Player, error)

	AddDebt(ctx context.Context, discordId string, guildId string, amount int64, reason string) error
	SubDebt(ctx context.Context, discordId string, guildId string, amount int64) error

	SetBotSetup(
		ctx context.Context,
		guildId string,
		channelId string,
		registrationMessageId string,
		debtsMessageId string,
	) error
	GetBotSetup(ctx context.Context, guildId string) (*models.BotSetup, error)
	GetAllBotSetups(ctx context.Context) ([]models.BotSetup, error)
	DeleteBotSetup(ctx context.Context, guildId string) error
}

var (
	ErrDatabase            = errors.New("some database error occurred")
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrPlayerDoesNotExist  = errors.New("player does not exist")

	ErrBotSetupAlreadyExists = errors.New("bot-setup already exists")
	ErrBotSetupDoesNotExist  = errors.New("bot-setup does not exist")
)

type service struct {
	db db.Database
}

var _ Service = (*service)(nil)

func NewSlashTenK(db db.Database) Service {
	return &service{db: db}
}

func (s service) AddPlayer(
	ctx context.Context,
	discordId string,
	discordName string,
	guildId string,
	nick string,
) error {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	tx, err := conn.StartTransaction(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	doesAlreadyExist, err := tx.Queries().DoesPlayerExist(
		ctx,
		sqlc.DoesPlayerExistParams{DiscordID: discordId, GuildID: guildId},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	if doesAlreadyExist {
		return fmt.Errorf("%w: %s(%s)@%s", ErrPlayerAlreadyExists, discordName, discordId, guildId)
	}

	_, err = tx.Queries().AddPlayer(
		ctx, sqlc.AddPlayerParams{
			DiscordID:   discordId,
			DiscordName: discordName,
			GuildID:     guildId,
			Name:        nick,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	return nil
}

func (s service) DeletePlayer(ctx context.Context, discordId string, guildId string) error {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	tx, err := conn.StartTransaction(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	doesExist, err := tx.Queries().DoesPlayerExist(
		ctx,
		sqlc.DoesPlayerExistParams{DiscordID: discordId, GuildID: guildId},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	if !doesExist {
		return fmt.Errorf("%w: %s@%s", ErrPlayerDoesNotExist, discordId, guildId)
	}

	id, err := tx.Queries().GetIdOfPlayer(
		ctx, sqlc.GetIdOfPlayerParams{
			DiscordID: discordId,
			GuildID:   guildId,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	err = tx.Queries().DeletePlayer(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	return nil
}

func (s service) GetAllPlayers(ctx context.Context, guildId string) ([]models.Player, error) {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	allPlayers, err := conn.Queries().GetAllPlayers(ctx, guildId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	return fromdb.FromAllPlayers(allPlayers), nil
}

func (s service) GetPlayer(ctx context.Context, discordId string, guildId string) (*models.Player, error) {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	player, err := conn.Queries().GetPlayer(
		ctx, sqlc.GetPlayerParams{
			DiscordID: discordId,
			GuildID:   guildId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	res := fromdb.FromPlayerWithDebt(player)
	return &res, nil
}

func (s service) AddDebt(ctx context.Context, discordId string, guildId string, amount int64, reason string) error {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	tx, err := conn.StartTransaction(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	queries := tx.Queries()

	player, err := queries.GetPlayer(
		ctx, sqlc.GetPlayerParams{
			DiscordID: discordId,
			GuildID:   guildId,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	currentPlayer := fromdb.FromPlayerWithDebt(player)

	newAmount := currentPlayer.Debt.Amount + amount

	err = queries.SetDebt(
		ctx, sqlc.SetDebtParams{
			Amount: newAmount,
			UserID: currentPlayer.Id,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	return nil
}

func (s service) SubDebt(ctx context.Context, discordId string, guildId string, amount int64) error {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	tx, err := conn.StartTransaction(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	queries := tx.Queries()

	player, err := queries.GetPlayer(
		ctx, sqlc.GetPlayerParams{
			DiscordID: discordId,
			GuildID:   guildId,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	currentPlayer := fromdb.FromPlayerWithDebt(player)

	newAmount := currentPlayer.Debt.Amount - amount
	if newAmount < 0 {
		newAmount = 0
	}

	err = queries.SetDebt(
		ctx, sqlc.SetDebtParams{
			Amount: newAmount,
			UserID: currentPlayer.Id,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	return nil
}

func (s service) SetBotSetup(
	ctx context.Context,
	guildId string,
	channelId string,
	registrationMessageId string,
	debtsMessageId string,
) error {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	tx, err := conn.StartTransaction(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	doesExist, err := tx.Queries().DoesBotSetupExist(ctx, guildId)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	if doesExist {
		return ErrBotSetupAlreadyExists
	}

	_, err = tx.Queries().PutBotSetup(
		ctx, sqlc.PutBotSetupParams{
			GuildID:               guildId,
			ChannelID:             channelId,
			DebtsMessageID:        debtsMessageId,
			RegistrationMessageID: registrationMessageId,
		},
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	return nil
}

func (s service) GetBotSetup(ctx context.Context, guildId string) (*models.BotSetup, error) {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	tx, err := conn.StartTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	doesExist, err := tx.Queries().DoesBotSetupExist(ctx, guildId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	if !doesExist {
		return nil, ErrBotSetupDoesNotExist
	}

	botSetup, err := tx.Queries().GetBotSetup(ctx, guildId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	res := fromdb.FromBotSetup(botSetup)
	return &res, nil
}

func (s service) GetAllBotSetups(ctx context.Context) ([]models.BotSetup, error) {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	botSetups, err := conn.Queries().GetAllBotSetups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	return fromdb.FromBotSetups(botSetups), nil
}

func (s service) DeleteBotSetup(ctx context.Context, guildId string) error {
	conn, err := s.db.Connect(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}
	defer conn.Close(ctx)

	doesExist, err := conn.Queries().DoesBotSetupExist(ctx, guildId)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	if !doesExist {
		return ErrBotSetupDoesNotExist
	}

	err = conn.Queries().DeleteBotSetup(ctx, guildId)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDatabase, err)
	}

	return nil
}
