package media

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyProcessing = errors.New("media.ingest.alreadyProcessing")
	ErrAlreadyProcessed  = errors.New("media.ingest.alreadyProcessed")
	ErrIngestFailed      = errors.New("media.ingest.failed")
)

var (
	errInvalidDimension = fmt.Errorf("%w: invalid dimensions", ErrIngestFailed)
	errInvalidFormat    = fmt.Errorf("%w: invalid format", ErrIngestFailed)
)
