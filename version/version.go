package version

import (
	"fmt"
	"runtime"
)

// verion info, values set by linker using ldflag -X
var (
	version string //nolint:gochecknoglobals
	date    string //nolint:gochecknoglobals
	commit  string //nolint:gochecknoglobals
)

func Info() string {
	return fmt.Sprintf("%s@%s [%s].[%s].[%s]",
		version,
		commit,
		date,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
