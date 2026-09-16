package media

import "errors"

var (
	ErrAlreadyProcessing = errors.New("media.ingest.alreadyProcessing")
	ErrAlreadyProcessed  = errors.New("media.ingest.alreadyProcessed")
	ErrIngestFailed      = errors.New("media.ingest.failed")
)

var (
	errInvalidDimension = errors.New("media.ingest.invalidDimension")
	errInvalidFormat    = errors.New("media.ingest.invalidFormat")
)
