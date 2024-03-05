package executil

import (
	"context"
	"os/exec"
	"sync"
)

func RunCtx(
	ctx context.Context,
	cmd *exec.Cmd,
	runClb func() error,
) error {
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
		case <-done:
		}
	}()
	err := runClb()
	close(done)
	wg.Wait()
	return err
}
