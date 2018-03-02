# xlog [![GoDoc](https://godoc.org/github.com/hlandau/xlog?status.svg)](https://godoc.org/github.com/hlandau/xlog)

Yet another logging package for Go.

The main thing about this package is that it's good for usage in libraries and
doesn't involve itself with policy. Essentially, if, as a library, you want to
log something, you write this:

```go
var log, Log = xlog.NewQuiet("my.logger.name")

func Foo() {
  log.Debugf("Bar")
}
```

The `log` variable is what you use to log, and the `Log` variable is exported
from the package and provides methods for controlling the log site. (These are
actually two interfaces to the same object which enforce this pattern.)

The idea is that consuming code can call somepkg.Log to control where it logs
to, at what verbosity, etc.

You should instantiate with `NewQuiet` if you are a library, as this suppresses
most log messages by default. Loggers created with `New` don't suppress
any log messages by default.

xlog uses a traditional Print/Printf interface. It also has the following
conveniences:

  - Methods ending in `e`, such as `Debuge`, take an error as their first
    argument and are no-ops if it is nil.

  - `Fatal` and `Panic` call os.Exit(1) and Panic, respectively.
    The `e` variants of these are no-ops if the error is nil, providing
    a simple assertion mechanism.

xlog uses syslog severities (Emergency, Alert, Critical, Error, Warning,
Notice, Info, Debug) and also provides a Trace severity which is even less
severe than Debug. You should generally not emit Alert or Emergency severities
from your code, as these are of system-level significance.

You can visit all registered log sites in order to configure loggers
programmatically.

Loggers should be named via a dot-separated hierarchy with names in lowercase.
If you have a repository called `foo` and a subpackage `baz`, naming the logger
for that subpackage `foo.baz` might be reasonable.

Loggers are arranged in a hierarchy terminating in the root logger, which is
configured to log to stderr by default. You can create a logger under another
logger using `NewUnder`.

## Licence

    © 2014—2016 Hugo Landau <hlandau@devever.net>  MIT License

[Licenced under the licence with SHA256 hash
`fd80a26fbb3f644af1fa994134446702932968519797227e07a1368dea80f0bc`, a copy of
which can be found
here.](https://raw.githubusercontent.com/hlandau/acme/master/_doc/COPYING.MIT)
