package mocks

import "time"

type FakeTimeProvider struct {
	NowTime time.Time
}

func (p *FakeTimeProvider) Now() time.Time {
	return p.NowTime
}
