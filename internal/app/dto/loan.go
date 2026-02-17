package dto

type IssueLoanInput struct {
	CopyID   string
	MemberID string
}

type RenewLoanInput struct {
	LoanID string
}

type ReturnLoanInput struct {
	LoanID string
}
