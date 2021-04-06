package scheduler

import (
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/scheduler"
)

type Provider struct{}

func (s Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		scheduler.Provider(s.jobCreator),
	}
}

func (s Provider) Register(cc infra.Binder) {}
func (s Provider) Boot(cc infra.Resolver)   {}

func (s Provider) jobCreator(cc infra.Resolver, creator scheduler.JobCreator) {

}
