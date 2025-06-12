package worker

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/utils/slices"
)

type JobType int

const (
	JobAtTime JobType = iota
	JobEveryDuration
)

type JobConfig struct {
	Type     JobType
	AtTimes  []time.Time   // e.g., "05:00"
	Interval time.Duration // e.g., 5 * time.Second
	Job      func()
	About    string
}

type Scheduler struct {
	raw gocron.Scheduler
	log *zerolog.Event
}

func NewScheduler() *Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create scheduler")
	}

	return &Scheduler{
		raw: scheduler,
		log: log.Info(),
	}
}

func (s *Scheduler) Start() {
	s.raw.Start()
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

	if cfg.Type == JobAtTime {
		_, err := s.raw.NewJob(
			gocron.DailyJob(1, func() []gocron.AtTime {
				return slices.Map(cfg.AtTimes, func(t time.Time) gocron.AtTime {
					return gocron.NewAtTime(uint(t.Hour()), uint(t.Minute()), uint(t.Second()))
				})
			}),
			gocron.NewTask(wrappedJob),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to add job")
		}
	} else if cfg.Type == JobEveryDuration {
		_, err := s.raw.NewJob(
			gocron.DurationJob(cfg.Interval),
			gocron.NewTask(wrappedJob),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to add job")
		}
	}
}

func (s *Scheduler) Shutdown() {
	s.raw.Shutdown()
}
