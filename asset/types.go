package asset

// StatementType represents the type of statement.
type StatementType int32

const (
	// StatementTypeDaily is a daily statement.
	StatementTypeDaily StatementType = 1
	// StatementTypeMonthly is a monthly statement.
	StatementTypeMonthly StatementType = 2
)

// StatementItem represents a single statement entry from the list API.
type StatementItem struct {
	// Date of the statement (integer, e.g. 20250301).
	Date int32
	// FileKey used to request the download URL.
	FileKey string
}

// GetStatementList contains parameters for listing statements.
type GetStatementList struct {
	// StatementType: 1 = daily (default), 2 = monthly.
	StatementType StatementType
	// Page number for pagination.
	Page int32
	// PageSize is the number of results per page.
	PageSize int32
}

// GetStatementDownloadURL contains parameters for getting a statement download URL.
type GetStatementDownloadURL struct {
	// FileKey obtained from the list statements endpoint.
	FileKey string
}
