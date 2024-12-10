package testutil

import (
	"go.uber.org/mock/gomock"
	"slash10k/pkg/config"
	"slash10k/pkg/db"
	mockdb "slash10k/pkg/mocks"
	sqlc "slash10k/sql/gen"
	"testing"
)

func QueriesMock(c *gomock.Controller) (db.Database, *mockdb.MockQueries) {
	mockDb := mockdb.NewMockDatabase(c)
	mockConn := mockdb.NewMockConnection(c)
	mockTx := mockdb.NewMockTransaction(c)
	mockQueries := mockdb.NewMockQueries(c)
	mockDb.EXPECT().
		Connect(gomock.Any()).
		MinTimes(1).
		Return(mockConn, nil)
	mockConn.EXPECT().
		Queries().
		AnyTimes().
		Return(mockQueries)
	mockConn.EXPECT().
		StartTransaction(gomock.Any()).
		AnyTimes().
		Return(mockTx, nil)
	mockTx.EXPECT().
		Queries().
		AnyTimes().
		Return(mockQueries)
	mockConn.EXPECT().
		Close(gomock.Any()).
		MinTimes(1).
		Return(nil)
	mockTx.EXPECT().
		Commit(gomock.Any()).
		AnyTimes().
		Return(nil)
	return mockDb, mockQueries
}

func WithoutError(t *testing.T, a interface{}, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestGuildIdString() string {
	return string(config.TorfstackServerGuildIdString)
}

func AddPlayerParams(discordId string) sqlc.AddPlayerParams {
	return sqlc.AddPlayerParams{
		DiscordID: discordId,
		GuildID:   TestGuildIdString(),
	}
}

func SetDebtParams(id int32, amount int64) sqlc.SetDebtParams {
	return sqlc.SetDebtParams{
		Amount: amount,
		UserID: id,
	}
}

func GetIdOfPlayerParams(discordId string) sqlc.GetIdOfPlayerParams {
	return sqlc.GetIdOfPlayerParams{
		DiscordID: discordId,
		GuildID:   TestGuildIdString(),
	}
}

func GetAllPlayersParams() string {
	return TestGuildIdString()
}

func GetBotSetupParams() string {
	return TestGuildIdString()
}
