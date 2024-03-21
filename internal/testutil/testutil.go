package testutil

import (
	"go.uber.org/mock/gomock"
	"slash10k/internal/db"
	mockdb "slash10k/internal/mocks"
)

func QueriesMock(c *gomock.Controller) (db.Database, *mockdb.MockQueries) {
	mockDb := mockdb.NewMockDatabase(c)
	mockConn := mockdb.NewMockConnection(c)
	mockQueries := mockdb.NewMockQueries(c)
	mockDb.EXPECT().
		Connect(gomock.Any()).
		MinTimes(1).
		Return(mockConn, nil)
	mockConn.EXPECT().
		Queries().
		MinTimes(1).
		Return(mockQueries)
	mockConn.EXPECT().
		Close(gomock.Any()).
		MinTimes(1).
		Return(nil)
	return mockDb, mockQueries
}
