package command

import (
	"context"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/rs/zerolog/log"
	"slash10k/pkg/domain"
)

func AddDebt(discord *state.State, service domain.Service) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		name := options.Find("name").String()
		amount, err := options.Find("amount").IntValue()
		if err != nil || amount < 0 {
			log.Error().Msgf("could not parse amount: %v, err: %s", amount, err)
			return ephemeralMessage("Amount needs to be a non-negative number!")
		}
		reason := options.Find("reason").String()

		err = service.AddDebt(ctx, data.Event.User.ID.String(), data.Event.GuildID.String(), amount, reason)
		if err != nil {
			log.Error().Msgf("could not add debt: %s", err)
			return ephemeralMessage("Could not update debt")
		}

		updateDebtsMessage(ctx, discord, service, data.Event.GuildID.String())

		return ephemeralMessage(fmt.Sprintf("Added %v to %v, because '%v'", amount, name, reason))
	}
}

func SubDebt(discord *state.State, service domain.Service) func(
	ctx context.Context,
	data cmdroute.CommandData,
) *api.InteractionResponseData {
	return func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		options := data.Options
		name := options.Find("name").String()
		amount, err := options.Find("amount").IntValue()
		if err != nil || amount < 0 {
			log.Error().Msgf("could not parse amount: %v, err: %s", amount, err)
			return ephemeralMessage("Amount needs to be a non-negative number!")
		}

		err = service.SubDebt(ctx, data.Event.User.ID.String(), data.Event.GuildID.String(), amount)
		if err != nil {
			log.Error().Msgf("could not subtract debt: %s", err)
			return ephemeralMessage("Could not update debt")
		}

		updateDebtsMessage(ctx, discord, service, data.Event.GuildID.String())

		return ephemeralMessage(fmt.Sprintf("Removed %v from %v", amount, name))
	}
}
