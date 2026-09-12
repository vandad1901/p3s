package media

import "errors"

var (
	errAlreadyProcessing = errors.New("media.ingest.alreadyProcessing")
	errAlreadyProcessed  = errors.New("media.ingest.alreadyProcessed")
	errIngestFailed      = errors.New("media.ingest.failed")
)

var (
	errInvalidDimension = errors.New("media.ingest.invalidDimension")
	errInvalidFormat    = errors.New("media.ingest.invalidFormat")
)
