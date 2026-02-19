package service

import (
	"context"
	"io"

	"github.com/eenemeene/kitamanager-go/internal/isbj"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

// GovernmentFundingBillService handles government funding bill file processing.
type GovernmentFundingBillService struct {
	childStore store.ChildStorer
}

// NewGovernmentFundingBillService creates a new GovernmentFundingBillService.
func NewGovernmentFundingBillService(childStore store.ChildStorer) *GovernmentFundingBillService {
	return &GovernmentFundingBillService{childStore: childStore}
}

// ProcessISBJ parses an ISBJ Excel file and returns enriched government funding bill data.
func (s *GovernmentFundingBillService) ProcessISBJ(ctx context.Context, orgID uint, reader io.Reader) (*models.GovernmentFundingBillResponse, error) {
	output, err := isbj.ParseFromReader(reader)
	if err != nil {
		return nil, err
	}

	converted, err := isbj.Convert(output)
	if err != nil {
		return nil, err
	}

	// Collect voucher numbers for matching
	voucherNumbers := make([]string, 0, len(converted.Children))
	for _, child := range converted.Children {
		voucherNumbers = append(voucherNumbers, child.VoucherNumber)
	}

	// Look up contracts by voucher number
	contractMap := make(map[string]models.ChildContract)
	if len(voucherNumbers) > 0 {
		contracts, err := s.childStore.FindContractsByVoucherNumbers(ctx, orgID, voucherNumbers)
		if err != nil {
			return nil, err
		}
		for _, c := range contracts {
			if c.VoucherNumber != nil {
				contractMap[*c.VoucherNumber] = c
			}
		}
	}

	// Build response
	matchedCount := 0
	children := make([]models.GovernmentFundingBillChildResponse, 0, len(converted.Children))
	for _, child := range converted.Children {
		resp := models.GovernmentFundingBillChildResponse{
			VoucherNumber: child.VoucherNumber,
			ChildName:     child.ChildName,
			BirthDate:     child.BirthDate,
			District:      child.District,
			TotalAmount:   child.TotalAmount,
			Amounts:       convertBillAmounts(child.Amounts),
		}

		if contract, ok := contractMap[child.VoucherNumber]; ok {
			resp.ChildID = &contract.ChildID
			resp.ContractID = &contract.ID
			resp.Matched = true
			matchedCount++
		}

		children = append(children, resp)
	}

	return &models.GovernmentFundingBillResponse{
		FacilityName:      converted.FacilityName,
		FacilityTotal:     converted.FacilityTotal,
		ContractBooking:   converted.ContractBooking,
		CorrectionBooking: converted.CorrectionBooking,
		ChildrenCount:     converted.ChildrenCount,
		MatchedCount:      matchedCount,
		UnmatchedCount:    converted.ChildrenCount - matchedCount,
		Surcharges:        convertBillAmounts(converted.Surcharges),
		Children:          children,
	}, nil
}

func convertBillAmounts(amounts []isbj.SettlementAmount) []models.GovernmentFundingBillAmount {
	result := make([]models.GovernmentFundingBillAmount, len(amounts))
	for i, a := range amounts {
		result[i] = models.GovernmentFundingBillAmount{
			Key:    a.Key,
			Value:  a.Value,
			Amount: a.Amount,
		}
	}
	return result
}
