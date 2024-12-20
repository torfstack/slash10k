package db_test

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"slash10k/pkg/db"
	"slash10k/pkg/testutil"
	sqlc "slash10k/sql/gen"
	"slices"
	"testing"
	"time"
)

func Test_Connection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db test")
		return
	}
	type test struct {
		name           string
		withConnection func(*testing.T, db.Connection, context.Context)
	}
	tests := []test{
		{
			name: "can add players and get correct total",
			withConnection: func(t *testing.T, conn db.Connection, ctx context.Context) {
				p1, _ := conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("torfstack"))
				p2, _ := conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("neruh"))
				id1, _ := conn.Queries().GetIdOfPlayer(ctx, testutil.GetIdOfPlayerParams("torfstack"))
				id2, _ := conn.Queries().GetIdOfPlayer(ctx, testutil.GetIdOfPlayerParams("neruh"))
				if p1.ID != id1 || p2.ID != id2 {
					t.Fatalf("Expected IDs to match, got %v!=%v or %v!=%v", p1.ID, id1, p2.ID, id2)
				}
				n, _ := conn.Queries().NumberOfPlayers(ctx)
				if n != 2 {
					t.Fatalf("Expected 2 players, got %d", n)
				}
			},
		},
		{
			name: "adding debts to multiple players and retrieving",
			withConnection: func(t *testing.T, conn db.Connection, ctx context.Context) {
				p1, _ := conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("torfstack"))
				p2, _ := conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("neruh"))
				p3, _ := conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("scurvy"))
				_ = conn.Queries().SetDebt(ctx, sqlc.SetDebtParams{Amount: 30000, UserID: p1.ID})
				_ = conn.Queries().SetDebt(ctx, sqlc.SetDebtParams{Amount: 50000, UserID: p2.ID})
				_ = conn.Queries().SetDebt(ctx, sqlc.SetDebtParams{Amount: 80000, UserID: p3.ID})
				allDebts, _ := conn.Queries().GetAllPlayers(ctx, testutil.GetAllPlayersParams())
				if len(allDebts) != 3 {
					t.Fatalf("Expected 3 debts, got %d", len(allDebts))
				}
				amounts := []int64{allDebts[0].Debt.Amount, allDebts[1].Debt.Amount, allDebts[2].Debt.Amount}
				if slices.Contains(amounts, int64(30000)) == false ||
					slices.Contains(amounts, int64(50000)) == false ||
					slices.Contains(amounts, int64(80000)) == false {
					t.Fatalf("Expected debts of 30000, 50000, 80000, got %v", amounts)
				}
			},
		},
		{
			name: "put bot setup and retrieve",
			withConnection: func(t *testing.T, conn db.Connection, ctx context.Context) {
				_, err := conn.Queries().PutBotSetup(
					ctx, sqlc.PutBotSetupParams{
						GuildID:               testutil.TestGuildIdString(),
						ChannelID:             "channel-id",
						DebtsMessageID:        "debts-message-id",
						RegistrationMessageID: "registration-message-id",
					},
				)
				if err != nil {
					t.Fatalf("Could not put bot setup: %s", err)
				}
				botSetup, err := conn.Queries().GetBotSetup(ctx, testutil.GetBotSetupParams())
				if err != nil {
					t.Fatalf("Could not get bot setup: %s", err)
				}
				if botSetup.ChannelID != "channel-id" ||
					botSetup.GuildID != testutil.TestGuildIdString() ||
					botSetup.DebtsMessageID != "debts-message-id" ||
					botSetup.RegistrationMessageID != "registration-message-id" {
					t.Fatalf(
						"Expected bot setup to be %v, channel-id, debts-message-id, registration-message-id, got %v, %v, %v, %v",
						testutil.TestGuildIdString(),
						botSetup.GuildID,
						botSetup.ChannelID,
						botSetup.DebtsMessageID,
						botSetup.RegistrationMessageID,
					)
				}
			},
		},
		{
			name: "add more than 10 journal entries and retrieve them",
			withConnection: func(t *testing.T, conn db.Connection, ctx context.Context) {
				p, _ := conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("torfstack"))
				for i := range 5 {
					_, _ = conn.Queries().AddJournalEntry(
						ctx, sqlc.AddJournalEntryParams{
							Amount:      int64(i * 10000),
							Description: fmt.Sprintf("added %v", i*10000),
							UserID:      p.ID,
						},
					)
				}
				entries, _ := conn.Queries().GetJournalEntries(ctx, p.ID)
				if len(entries) != 5 {
					t.Fatalf("Expected 5 journal entries, got %d", len(entries))
				}
				for i := range 6 {
					_, _ = conn.Queries().AddJournalEntry(
						ctx, sqlc.AddJournalEntryParams{
							Amount:      int64(i * 10000),
							Description: fmt.Sprintf("added %v", i*10000),
							UserID:      p.ID,
						},
					)
				}
				entries2, _ := conn.Queries().GetJournalEntries(ctx, p.ID)
				if len(entries2) > 10 {
					t.Fatalf("Expected 10 journal entries, got %d", len(entries))
				}
			},
		},
		{
			name: "update journal entry and retrieve it",
			withConnection: func(t *testing.T, conn db.Connection, ctx context.Context) {
				p, _ := conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("torfstack"))
				j, _ := conn.Queries().AddJournalEntry(
					ctx, sqlc.AddJournalEntryParams{
						Amount:      int64(10000),
						Description: fmt.Sprintf("added %v", 10000),
						UserID:      p.ID,
					},
				)
				_, _ = conn.Queries().UpdateJournalEntry(
					ctx, sqlc.UpdateJournalEntryParams{
						Amount:      int64(25000),
						Description: fmt.Sprintf("added %v", 25000),
						ID:          j.ID,
					},
				)
				entries, _ := conn.Queries().GetJournalEntries(ctx, p.ID)
				if len(entries) != 1 || entries[0].Amount != 25000 {
					t.Fatalf("Expected journal entry to be 25000, got %d", entries[0].Amount)
				}
			},
		},
		{
			name: "can not add player with same name discord_id and guild_id twice",
			withConnection: func(t *testing.T, conn db.Connection, ctx context.Context) {
				_, e := conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("torfstack"))
				if e != nil {
					t.Fatalf("Could not add player: %s", e)
				}
				_, e = conn.Queries().AddPlayer(ctx, testutil.AddPlayerParams("torfstack"))
				if e == nil {
					t.Fatalf("Expected error, got nil")
				}
			},
		},
	}
	ctx := context.Background()
	cont, err := setupDatabase(t)
	if err != nil {
		t.Fatalf("Could not set up database: %s", err)
		return
	}
	err = cont.Snapshot(ctx)
	if err != nil {
		t.Fatalf("Could not snapshot database: %s", err)
	}
	var connStr string
	connStr, err = cont.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Could not get connection string: %s", err)
	}
	t.Cleanup(
		func() {
			err = cont.Terminate(ctx)
			if err != nil {
				t.Fatalf("Could not terminate database: %s", err)
			}
		},
	)

	for _, tc := range tests {
		t.Run(
			tc.name, func(t *testing.T) {
				t.Cleanup(
					func() {
						if err = cont.Restore(ctx); err != nil {
							t.Fatalf("Could not restore database: %s", err)
						}
					},
				)

				d := db.NewDatabase(connStr)
				conn, err := d.Connect(ctx)
				if err != nil {
					t.Fatalf("Could not get connection: %s", err)
				}
				defer conn.Close(ctx)

				tc.withConnection(t, conn, ctx)
			},
		)
	}
}

func setupDatabase(t *testing.T) (*postgres.PostgresContainer, error) {
	ctx := context.Background()

	postgresContainer, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16.2-alpine"),
		postgres.WithDatabase("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		t.Fatalf("Could not start postgres container: %s", err)
		return nil, err
	}

	connectionString, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Could not get connection string: %s", err)
		return nil, err
	}

	err = db.Migrate(ctx, connectionString, db.WithMigrationsDir("../../sql/migrations"))
	if err != nil {
		t.Fatalf("Could not run migrations: %s", err)
		return nil, err
	}

	return postgresContainer, nil
}
