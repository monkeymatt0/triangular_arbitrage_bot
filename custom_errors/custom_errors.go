package custom_errors

import "fmt"

const (
	WRONG_SIDE = "Your side is not mapped"
)

type ExchangeError struct {
	Context string
}

// ERRORS
type SideError struct {
	// Embedded struct
	Err ExchangeError
}

// Errors function implementations
func (se *SideError) Error() string {
	return fmt.Sprintf("----------\ncontext: %s\nside error: %s\nplease use one of the following: ['BUY', 'SELL']\n----------\n", se.Err.Context, WRONG_SIDE)
}
