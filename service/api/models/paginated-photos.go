package models

// PaginatedPhotos - A list of paginated photo entries
type PaginatedPhotos struct {

	// The page of items for this page. This will be an empty array if there are no results.
	Entries *[]Photo `json:"entries,omitempty"`

	// The offset used for this page of results
	Offset int `json:"offset,omitempty"`

	// The limit used for this page of results. This will be the same as the `limit` query parameter unless it exceeded the maximum value allowed for this API endpoint.
	Limit int `json:"limit,omitempty"`

	// One greater than the `offset` of the last item in the entire collection. The total number of items in the collection may be less than `totalCount`
	TotalCount int `json:"totalCount"`
}
