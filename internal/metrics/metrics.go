package metrics

import (
	"lob/internal/engine"
	"strconv"
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

func (m *Metrics) RecordTrades() {
	m.TotalTrades.Inc()
}

func (m *Metrics) SetActiveBuyOrders(n uint64) {
	m.ActiveOrdersBuy.Set(float64(n))
}

func (m *Metrics) SetActiveSellOrders(n uint64) {
	m.ActiveOrdersSell.Set(float64(n))
}

func (m *Metrics) RecordLatency(duration time.Duration) {
	m.OrderLatency.Observe(duration.Seconds())
}

func (m *Metrics) SetVolumePL(price uint64, side engine.Side, volume uint64) {

	stringPrice := strconv.FormatUint(price, 10)
	stringSide := strconv.FormatInt(int64(side), 10)

	m.VolumePL.WithLabelValues(stringPrice, stringSide).Set(float64(volume))
}
