package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"cloud.google.com/go/civil"
	"github.com/gocarina/gocsv"
	"github.com/siddmoitra/btech-minicash-2-poc/backend/domain"
	"github.com/siddmoitra/btech-minicash-2-poc/utils"
)

var loanApplications []*domain.LoanApplication
var loanPayments []*domain.LoanPayment
var interestApportionments []*domain.InterestApportionment
var userProfiles []*domain.UserProfile

func calculateEMI(principle float64, tenor int) (float64, float64) {
	for _, ia := range interestApportionments {
		if ia.CustomerSegment == "A076X" && ia.Tenor == tenor {
			roi := ia.Percent
			rate := roi / 100
			formula_1_plus_r_power_n := math.Pow(1+rate, float64(tenor))
			emi := principle *
				rate *
				(formula_1_plus_r_power_n / (formula_1_plus_r_power_n - 1))
			return utils.FloatToMoney(math.Round(emi)), roi
		}
	}

	return 0.0, 0.0
}

func MiniCashLoadData() {
	if loanApplications == nil {
		bytes, err := utils.LoadJSON("loan_applications.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bytes, &loanApplications)
	}

	if loanPayments == nil {
		bytes, err := utils.LoadJSON("loan_payments.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bytes, &loanPayments)
	}

	if userProfiles == nil {
		bytes, err := utils.LoadJSON("user_profiles.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(bytes, &userProfiles)
	}

	if interestApportionments == nil {
		bytes := utils.LoadCSV("interest_apportionments.csv")

		if err := gocsv.UnmarshalString(bytes, &interestApportionments); err != nil { // Load products from file
			panic(err)
		}
	}
}

func MiniCashWriteData() {
	bytesLoanApplications, _ := json.MarshalIndent(loanApplications, "", " ")
	bytesLoanPayments, _ := json.MarshalIndent(loanPayments, "", " ")
	bytesUserProfiles, _ := json.MarshalIndent(userProfiles, "", " ")
	utils.WriteJSON(bytesLoanApplications, "loan_applications.json")
	utils.WriteJSON(bytesLoanPayments, "loan_payments.json")
	utils.WriteJSON(bytesUserProfiles, "user_profiles.json")
}

func CreateUserProfile(userProfile *domain.UserProfile) {
	userProfiles = append(userProfiles, userProfile)
}

func GetUserProfileByID(id string) *domain.UserProfile {
	for _, up := range userProfiles {
		if up.ID == id {
			return up
		}
	}
	return nil
}

func AssignUserCreditForUser(userProfile *domain.UserProfile) {
	// 50% of annual income
	userProfile.CreditLimit = math.Round((0.5 * float64(userProfile.AnnualIncome) * 100) / 100)
}

func GetLoanApplicationsForUser(userProfile *domain.UserProfile) []*domain.LoanApplication {
	loanApplicationsForUser := make([]*domain.LoanApplication, 0)
	for _, la := range loanApplications {
		if la.UserID == userProfile.ID {
			loanApplicationsForUser = append(loanApplicationsForUser, la)
		}
	}

	return loanApplicationsForUser
}

func GetUserCreditWallet(userProfile *domain.UserProfile) *domain.EWallet {
	loanApplicationsForUser := GetLoanApplicationsForUser(userProfile)

	totalPaymentsTowardsLoans := 0.0
	for _, la := range loanApplicationsForUser {
		totalPaymentsTowardsLoans += la.TotalPayment
	}

	eWallet := new(domain.EWallet)
	eWallet.CreditLimit = userProfile.CreditLimit
	eWallet.UsedCredit = totalPaymentsTowardsLoans
	eWallet.LoyaltyPoints = int(math.Round((float64(eWallet.UsedCredit) * 100) / 100))

	return eWallet
}

func CreateLoanApplication(salesOrder *domain.SalesOrder, userProfile *domain.UserProfile, tenor int) *domain.LoanApplication {

	totalDiscountedPrice := 0.0
	for _, p := range salesOrder.LineItems {
		totalDiscountedPrice += utils.DiscountedPrice(p.Product.Price, p.Product.DiscountPercentMiniCash)
	}

	now := time.Now()

	loanApplication := new(domain.LoanApplication)
	loanApplication.ID = fmt.Sprint(time.Now().UTC().UnixNano())
	loanApplication.StartDate = civil.DateOf(time.Now())
	loanApplication.EndDate = civil.DateOf(now.AddDate(0, tenor, 0))
	loanApplication.SalesOrder = *salesOrder
	loanApplication.TenorMonths = tenor
	loanApplication.TotalPayment = utils.FloatToMoney(totalDiscountedPrice)
	loanApplication.UserID = userProfile.ID
	loanApplication.InstallmentValue, loanApplication.RateOfInterest = calculateEMI(totalDiscountedPrice, tenor)

	loanApplications = append(loanApplications, loanApplication)

	return loanApplication
}
