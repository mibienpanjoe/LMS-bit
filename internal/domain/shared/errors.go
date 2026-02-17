package shared

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrCopyNotAvailable   = errors.New("copy is not available")
	ErrMemberNotEligible  = errors.New("member is not eligible to borrow")
	ErrLoanLimitReached   = errors.New("member has reached active loan limit")
	ErrLoanAlreadyClosed  = errors.New("loan is already returned")
	ErrRenewalLimit       = errors.New("renewal limit reached")
	ErrLoanAlreadyOverdue = errors.New("overdue loan cannot be renewed")
)
