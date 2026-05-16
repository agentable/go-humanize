package humanize

type invalidError string

func (e invalidError) Error() string {
	return string(e)
}

// ErrInvalid reports malformed ParseBytes or ParseDuration input.
const ErrInvalid invalidError = "invalid input"
