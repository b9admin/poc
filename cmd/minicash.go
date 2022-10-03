package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"

	"github.com/AlecAivazis/survey/v2"

	"github.com/siddmoitra/btech-minicash-2-poc/backend/domain"
	"github.com/siddmoitra/btech-minicash-2-poc/backend/handler"
	"github.com/siddmoitra/btech-minicash-2-poc/utils"
)

func displaySalesOrder(salesOrder *domain.SalesOrder) {

	fmt.Println("Sales Order Details")
	fmt.Println("-------------------")

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Product Name", "Sales Price (in EGP)", "MiniCash Promotional Discount (%)", "Buying Price (in EGP)")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	totalPriceOnDiscount := 0.0
	for _, lineItem := range salesOrder.LineItems {
		totalPriceOnDiscount += utils.DiscountedPrice(lineItem.Product.Price, lineItem.Product.DiscountPercentMiniCash)
		tbl.AddRow(lineItem.Product.Name, utils.FloatToMoney(lineItem.Product.Price), lineItem.Product.DiscountPercentMiniCash, utils.DiscountedPrice(lineItem.Product.Price, lineItem.Product.DiscountPercentMiniCash))
	}

	tbl.AddRow("Total Order Price:", utils.FloatToMoney(salesOrder.TotalOrderPrice), "-", utils.FloatToMoney(totalPriceOnDiscount))
	tbl.Print()
	fmt.Println("")
}

func displayEWallet(eWallet *domain.EWallet) {
	fmt.Println("User E-Wallet")
	fmt.Println("-------------")

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("", "")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.AddRow("Max Credit Limit", fmt.Sprintf("%f EGP", utils.FloatToMoney(eWallet.CreditLimit)))
	tbl.AddRow("Used Credit", fmt.Sprintf("%f EGP", utils.FloatToMoney(eWallet.UsedCredit)))
	tbl.AddRow("Loyalty Points", fmt.Sprintf("%d", eWallet.LoyaltyPoints))
	tbl.Print()

	fmt.Println("")
}

func displayLoanApplication(loanApplication *domain.LoanApplication) {
	fmt.Println("Loan Application")
	fmt.Println("----------------")

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("", "")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.AddRow("Loan Amount (Pinciple in EGP)", fmt.Sprintf("%f EGP", utils.FloatToMoney(loanApplication.TotalPayment)))
	tbl.AddRow("Rate of Interest (in %% per month)", utils.FloatToMoney(loanApplication.RateOfInterest))
	tbl.AddRow("Tenor (in months)", loanApplication.TenorMonths)
	tbl.AddRow("Start Date", utils.DateToString(&loanApplication.StartDate))
	tbl.AddRow("End Date", utils.DateToString(&loanApplication.StartDate))
	tbl.AddRow("Monthly Installment (in EGP)", utils.FloatToMoney(loanApplication.InstallmentValue))
	tbl.AddRow("Total Payment (in EGP)", utils.FloatToMoney(loanApplication.InstallmentValue*float64(loanApplication.TenorMonths)))
	tbl.AddRow("Extra Payment towards Interest (in EGP)", loanApplication.InstallmentValue*float64(loanApplication.TenorMonths)-loanApplication.TotalPayment)
	tbl.Print()

	fmt.Println("")
}

func displayLoanApplicationsOfUserMinimal(loanApplications []*domain.LoanApplication) {
	fmt.Println("Loan Applications")
	fmt.Println("-----------------")

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Application ID", "Total Order Value", "Number of Paid Installments", "Start Date", "End Date")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, la := range loanApplications {
		tbl.AddRow(la.ID, fmt.Sprintf("%f EGP", utils.FloatToMoney(la.TotalPayment)), 0, utils.DateToString(&la.StartDate), utils.DateToString(&la.EndDate))
	}

	tbl.Print()

	fmt.Println("")
}

func questionExistingCustomer() (bool, error) {
	existingCustomer := false
	questionExistingCustomer := []*survey.Question{
		{
			Name: "existingCustomer",
			Prompt: &survey.Confirm{
				Message: "Are you an existing customer ?",
			},
		},
	}
	err := survey.Ask(questionExistingCustomer, &existingCustomer)
	if err != nil {
		return existingCustomer, err
	}

	return existingCustomer, nil
}

func questionLogin() (string, error) {

	loginCredentials := struct {
		LoginId  string `survey:"loginId"`
		Password string `survey:"password"`
	}{}

	questionLogin := []*survey.Question{
		{
			Name: "loginId",
			Prompt: &survey.Input{
				Message: "Enter your customer ID or National ID to login ?",
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Enter your password ?",
			},
		},
	}
	err := survey.Ask(questionLogin, &loginCredentials)
	if err != nil {
		return "", err
	}

	return loginCredentials.LoginId, nil
}

func onboardAndGetUserProfile() (*domain.UserProfile, error) {
	userProfile := new(domain.UserProfile)
	userOnboarding := struct {
		Id           string
		Name         string
		Password     string
		DocumentType string `survey:"documentType"`
		DocumentId   string `survey:"documentId"`
		AnnualIncome int    `survey:"annualIncome"`
	}{}

	questionOnboardUser := []*survey.Question{
		{
			Name: "id",
			Prompt: &survey.Input{
				Message: "Enter your desired username ?",
			},
		},
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Enter your name ?",
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Enter your password (for future) ?",
			},
		},
		{
			Name: "documentType",
			Prompt: &survey.Select{
				Message: "Choose the document type (for identification) ?",
				Options: []string{"National ID", "Passport"},
			},
		},
		{
			Name: "documentId",
			Prompt: &survey.Input{
				Message: "Enter the document number (for identification) ?",
			},
		},
		{
			Name: "annualIncome",
			Prompt: &survey.Input{
				Message: "Enter the annual income (for credit limit) ?",
			},
		},
	}
	err := survey.Ask(questionOnboardUser, &userOnboarding)
	if err != nil {
		return nil, err
	}

	userProfile.ID = userOnboarding.Id
	userProfile.Name = userOnboarding.Name
	userProfile.LegalDocument.DocumentType = userOnboarding.DocumentType
	userProfile.LegalDocument.DocumentNumber = userOnboarding.DocumentId
	userProfile.AnnualIncome = userOnboarding.AnnualIncome

	handler.CreateUserProfile(userProfile)

	return userProfile, nil
}

func loginAndGetUserProfile() (*domain.UserProfile, error) {
	id, err := questionLogin()
	if err != nil {
		return nil, err
	}
	userProfile := handler.GetUserProfileByID(id)
	if userProfile == nil {
		fmt.Println("X Wrong user credentials. Try again.")
		return loginAndGetUserProfile()
	}

	return userProfile, nil
}

func assignUserCreditNewCustomer(userProfile *domain.UserProfile) error {
	handler.AssignUserCreditForUser(userProfile)
	return nil
}

func questionTenor() (int, error) {
	var tenor string
	questionTenor := []*survey.Question{
		{
			Name: "tenor",
			Prompt: &survey.Select{
				Message: "Select your installment plan (in months) ?",
				Options: []string{"12", "24", "36", "48"},
			},
		},
	}
	err := survey.Ask(questionTenor, &tenor)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(tenor)
}

func MiniCashInitiatePayment(salesOrder *domain.SalesOrder) (string, error) {

	utils.ShowUserInSystem("End-Consumer now in MiniCash 2.0")

	existingCustomer, err := questionExistingCustomer()
	if err != nil {
		return "", err
	}

	var userProfile *domain.UserProfile
	if existingCustomer {
		userProfile, err = loginAndGetUserProfile()
		if err != nil {
			return "", err
		}
	} else {
		fmt.Println("-> End-Consumer now in MiniCash 2.0 new user onboarding workflow")
		utils.PressEnterToContinue()

		userProfile, err = onboardAndGetUserProfile()
		if err != nil {
			return "", err
		}
	}

	if !existingCustomer {
		utils.LineBreak(2)
		utils.SpinnerMessage("  Now fetching credit risk, AML, score and other related details using the identified document", time.Second*5)
		fmt.Println("An inital customer score and segment assigned ...")
		fmt.Println("- In real world, this process might branch out to an exception path where the credit limit will be put on hold and MiniCash 2.0 will respond to the e-Commerce as payment in progress")
		fmt.Println("- The e-Commerce solution MUST then decide how to process the order, MiniCash 2.0 will NEVER make that decision for the e-Commerce, but only respond with PAYMENT_STATUS = [SUCCESS | FAIL | ON_HOLD]")
		fmt.Println("- If the payment status is ON_HOLD, MiniCash 2.0 will send a new payment status based on asynchronous process (manual or automated), a message to e-Commerce on the new payment status")
		utils.LineBreak(2)

		assignUserCreditNewCustomer(userProfile)
	}

	eWallet := handler.GetUserCreditWallet(userProfile)

	displayEWallet(eWallet)
	utils.PressEnterToContinue()
	displaySalesOrder(salesOrder)
	utils.PressEnterToContinue()

	if eWallet.CreditLimit-eWallet.UsedCredit < salesOrder.TotalOrderPrice {
		return domain.PAYMENT_FAILED, nil
	}

	tenor, err := questionTenor()
	if err != nil {
		return "", err
	}
	loanApplication := handler.CreateLoanApplication(salesOrder, userProfile, tenor)
	displayLoanApplication(loanApplication)

	return domain.PAYMENT_SUCCESS, nil
}

func MiniCashCustomerWallet() error {
	id, err := questionLogin()
	if err != nil {
		return err
	}
	userProfile := handler.GetUserProfileByID(id)

	eWallet := handler.GetUserCreditWallet(userProfile)
	displayEWallet(eWallet)

	loanApplicationsOfUser := handler.GetLoanApplicationsForUser(userProfile)
	displayLoanApplicationsOfUserMinimal(loanApplicationsOfUser)

	return nil
}
