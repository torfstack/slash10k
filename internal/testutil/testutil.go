package testutil

import (
	"go.uber.org/mock/gomock"
	"slash10k/internal/db"
	mockdb "slash10k/internal/mocks"
	"testing"
)

func QueriesMock(c *gomock.Controller) (db.Database, *mockdb.MockQueries) {
	mockDb := mockdb.NewMockDatabase(c)
	mockConn := mockdb.NewMockConnection(c)
	mockTx := mockdb.NewMockTransaction(c)
	mockQueries := mockdb.NewMockQueries(c)
	mockDb.EXPECT().
		Connect(gomock.Any(), gomock.Any()).
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
