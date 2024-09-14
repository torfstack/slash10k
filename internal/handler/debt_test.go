package handler

import (
	"bytes"
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

type AddDebtTestParam struct {
	Player      string
	Amount      string
	Description string
}

func TestAddDebt(t *testing.T) {
	tests := []struct {
		name        string
		params      []AddDebtTestParam
		withQueries func(q *mock_db.MockQueries)
		wantStatus  int
	}{
		{
			name: "adding debt to player 'torfstack'",
			params: []AddDebtTestParam{
				{
					Player: "torfstack",
					Amount: "10000",
				},
			},
			withQueries: func(q *mock_db.MockQueries) {
				q.EXPECT().
					GetIdOfPlayer(gomock.Any(), "torfstack").
					Return(int32(1), nil)
				q.EXPECT().
					GetDebt(gomock.Any(), gomock.AnyOf(db.IdType(1))).
					Return(sqlc.Debt{
						Amount: 20000,
						UserID: db.IdType(1),
					}, nil)
				q.EXPECT().
					SetDebt(gomock.Any(), sqlc.SetDebtParams{
						Amount: 30000,
						UserID: db.IdType(1),
					})
				q.EXPECT().
					GetAllDebts(gomock.Any())
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "adding debt to player 'torfstack' with description",
			params: []AddDebtTestParam{
				{
					Player:      "torfstack",
					Amount:      "10000",
					Description: "Trash-AFK",
				},
			},
			withQueries: func(q *mock_db.MockQueries) {
				q.EXPECT().
					GetIdOfPlayer(gomock.Any(), "torfstack").
					Return(int32(1), nil)
				q.EXPECT().
					GetDebt(gomock.Any(), gomock.Any()).
					Return(sqlc.Debt{
						Amount: 20000,
						UserID: db.IdType(1),
					}, nil)
				q.EXPECT().
					SetDebt(gomock.Any(), sqlc.SetDebtParams{
						Amount: 30000,
						UserID: db.IdType(1),
					})
				q.EXPECT().
					GetAllDebts(gomock.Any())
				q.EXPECT().
					AddJournalEntry(gomock.Any(), sqlc.AddJournalEntryParams{
						Amount:      10000,
						Description: "Trash-AFK",
						UserID:      db.IdType(1),
					})
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "subtracting debt from player 'neruh' lowers debt and removes journal entry",
			params: []AddDebtTestParam{
				{
					Player: "neruh",
					Amount: "-30000",
				},
			},
			withQueries: func(q *mock_db.MockQueries) {
				q.EXPECT().
					GetIdOfPlayer(gomock.Any(), "neruh").
					Return(int32(2), nil)
				q.EXPECT().
					GetDebt(gomock.Any(), gomock.Any()).
					Return(sqlc.Debt{
						Amount: 70000,
						UserID: db.IdType(2),
					}, nil)
				q.EXPECT().
					SetDebt(gomock.Any(), sqlc.SetDebtParams{
						Amount: 40000,
						UserID: db.IdType(2),
					})
				q.EXPECT().
					GetAllDebts(gomock.Any())
				q.EXPECT().
					GetJournalEntries(gomock.Any(), db.IdType(2)).
					Return([]sqlc.DebtJournal{
						{
							ID:          1,
							Amount:      10000,
							Description: "Trash-AFK",
						},
						{
							ID:          2,
							Amount:      60000,
							Description: "Boss reset fail :(",
						},
					}, nil)
				q.EXPECT().
					DeleteJournalEntry(gomock.Any(), int32(1))
				q.EXPECT().
					UpdateJournalEntry(gomock.Any(), sqlc.UpdateJournalEntryParams{
						Amount:      40000,
						Description: "Boss reset fail :(",
						ID:          2,
					})
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			d, q := testutil.QueriesMock(c)
			if tt.withQueries != nil {
				tt.withQueries(q)
			}

			e := echo.New()
			for _, param := range tt.params {
				var jsonData []byte
				if param.Description != "" {
					jsonData = []byte(`{
						"description": "` + param.Description + `"
					}`)
				}
				req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonData))
				req.Header.Add("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				ctx := e.NewContext(req, rec)
				ctx.SetPath("/:player/debt/:amount")
				ctx.SetParamNames("player", "amount")
				ctx.SetParamValues(param.Player, param.Amount)
				_ = AddDebt(d)(ctx)
				if ctx.Response().Status != tt.wantStatus {
					t.Fatalf("expected status %d, got %d", tt.wantStatus, ctx.Response().Status)
				}
			}

			c.Finish()
		})
	}
}
