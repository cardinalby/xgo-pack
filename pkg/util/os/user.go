package osutil

import (
	"fmt"
	"os/user"
	"runtime"
	"strconv"

	"github.com/cardinalby/xgo-pack/pkg/consts"
)

func GetLinuxUser() (uid int, gid int, ok bool, err error) {
	if //goland:noinspection GoBoolExpressions
	runtime.GOOS != string(consts.OsLinux) {
		return 0, 0, false, nil
	}
	usr, err := user.Current()
	if err != nil {
		return 0, 0, false, fmt.Errorf("error getting current user: %w", err)
	}
	uid, err = strconv.Atoi(usr.Uid)
	if err != nil {
		return 0, 0, false, fmt.Errorf("error converting current user uid '%s' to int: %w", usr.Uid, err)
	}
	gid, err = strconv.Atoi(usr.Gid)
	if err != nil {
		return 0, 0, false, fmt.Errorf("error converting current user gid '%s' to int: %w", usr.Gid, err)
	}
	return uid, gid, true, nil
}
