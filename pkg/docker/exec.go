package docker

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"

	executil "github.com/cardinalby/xgo-pack/pkg/util/exec"
)

const dockerBin = "docker"

func Exec(ctx context.Context, args ...string) error {
	execCmd := exec.Command(dockerBin, args...)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		_ = execCmd.Process.Kill()
	}()

	return executil.RunCtx(ctx, execCmd, func() error {
		out, err := execCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("docker %v: %w. %v", args, err, string(out))
		}
		return nil
	})
}

func ExecRes(ctx context.Context, args ...string) (stdout string, err error) {
	execCmd := exec.Command(dockerBin, args...)
	var stdoutBuff, stderrBuff bytes.Buffer
	execCmd.Stdout = &stdoutBuff
	execCmd.Stderr = &stderrBuff

	err = executil.RunCtx(ctx, execCmd, func() error {
		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("docker %v: %w. %v", args, err, stderrBuff.String())
		}
		return nil
	})

	return stdoutBuff.String(), err
}
