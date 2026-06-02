package counter

import (
	"fmt"
	"strings"
)

// IDToSymbol converts a counter_id like "ST/US/TSLA" to a symbol like "TSLA.US".
// Returns the input unchanged if it is not in the expected format.
func IDToSymbol(counterID string) string {
	parts := strings.SplitN(counterID, "/", 3)
	if len(parts) == 3 {
		return fmt.Sprintf("%s.%s", parts[2], parts[1])
	}
	return counterID
}
