package buildctx

import (
	"context"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/cardinalby/xgo-pack/pkg/pipeline/config/cfgtypes"
	"github.com/cardinalby/xgo-pack/pkg/util/logging"
	"github.com/stretchr/testify/require"
)

type testArtifact struct {
	x int
}

func (t *testArtifact) GetPath() string {
	return strconv.Itoa(t.x)
}

func TestArtifacts_Get(t *testing.T) {
	for j := 0; j < 100; j++ {
		artifacts := NewArtifacts(logging.NewNopLogger())
		ctx := NewContext(context.Background(), cfgtypes.Config{}, artifacts, logging.NewNopLogger())
		var artifactsCount = 10
		for i := 0; i < artifactsCount; i++ {
			i := i
			artifacts.RegisterBuilder(Kind(strconv.Itoa(i)), func(ctx Context) (Artifact, error) {
				if i > 0 {
					depStr := strconv.Itoa(rand.Intn(i))
					depA, err := artifacts.Get(ctx, Kind(depStr))
					require.NoError(t, err)
					require.Equal(t, depA.GetPath(), depStr)
				}
				return &testArtifact{x: i}, nil
			})
		}
		wg := sync.WaitGroup{}
		for i := 0; i < 200; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				numStr := strconv.Itoa(i % artifactsCount)
				a, err := artifacts.Get(ctx, Kind(numStr))
				require.NoError(t, err)
				require.Equal(t, a.GetPath(), numStr)
			}()
		}
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()
		select {
		case <-time.After(5 * time.Second):
			_ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
			t.Fatal("deadlock")
		case <-done:
		}
	}
}
