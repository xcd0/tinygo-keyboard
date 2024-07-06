package keyboard

import "time"

type Options struct {
	InvertButtonState bool
	InvertDiode       bool
	MatrixScanPeriod  time.Duration
}

type Option func(*Options)

func InvertButtonState(b bool) Option {
	return func(o *Options) {
		o.InvertButtonState = b
	}
}

func InvertDiode(b bool) Option {
	return func(o *Options) {
		o.InvertDiode = b
	}
}

// MatrixScanPeriod sets the total period for scanning the entire keyboard matrix.
// This period determines the overall frequency at which the complete keyboard state is updated.
func MatrixScanPeriod(period time.Duration) Option {
	return func(o *Options) {
		o.MatrixScanPeriod = period
	}
}
