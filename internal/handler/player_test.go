package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"slash10k/internal/db"
	mock_db "slash10k/internal/mocks"
	"slash10k/internal/testutil"
	sqlc "slash10k/sql/gen"
	"testing"
)

func TestAddPlayer(t *testing.T) {
	tests := []struct {
		name       string
		players    []string
		withDb     func(c *mock_db.MockQueries)
		wantStatus int
	}{
		{
			name: "adding player 'torfstack'",
			players: []string{
				"torfstack",
			},
			withDb: func(c *mock_db.MockQueries) {
				c.EXPECT().NumberOfPlayers(gomock.Any()).Return(int64(0), nil)
				c.EXPECT().
					AddPlayer(gomock.Any(), "torfstack").
					Return(sqlc.Player{
						ID:   3,
						Name: "torfstack",
					}, nil)
				c.EXPECT().
					SetDebt(gomock.Any(), sqlc.SetDebtParams{
						Amount: 0,
						UserID: db.IdType(3),
					})
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "adding player 'torfstack' and 'neruh'",
			players: []string{
				"torfstack",
				"neruh",
			},
			withDb: func(c *mock_db.MockQueries) {
				c.EXPECT().NumberOfPlayers(gomock.Any()).MinTimes(2).Return(int64(0), nil)
				c.EXPECT().
					AddPlayer(gomock.Any(), "torfstack").
					Return(sqlc.Player{
						ID:   1,
						Name: "torfstack",
					}, nil)
				c.EXPECT().
					AddPlayer(gomock.Any(), "neruh").
					Return(sqlc.Player{
						ID:   2,
						Name: "neruh",
					}, nil)
				c.EXPECT().
					SetDebt(gomock.Any(), sqlc.SetDebtParams{
						Amount: 0,
						UserID: db.IdType(1),
					})
				c.EXPECT().
					SetDebt(gomock.Any(), sqlc.SetDebtParams{
						Amount: 0,
						UserID: db.IdType(2),
					})
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "adding player 'torfstack' after 100 players",
			players: []string{
				"torfstack",
			},
			withDb: func(c *mock_db.MockQueries) {
				c.EXPECT().NumberOfPlayers(gomock.Any()).Return(int64(100), nil)
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			d, q := testutil.QueriesMock(c)
			if tt.withDb != nil {
				tt.withDb(q)
			}

			e := echo.New()
			for _, player := range tt.players {
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rec := httptest.NewRecorder()
				ctx := e.NewContext(req, rec)
				ctx.SetPath("/:name")
				ctx.SetParamNames("name")
				ctx.SetParamValues(player)
				_ = AddPlayer(d)(ctx)
				if ctx.Response().Status != tt.wantStatus {
					t.Fatalf("expected status %d, got %d", tt.wantStatus, ctx.Response().Status)
				}
			}

			c.Finish()
		})
	}
}
