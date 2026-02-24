package metrics

import (
	"lob/internal/engine"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//metrics to implement:
//Total trades
//Trades per second
//Volume at price level
//Active orders on each book
//Order processing latency

type Metrics struct {
	TotalTrades      prometheus.Counter
	VolumePL         prometheus.GaugeVec
	ActiveOrdersBuy  prometheus.Gauge
	ActiveOrdersSell prometheus.Gauge
	OrderLatency     prometheus.Histogram
}

func NewMetrics(reg prometheus.Registerer) *Metrics {

	m := &Metrics{
		TotalTrades: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "lob_total_trades",
			Help: "The total number of trades processed",
		}),
		VolumePL: *promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "lob_volume_pl",
			Help: "The current volume of trades at each Price Level",
		}, []string{"price", "side"}),
		ActiveOrdersBuy: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Name: "lob_active_buy_orders",
			Help: "The total number of active orders on the buy book",
		}),
		ActiveOrdersSell: promauto.With(reg).NewGauge(prometheus.GaugeOpts{
			Name: "lob_active_sell_orders",
			Help: "The total number of active orders on the sell book",
		}),
		OrderLatency: promauto.With(reg).NewHistogram(prometheus.HistogramOpts{
			Name: "lob_order_latency",
			Help: "The latency of order processing",
		}),
	}

	return m

}

func (m *Metrics) RecordTrades()                                             {}
func (m *Metrics) SetActiveBuyOrders(n uint64)                               {}
func (m *Metrics) SetActiveSellOrders(n uint64)                              {}
func (m *Metrics) RecordLatency(time.Time)                                   {}
func (m *Metrics) SetVolumePL(price uint64, side engine.Side, volume uint64) {}
