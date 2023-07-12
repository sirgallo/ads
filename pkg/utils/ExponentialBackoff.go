package utils

import "errors"
import "math"
import "time"


const DefaultMaxRetries = -1 // let's use this to represent unlimited retries

type ExpBackoffOpts struct {
	TimeoutInMicroseconds int
	MaxRetries *int // optional field, use a pointer
}

type ExponentialBackoffStrat struct {
	depth int
	initialTimeout int
	currentTimeout int
	maxRetries *int
}


func NewExponentialBackoffStrat(opts ExpBackoffOpts) *ExponentialBackoffStrat {
	maxRetries := DefaultMaxRetries
	if opts.MaxRetries != nil {
		maxRetries = *opts.MaxRetries
	}

	return &ExponentialBackoffStrat{
		depth: 1, 
		initialTimeout: opts.TimeoutInMicroseconds,
		currentTimeout: opts.TimeoutInMicroseconds,
		maxRetries: &maxRetries,
	}
}

func (expStrat *ExponentialBackoffStrat) PerformBackoff() error {
	if expStrat.depth > *expStrat.maxRetries && *expStrat.maxRetries != DefaultMaxRetries { 
		return errors.New("process reached max retries on exponential backoff") 
	}

	time.Sleep(time.Duration(expStrat.currentTimeout) * time.Microsecond)

	expStrat.currentTimeout = int(math.Pow(float64(2), float64(expStrat.depth - 1))) * expStrat.currentTimeout
	expStrat.depth++

	return nil
}

func (expStrat *ExponentialBackoffStrat) Reset() {
	expStrat.depth = 1
	expStrat.currentTimeout = expStrat.initialTimeout
}