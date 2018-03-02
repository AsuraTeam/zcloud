package xlog

import "strings"

// Log message severity. This is the syslog severity order.
//
// Note that Emergency and Alert are system-level severities.  Generally
// speaking, application programs should not emit log messages at such
// severities unless they are programs which monitor the system for
// system-level issues. i.e., programs should never emit Emergency or Alert
// messages regarding issues with their own operation.
//
// Programs suffering from critical failures should emit log messages at the
// Critical severity. The Panic*() and Fatal*() log message functions in this
// package emit log messages at the Critical level.
//
// The Error severity should be used when errors occur which do not constitute
// a critical or unrecoverable failure of the program.
//
// Any severity less severe than Debug is not part of the syslog severity
// order. These are converted to messages of Debug severity when exported
// to e.g. syslog.
//
// Trace should be used for extremely verbose debugging information which
// is likely to be used mainly for debugging and is of such verbosity that it
// may overwhelm a programmer unless enabled only for a few specific
// facilities.
type Severity int

const (
	SevEmergency Severity = iota
	SevAlert
	SevCritical
	SevError
	SevWarn
	SevNotice
	SevInfo
	SevDebug
	SevTrace
	SevNone Severity = -1 // (Do not use.)
)

var severityString = map[Severity]string{
	SevEmergency: "EMERGENCY", // EM EMR EMER
	SevAlert:     "ALERT",     // AL ALR ALER
	SevCritical:  "CRITICAL",  // CR CRT CRIT
	SevError:     "ERROR",     // ER ERR ERRO
	SevWarn:      "WARN",      // WA WRN WARN
	SevNotice:    "NOTICE",    // NO NOT NOTC
	SevInfo:      "INFO",      // IN INF INFO
	SevDebug:     "DEBUG",     // DE DBG DEBG
	SevTrace:     "TRACE",     // TR TRC TRAC
}

var ansiSeverityString = map[Severity]string{
	SevEmergency: "\x1B[41;37mEMERGENCY\x1B[0m",
	SevAlert:     "\x1B[41;37mALERT\x1B[0m",
	SevCritical:  "\x1B[41;37mCRITICAL\x1B[0m",
	SevError:     "\x1B[31mERROR\x1B[0m",
	SevWarn:      "\x1B[33mWARN\x1B[0m",
	SevNotice:    "NOTICE\x1B[0m",
	SevInfo:      "INFO\x1B[0m",
	SevDebug:     "DEBUG\x1B[0m",
	SevTrace:     "TRACE\x1B[0m",
}

var severityValue = map[string]Severity{}

func init() {
	for k, v := range severityString {
		severityValue[v] = k
	}
}

// Returns the severity as an uppercase unabbreviated string.
func (severity Severity) String() string {
	return severityString[severity]
}

// Parse a severity string.
func ParseSeverity(severity string) (s Severity, ok bool) {
	s, ok = severityValue[strings.ToUpper(severity)]
	return
}

// Returns the syslog-compatible severity. Converts severities
// less severe than Debug to Debug.
func (severity Severity) Syslog() Severity {
	if severity > SevDebug {
		return SevDebug
	}
	return severity
}
