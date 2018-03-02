// Package xlog provides a hierarchical, configurable logging system suitable
// for use in libraries.
package xlog

import "fmt"
import "sync"

// Site is the interface exposed to the externals of a package, which uses it
// to configure the logger. This is the other side of a Logger.
type Site interface {
	// The facility name.
	Name() string

	// Sets the severity condition.
	SetSeverity(severity Severity)

	// Set the sink which will receive all messages from this logger.
	SetSink(sink Sink)
}

var loggersMutex sync.RWMutex
var loggers = map[string]*logger{}

// Creates a logger which is subordinate to another logger.
//
// All log messages will be forwarded through the parent logger, meaning that
// any filtration or forwarding settings set on the parent logger will
// also apply to this one.
//
// The name of the logger is formed by appending the name given to the name
// of the parent logger and ".". If site is nil, behaves like New().
func NewUnder(name string, site Site) (Logger, Site) {
	if site == nil {
		return New(name)
	}

	sink, ok := site.(Sink)
	if !ok {
		panic("site does not implement sink")
	}

	l, s := New(site.Name() + "." + name)
	s.SetSink(sink)
	return l, s
}

// Creates a new logger.
//
// While there are no particular restrictions on facility names, the preferred
// convention for the facility name is a dot-separated hierarchy of words
// matching [a-zA-Z0-9_-]+. Hyphens are preferred over underscores, and
// uppercase should be avoided in most cases.
//
// The facility name should reflect the package and, if the package's
// status as a subpackage is of particular significance or grouping
// is desirable, a parent package.
//
// For example, if you have a package foo which has 10 subpackages which
// implement different parts of foo, you might give them facility names like
// "foo.alpha", "foo.beta", "foo.gamma", etc.
//
// Typical usage:
//
//     var log, Log = xlog.New("facility name")
//
func New(name string) (Logger, Site) {
	loggersMutex.Lock()
	defer loggersMutex.Unlock()

	if _, ok := loggers[name]; ok {
		panic(fmt.Sprintf("Logger name conflict: logger with name %s already exists", name))
	}

	log := &logger{
		parent:      rootLogger,
		maxSeverity: SevTrace,
		name:        name,
	}

	loggers[name] = log

	return Logger{log}, log
}

// Like New, but the created logger by default doesn't output anything but the
// most severe errors. Intended for use by libraries so that consuming code
// needs to opt in to log output by that library.
func NewQuiet(name string) (Logger, Site) {
	l, s := New(name)
	s.SetSeverity(SevCritical)
	return l, s
}

type logger struct {
	maxSeverity Severity
	name        string

	parent Sink
}

// Sink is implemented by objects that can receive log messages from loggers
// deeper in the hierarchy.
type Sink interface {
	ReceiveLocally(sev Severity, format string, params ...interface{})
	ReceiveFromChild(sev Severity, format string, params ...interface{})
}

func init() {
	RootSink.Add(StderrSink)
}

var rootLogger = &logger{
	parent:      &RootSink,
	maxSeverity: SevTrace,
}

// The root logger.
var Root Site = rootLogger

// The sink which is used by default by the root logger.
var RootSink MultiSink

func (l *logger) Name() string {
	return l.name
}

func (l *logger) SetSeverity(sev Severity) {
	l.maxSeverity = sev
}

func (l *logger) SetSink(sink Sink) {
	l.parent = sink
}

func (l *logger) ReceiveLocally(sev Severity, format string, params ...interface{}) {
	format = l.localPrefix() + format // XXX unsafe format string
	l.remoteLogf(sev, format, params...)
}

func (l *logger) remoteLogf(sev Severity, format string, params ...interface{}) {
	if sev > l.maxSeverity {
		return
	}

	if l.parent != nil {
		l.parent.ReceiveFromChild(sev, format, params...)
	}
}

func (l *logger) ReceiveFromChild(sev Severity, format string, params ...interface{}) {
	l.remoteLogf(sev, format, params...)
}

func (l *logger) localPrefix() string {
	if l.name != "" {
		return l.name + ": "
	}
	return ""
}

// Calls a function for every Site which has been created.
//
// Do not attempt to create new loggers from the callback.
func VisitSites(siteFunc func(s Site) error) error {
	loggersMutex.RLock()
	defer loggersMutex.RUnlock()

	for _, v := range loggers {
		err := siteFunc(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// LogClosure can be used to pass a function that returns a string
// to a log method call. This is useful if the computation of a log message
// is expensive and the message will often be filtered.
type LogClosure func() string

func (c LogClosure) String() string {
	return c()
}
