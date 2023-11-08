package utils

import "time"

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (p *RealTimeProvider) Now() time.Time {
	return time.Now()
}
