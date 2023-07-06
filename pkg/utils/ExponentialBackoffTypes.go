package utils


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