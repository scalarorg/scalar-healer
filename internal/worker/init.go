package worker

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type JobType int

const (
	JobAtTime JobType = iota
	JobEveryDuration
)

type JobConfig struct {
	Type     JobType
	AtTime   time.Time     // e.g., "05:00"
	Interval time.Duration // e.g., 5 * time.Second
	Job      func()
	About    string
}

type Scheduler struct {
	raw  *gocron.Scheduler
	jobs []JobConfig
	log  *zerolog.Event
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		raw: gocron.NewScheduler(time.UTC),
		log: log.Info(),
	}
}

func (s *Scheduler) Start() {
	for _, job := range s.jobs {
		switch job.Type {
		case JobAtTime:
			s.raw.Every(1).Day().At(job.AtTime).Do(job.Job)
		case JobEveryDuration:
			s.raw.Every(job.Interval).Do(job.Job)
		}
	}
	s.raw.StartAsync()
}

func (s *Scheduler) AddJob(cfg JobConfig) {
	// Wrap the original job with panic recovery
	wrappedJob := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Interface("panic", r).
					Msg("Recovered from panic in scheduled job")
			}
		}()

		s.log.Msgf("Running job %s", cfg.About)
		cfg.Job()
	}
	cfg.Job = wrappedJob
	s.jobs = append(s.jobs, cfg)
}

func (s *Scheduler) Len() int {
	return len(s.jobs)
}

func (s *Scheduler) Shutdown() {
	s.raw.Stop()
}
