package xlog

import "os"
import "fmt"

// Logger is the interface exposed to the internals of a package, which uses it
// to log messages. This is the other side of a Site.
//
// The 'f' functions work like Printf. The suffix-free functions work like
// 'Print'. The 'e' functions are no-ops if nil is passed for the error;
// otherwise, they print the error as well as the arguments specified.
//
// Fatal calls os.Exit(1) like "log", and Panic panics. Both emit
// Critical-severity log messages before doing so.
type Logger struct {
	Sink
}

//

func (l Logger) Tracef(format string, params ...interface{}) {
	l.ReceiveLocally(SevTrace, format, params...)
}

func (l Logger) Debugf(format string, params ...interface{}) {
	l.ReceiveLocally(SevDebug, format, params...)
}

func (l Logger) Infof(format string, params ...interface{}) {
	l.ReceiveLocally(SevInfo, format, params...)
}

func (l Logger) Noticef(format string, params ...interface{}) {
	l.ReceiveLocally(SevNotice, format, params...)
}

func (l Logger) Warnf(format string, params ...interface{}) {
	l.ReceiveLocally(SevWarn, format, params...)
}

func (l Logger) Errorf(format string, params ...interface{}) {
	l.ReceiveLocally(SevError, format, params...)
}

func (l Logger) Criticalf(format string, params ...interface{}) {
	l.ReceiveLocally(SevCritical, format, params...)
}

func (l Logger) Alertf(format string, params ...interface{}) {
	l.ReceiveLocally(SevAlert, format, params...)
}

func (l Logger) Emergencyf(format string, params ...interface{}) {
	l.ReceiveLocally(SevEmergency, format, params...)
}

func (l Logger) Panicf(format string, params ...interface{}) {
	l.Criticalf("panic: "+format, params...)
	panic(fmt.Sprintf(format, params...))
}

func (l Logger) Fatalf(format string, params ...interface{}) {
	l.Criticalf("fatal: "+format, params...)
	os.Exit(1)
}

func (l Logger) Traceef(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevTrace, err, format, params...)
}

func (l Logger) Debugef(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevDebug, err, format, params...)
}

func (l Logger) Infoef(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevInfo, err, format, params...)
}

func (l Logger) Noticeef(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevNotice, err, format, params...)
}

func (l Logger) Warnef(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevWarn, err, format, params...)
}

func (l Logger) Erroref(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevError, err, format, params...)
}

func (l Logger) Criticalef(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevCritical, err, format, params...)
}

func (l Logger) Alertef(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevAlert, err, format, params...)
}

func (l Logger) Emergencyef(err error, format string, params ...interface{}) {
	l.ReceiveLocallye(SevEmergency, err, format, params...)
}

func (l Logger) Panicef(err error, format string, params ...interface{}) {
	if err != nil {
		l.Criticalef(err, "panic: "+format, params...)
		panic(fmt.Sprintf(format, params...))
	}
}

func (l Logger) Fatalef(err error, format string, params ...interface{}) {
	if err != nil {
		l.Criticalef(err, "fatal: "+format, params...)
		os.Exit(1)
	}
}

// TODO: optimize this

func (l Logger) Trace(params ...interface{}) {
	l.Tracef("%s", fmt.Sprint(params...))
}

func (l Logger) Debug(params ...interface{}) {
	l.Debugf("%s", fmt.Sprint(params...))
}

func (l Logger) Info(params ...interface{}) {
	l.Infof("%s", fmt.Sprint(params...))
}

func (l Logger) Notice(params ...interface{}) {
	l.Noticef("%s", fmt.Sprint(params...))
}

func (l Logger) Warn(params ...interface{}) {
	l.Warnf("%s", fmt.Sprint(params...))
}

func (l Logger) Error(params ...interface{}) {
	l.Errorf("%s", fmt.Sprint(params...))
}

func (l Logger) Critical(params ...interface{}) {
	l.Criticalf("%s", fmt.Sprint(params...))
}

func (l Logger) Alert(params ...interface{}) {
	l.Alertf("%s", fmt.Sprint(params...))
}

func (l Logger) Emergency(params ...interface{}) {
	l.Emergencyf("%s", fmt.Sprint(params...))
}

func (l Logger) Panic(params ...interface{}) {
	l.Panicf("%s", fmt.Sprint(params...))
}

func (l Logger) Fatal(params ...interface{}) {
	l.Fatalf("%s", fmt.Sprint(params...))
}

//
func (l Logger) Tracee(err error, params ...interface{}) {
	if err != nil {
		l.Tracef("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Debuge(err error, params ...interface{}) {
	if err != nil {
		l.Debugf("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Infoe(err error, params ...interface{}) {
	if err != nil {
		l.Infof("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Noticee(err error, params ...interface{}) {
	if err != nil {
		l.Noticef("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Warne(err error, params ...interface{}) {
	if err != nil {
		l.Warnf("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Errore(err error, params ...interface{}) {
	if err != nil {
		l.Errorf("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Criticale(err error, params ...interface{}) {
	if err != nil {
		l.Criticalf("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Alerte(err error, params ...interface{}) {
	if err != nil {
		l.Alertf("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Emergencye(err error, params ...interface{}) {
	if err != nil {
		l.Emergencyf("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Fatale(err error, params ...interface{}) {
	if err != nil {
		l.Fatalf("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) Panice(err error, params ...interface{}) {
	if err != nil {
		l.Panicf("%s: %v", fmt.Sprint(params...), err)
	}
}

func (l Logger) ReceiveLocallye(severity Severity, err error, format string, params ...interface{}) {
	if err == nil {
		return
	}

	l.ReceiveLocally(severity, "%s: %v", fmt.Sprintf(format, params...), err)
}
