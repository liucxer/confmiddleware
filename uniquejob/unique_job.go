package uniquejob

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type Locker interface {
	Lock(key string, expiresIn time.Duration) (bool, error)
}

func NewUniqueJob(ctx context.Context, locker Locker, key string, expiresIn time.Duration, do func()) *UniqueJob {
	return &UniqueJob{
		key:       key,
		expiresIn: expiresIn,
		locker:    locker,
		do:        do,
		ctx:       ctx,
	}
}

type UniqueJob struct {
	key       string
	expiresIn time.Duration
	locker    Locker
	ctx       context.Context
	do        func()
}

func (job *UniqueJob) Run() {
	ok, err := job.locker.Lock(job.key, job.expiresIn)
	if err != nil {
		ctx := job.ctx
		if ctx != nil {
			ctx = context.Background()
		}
		logrus.WithContext(ctx).Error(err)
		return
	}

	if !ok {
		return
	}

	job.do()
}
