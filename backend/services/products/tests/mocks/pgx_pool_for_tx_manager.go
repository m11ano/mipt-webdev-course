package mocks

import (
	mock "github.com/stretchr/testify/mock"
)

func NewPgxPoolMockForTxManager() *PgxPool {
	mockPool := new(PgxPool)
	mockTx := new(PoolTxInterface)

	mockPool.On("BeginTx", mock.Anything, mock.AnythingOfType("pgx.TxOptions")).Return(mockTx, nil)

	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	return mockPool
}
