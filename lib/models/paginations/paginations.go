package models

// PaginationMetadata represents a notification message structure
type PaginationMetadata struct {
	A_CurrentPage       int
	B_Limit             int
	C_TotalItems        int64
	D_TotalPages        int
	E_NextPageIndex     *int
	F_NextPageURL       string
	G_PreviousPageIndex *int
	H_PreviousPageURL   string
	I_PageNumbers       []int
}
