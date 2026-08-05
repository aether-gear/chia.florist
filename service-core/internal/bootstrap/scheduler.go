package bootstrap

import (
	"context"
	"time"

	applogger "service-core/internal/common/logger"
	paymentJob "service-core/internal/modules/payment/infra/job"
)

// Job defines the contract for any background runnable scheduled job.
type Job interface {
	Start(ctx context.Context)
}

// Scheduler handles registration and lifecycle management of background jobs.
type Scheduler struct {
	jobs []Job
}

func NewScheduler(cfg Config, container *Container, logger applogger.Logger) *Scheduler {
	s := &Scheduler{
		jobs: make([]Job, 0),
	}

	if container == nil {
		return s
	}

	syncInterval := time.Duration(cfg.PaymentSync.IntervalMinutes) * time.Minute
	paymentSyncJob := paymentJob.NewPaymentSyncJob(
		&container.SyncPendingPayments,
		syncInterval,
		logger,
	)
	s.Register(paymentSyncJob)

	expiryInterval := time.Duration(cfg.PaymentExpiry.IntervalMinutes) * time.Minute
	paymentExpiryJob := paymentJob.NewPaymentExpiryJob(
		&container.ExpirePastDuePayments,
		expiryInterval,
		logger,
	)
	s.Register(paymentExpiryJob)

	return s
}

func (s *Scheduler) Register(job Job) {
	if job != nil {
		s.jobs = append(s.jobs, job)
	}
}

// Start launches all registered background jobs in separate goroutines.
func (s *Scheduler) Start(ctx context.Context) {
	for _, job := range s.jobs {
		go job.Start(ctx)
	}
}
