package cronjob

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type UseCase interface {
	CronSendTransactions(ctx context.Context) error
}

type CronJob struct {
	c  *cron.Cron
	uc UseCase
}

func New(uc UseCase) *CronJob {
	return &CronJob{
		c:  cron.New(),
		uc: uc,
	}
}

func (cj *CronJob) Run(ctx context.Context) error {
	_, err := cj.c.AddFunc("* * * * *", func() {
		log.Info().Msg("run cj.uc.CronSendTransactions")

		if localErr := cj.uc.CronSendTransactions(ctx); localErr != nil {
			log.Error().Msgf("failed cj.uc.CronSendTransactions, err: %v", localErr)
		}
	})
	if err != nil {
		return fmt.Errorf("failed adding cron job cj.uc.CronSendTransactions: %w", err)
	}

	cj.c.Start()
	return nil
}

func (cj *CronJob) Stop() {
	cj.c.Stop()
}

// https://crontab.guru/
