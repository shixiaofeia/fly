package logging

type PrefixLog struct {
	Prefix string
}

// NewPrefixLog  带前缀的log
func NewPrefixLog(prefix string) *PrefixLog {
	return &PrefixLog{Prefix: prefix + " "}
}

// Info
func (p *PrefixLog) Info(msg string) {
	Log.Info(p.Prefix + msg)
}

// Debug
func (p *PrefixLog) Debug(msg string) {
	Log.Debug(p.Prefix + msg)
}

// Warn
func (p *PrefixLog) Warn(msg string) {
	Log.Warn(p.Prefix + msg)
}

// Error
func (p *PrefixLog) Error(msg string) {
	Log.Error(p.Prefix + msg)
}
