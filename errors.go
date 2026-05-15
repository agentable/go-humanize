package humanize

import "errors"

// ErrInvalid reports malformed ParseBytes or ParseDuration input.
var ErrInvalid = errors.New("invalid input")
