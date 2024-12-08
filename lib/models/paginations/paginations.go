package models

// PaginationMetadata represents a notification message structure
type PaginationMetadata struct {
	A_CurrentPage int
	B_Limit       int
	C_TotalItems  int64
	D_TotalPages  int
	// G_PreviousPageIndex *int
	H_PreviousPageURL string
	// E_NextPageIndex     *int
	F_NextPageURL string
	I_Pages       []PageInfo
}

type PageInfo struct {
	Index    int
	Page_url string
}
