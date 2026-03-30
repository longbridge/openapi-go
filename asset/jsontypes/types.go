package jsontypes

// StatementItem is the JSON representation of a statement list entry.
type StatementItem struct {
	Dt      int32  `json:"dt"`
	FileKey string `json:"file_key"`
}

// StatementListResponse is the JSON response for the list statements API.
type StatementListResponse struct {
	List []*StatementItem `json:"list"`
}

// StatementDownloadURLResponse is the JSON response for the download URL API.
type StatementDownloadURLResponse struct {
	URL string `json:"url"`
}
