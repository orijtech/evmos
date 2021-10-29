package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func TestTPSCounter(t *testing.T) {
	buf := new(bytes.Buffer)
	tpc := newTPSCounter(&writerLogger{w: buf})
	tpc.reportPeriod = 5 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	go tpc.start(ctx)

	// Concurrently increment the counter.
	n := 50
	repeat := 5
	go func() {
		defer cancel()
		for i := 0; i < repeat; i++ {
			for j := 0; j < n; j++ {
				tpc.increment()
			}
			<-time.After(tpc.reportPeriod)
		}
	}()

	<-ctx.Done()

	// We expect that the TPS reported will be:
	// 100 / 5ms => 100 / 0.005s = 20,000 TPS
	wantTPS := float64(n) / (float64(tpc.reportPeriod) / float64(time.Second))
	lines := strings.Split(buf.String(), "\n")
	require.Equal(t, repeat+1, len(lines), "Expected exactly n repeats")
	want := strings.Repeat(fmt.Sprintf("Transactions per second %.0f\n", wantTPS), repeat)
	require.Equal(t, want, buf.String(), "Expecting the exact matches")
}

type writerLogger struct {
	w io.Writer
	log.Logger
}

var _ log.Logger = (*writerLogger)(nil)

func (wl *writerLogger) Info(msg string, keyVals ...interface{}) {
	fmt.Fprintf(wl.w, msg+" "+fmt.Sprintf("%v\n", keyVals...))
}
