package engine

import "time"

type NullStore struct{} // struct for a null store to allow engine to instantiate before postgres is working

type NullMetricRecorder struct{} //struct for a null metric recorder

// Funcs to satisfy MetricRecorder interface for NummMetricRecorder
func (nmr *NullMetricRecorder) RecordTrades()                                      {}
func (nmr *NullMetricRecorder) SetActiveBuyOrders(n uint64)                        {}
func (nmr *NullMetricRecorder) SetActiveSellOrders(n uint64)                       {}
func (nmr *NullMetricRecorder) RecordLatency(duration time.Duration)               {}
func (nmr *NullMetricRecorder) SetVolumePL(price uint64, side Side, volume uint64) {}
func (nmr *NullMetricRecorder) DeleteVolumePL(price uint64, side Side)             {}

func (n *NullStore) WriteTrade(trade Trade) error {
	return nil
}
