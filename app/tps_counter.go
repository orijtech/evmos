package app

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/tendermint/tendermint/libs/log"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var mTransactions = stats.Int64("transactions", "the number of transactions after .EndBlocker", "1")
var viewTransactions = &view.View{
	Name:        "evmos/app",
	Measure:     mTransactions,
	Description: "The transactions processed",
	TagKeys:     nil,
	Aggregation: view.Count(),
}

func ObservabilityViews() (views []*view.View) {
	views = append(views, viewTransactions)
	return views
}

type tpsCounter struct {
	nTxn         uint64
	logger       log.Logger
	reportPeriod time.Duration
}

func newTPSCounter(logger log.Logger) *tpsCounter {
	return &tpsCounter{logger: logger}
}

func (tpc *tpsCounter) increment() { atomic.AddUint64(&tpc.nTxn, 1) }

const defaultTPSReportPeriod = 10 * time.Second

func (tpc *tpsCounter) start(ctx context.Context) error {
	tpsReportPeriod := defaultTPSReportPeriod
	if tpc.reportPeriod > 0 {
		tpsReportPeriod = tpc.reportPeriod
	}
	ticker := time.NewTicker(tpsReportPeriod)
	defer ticker.Stop()

	var lastNTxn uint64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			// Report the number of transactions seen in the designated period of time.
			latestNTxn := atomic.LoadUint64(&tpc.nTxn)
			if latestNTxn == 0 {
				// No need to report it for the first time.
				continue
			}

			if latestNTxn >= lastNTxn {
				// Record these stats with OpenCensus so that they can be exported say to Prometheus, or any metrics backends.
				nTxn := int64(latestNTxn - lastNTxn)
				if nTxn >= 0 {
					stats.Record(ctx, mTransactions.M(nTxn))
				}
				// Record to our logger for easy examination in the logs.
				secs := float64(tpsReportPeriod) / float64(time.Second)
				tpc.logger.Info("Transactions per second", "tps", float64(latestNTxn-lastNTxn)/secs)
			}

			lastNTxn = latestNTxn
		}
	}

	return nil
}
