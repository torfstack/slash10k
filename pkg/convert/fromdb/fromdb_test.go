package fromdb

import (
	"github.com/jackc/pgx/v5/pgtype"
	"reflect"
	"slash10k/pkg/models"
	sqlc "slash10k/sql/gen"
	"testing"
	"time"
)

func TestFromAllPlayers(t *testing.T) {
	type args struct {
		allPlayers []sqlc.GetAllPlayersRow
	}
	tests := []struct {
		name string
		args args
		want []models.Player
	}{
		{
			name: "empty slice",
			args: args{allPlayers: []sqlc.GetAllPlayersRow{}},
			want: []models.Player{},
		},
		{
			name: "one player",
			args: args{
				allPlayers: []sqlc.GetAllPlayersRow{
					{
						Player: sqlc.Player{
							ID:          1,
							DiscordID:   "123",
							DiscordName: "torfstack",
							GuildID:     "456",
							Name:        "Torfstack",
						},
						Debt: sqlc.Debt{
							ID:          1,
							Amount:      1000,
							LastUpdated: pgtype.Timestamp{Time: time.Unix(0, 0)},
							UserID:      1,
						},
					},
				},
			},
			want: []models.Player{
				{
					Id:          1,
					DiscordId:   "123",
					DiscordName: "torfstack",
					GuildId:     "456",
					Name:        "Torfstack",
					Debt: models.Debt{
						Id:          1,
						Amount:      1000,
						LastUpdated: 0,
						UserId:      1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := FromAllPlayers(tt.args.allPlayers); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FromAllPlayers() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
