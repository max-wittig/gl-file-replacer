package version

import (
	"fmt"
	"runtime"
)

// GitCommit that was compiled. This will be filled in by the compiler.
var GitCommit string

// BuildDate is the date when this version was build
var BuildDate string

// Version is the main version number that is being run at the moment.
const Version = "0.0.1"

// GoVersion is the current go runtime version
var GoVersion = runtime.Version()

// OsArch is the os architecture this is running on
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
