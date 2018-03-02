package xlog

import "fmt"

// Interface compatible with "log/syslog".Writer.
type Syslogger interface {
	Alert(m string) error
	Crit(m string) error
	Debug(m string) error
	Emerg(m string) error
	Err(m string) error
	Info(m string) error
	Notice(m string) error
	Warning(m string) error
}

// A sink that logs to a "log/syslog".Writer-like interface.
type SyslogSink struct {
	s           Syslogger
	minSeverity Severity
}

// Create a new syslog sink. "log/syslog".Writer implements Syslogger.
func NewSyslogSink(syslogger Syslogger) *SyslogSink {
	return &SyslogSink{
		s:           syslogger,
		minSeverity: SevDebug,
	}
}

func (ss *SyslogSink) SetSeverity(sev Severity) {
	ss.minSeverity = sev
}

func (ss *SyslogSink) ReceiveLocally(sev Severity, format string, params ...interface{}) {
	ss.ReceiveFromChild(sev, format, params...)
}

func (ss *SyslogSink) ReceiveFromChild(sev Severity, format string, params ...interface{}) {
	if sev > ss.minSeverity {
		return
	}

	s := fmt.Sprintf(format, params...)
	switch sev {
	case SevEmergency:
		ss.s.Emerg(s)
	case SevAlert:
		ss.s.Alert(s)
	case SevCritical:
		ss.s.Crit(s)
	case SevError:
		ss.s.Err(s)
	case SevWarn:
		ss.s.Warning(s)
	case SevNotice:
		ss.s.Notice(s)
	case SevInfo:
		ss.s.Info(s)
	default:
		ss.s.Debug(s)
	}
}
