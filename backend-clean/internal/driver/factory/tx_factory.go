package factory

import "immortal-architecture-clean/backend/internal/port"

func NewTxFactory(tx port.TxManager) func() port.TxManager {
	return func() port.TxManager {
		return tx
	}
}
