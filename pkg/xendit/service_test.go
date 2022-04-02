package xendit

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/xendit/xendit-go/client"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"os"
	"strconv"
	"testing"
)

func TestService_CreateInvoiceSuccess(t *testing.T) {
	_ = godotenv.Load("../../.env")
	params := CreateInvoiceParams{
		PlaceID: 1,
		Items: []Item{
			{
				Name:  "test item",
				Price: 10000,
				Qty:   2,
			},
			{
				Name:  "test item",
				Price: 10000,
				Qty:   2,
			},
		},
		Description:         "test description",
		CustomerName:        "test customer name",
		CustomerPhoneNumber: "+628123456712",
	}

	totalAmountExpected := 0.0
	for _, item := range params.Items {
		totalAmountExpected += item.Price * float64(item.Qty)
	}

	for _, xenditFee := range util.XenditFeesDefault {
		totalAmountExpected += xenditFee.Value
	}

	xenCli := client.New(os.Getenv("XENDIT_TOKEN"))
	testService := NewXenditClient(xenCli)
	resp, err := testService.CreateInvoice(params)

	assert.NoError(t, err)
	assert.Equal(t, totalAmountExpected, resp.Amount)
	assert.Equal(t, params.CustomerName, resp.Customer.GivenNames)
	assert.Equal(t, params.CustomerPhoneNumber, resp.Customer.MobileNumber)
	assert.Equal(t, params.Description, resp.Description)
	assert.Equal(t, strconv.Itoa(params.PlaceID), resp.ExternalID)

	for i := 0; i < len(resp.Items); i++ {
		assert.Equal(t, params.Items[i].Name, resp.Items[i].Name)
		assert.Equal(t, params.Items[i].Price, resp.Items[i].Price)
		assert.Equal(t, params.Items[i].Qty, resp.Items[i].Quantity)
	}

	for i := 0; i < len(resp.Fees); i++ {
		assert.Equal(t, util.XenditFeesDefault[i].Value, resp.Fees[i].Value)
		assert.Equal(t, util.XenditFeesDefault[i].Type, resp.Fees[i].Type)
	}
}

func TestService_CreateInvoiceFailed(t *testing.T) {
	_ = godotenv.Load("../../.env")
	params := CreateInvoiceParams{
		PlaceID: 1,
		Items: []Item{
			{
				Name:  "test item",
				Price: 2000,
				Qty:   2,
			},
		},
		Description:         "test description",
		CustomerName:        "test customer name",
		CustomerPhoneNumber: "0812121212",
	}

	xenCli := client.New("failedSecretKey")
	testService := NewXenditClient(xenCli)
	resp, err := testService.CreateInvoice(params)

	assert.Nil(t, resp)
	assert.Equal(t, ErrXenditCreateInvoice, errors.Cause(err))
}

func TestService_CreateDisbursementSuccess(t *testing.T) {
	_ = godotenv.Load("../../.env")
	params := CreateDisbursementParams{
		ID:                1,
		BankAccountName:   "test bank account name",
		BankAccountNumber: "123412",
		Amount:            10000,
		Description:       "test description",
		Email:             []string{"test@email.com"},
	}

	xenCli := client.New(os.Getenv("XENDIT_TOKEN"))
	testService := NewXenditClient(xenCli)
	resp, err := testService.CreateDisbursement(params)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, strconv.Itoa(params.ID), resp.ExternalID)
	assert.Equal(t, params.BankAccountName, resp.AccountHolderName)
	assert.Equal(t, params.Email, resp.EmailTo)
	assert.Equal(t, []string{util.OOPEmail}, resp.EmailBCC)
	assert.Equal(t, params.Description, resp.DisbursementDescription)
	assert.Equal(t, params.Amount, resp.Amount)
	assert.Equal(t, util.BankBCA, resp.BankCode)
}

func TestService_CreateDisbursementFailed(t *testing.T) {
	_ = godotenv.Load("../../.env")
	params := CreateDisbursementParams{
		ID:                1,
		BankAccountName:   "test bank account name",
		BankAccountNumber: "123412",
		Amount:            10000,
		Description:       "test description",
		Email:             []string{"test@email.com"},
	}

	xenCli := client.New("wrongToken")
	testService := NewXenditClient(xenCli)
	resp, err := testService.CreateDisbursement(params)

	assert.Nil(t, resp)
	assert.Equal(t, ErrXenditCreateDisbursement, errors.Cause(err))
}

func TestService_GetInvoiceSuccess(t *testing.T) {
	_ = godotenv.Load("../../.env")

	params := CreateInvoiceParams{
		PlaceID: 1,
		Items: []Item{
			{
				Name:  "test item",
				Price: 10000,
				Qty:   2,
			},
			{
				Name:  "test item",
				Price: 10000,
				Qty:   2,
			},
		},
		Description:         "test description",
		CustomerName:        "test customer name",
		CustomerPhoneNumber: "+628123456712",
	}

	totalAmountExpected := 0.0
	for _, item := range params.Items {
		totalAmountExpected += item.Price * float64(item.Qty)
	}

	for _, xenditFee := range util.XenditFeesDefault {
		totalAmountExpected += xenditFee.Value
	}

	xenCli := client.New(os.Getenv("XENDIT_TOKEN"))
	testService := NewXenditClient(xenCli)

	resp, _ := testService.CreateInvoice(params)
	resp, err := testService.GetInvoice(resp.ID)

	assert.NoError(t, err)
	assert.Equal(t, totalAmountExpected, resp.Amount)
	assert.Equal(t, params.CustomerName, resp.Customer.GivenNames)
	assert.Equal(t, params.CustomerPhoneNumber, resp.Customer.MobileNumber)
	assert.Equal(t, params.Description, resp.Description)
	assert.Equal(t, strconv.Itoa(params.PlaceID), resp.ExternalID)

	for i := 0; i < len(resp.Items); i++ {
		assert.Equal(t, params.Items[i].Name, resp.Items[i].Name)
		assert.Equal(t, params.Items[i].Price, resp.Items[i].Price)
		assert.Equal(t, params.Items[i].Qty, resp.Items[i].Quantity)
	}

	for i := 0; i < len(resp.Fees); i++ {
		assert.Equal(t, util.XenditFeesDefault[i].Value, resp.Fees[i].Value)
		assert.Equal(t, util.XenditFeesDefault[i].Type, resp.Fees[i].Type)
	}
}

func TestService_GetInvoiceFailed(t *testing.T) {
	xenCli := client.New(os.Getenv("XENDIT_TOKEN"))
	testService := NewXenditClient(xenCli)
	resp, err := testService.GetInvoice("WrongID")

	assert.Nil(t, resp)
	assert.Equal(t, ErrXenditGetInvoice, errors.Cause(err))
}

func TestService_GetDisbursementSuccess(t *testing.T) {
	_ = godotenv.Load("../../.env")
	params := CreateDisbursementParams{
		ID:                1,
		BankAccountName:   "test bank account name",
		BankAccountNumber: "123412",
		Amount:            10000,
		Description:       "test description",
		Email:             []string{"test@email.com"},
	}

	xenCli := client.New(os.Getenv("XENDIT_TOKEN"))
	testService := NewXenditClient(xenCli)
	resp, _ := testService.CreateDisbursement(params)
	resp, err := testService.GetDisbursement(resp.ID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, strconv.Itoa(params.ID), resp.ExternalID)
	assert.Equal(t, params.BankAccountName, resp.AccountHolderName)
	assert.Equal(t, params.Email, resp.EmailTo)
	assert.Equal(t, []string{util.OOPEmail}, resp.EmailBCC)
	assert.Equal(t, params.Description, resp.DisbursementDescription)
	assert.Equal(t, params.Amount, resp.Amount)
	assert.Equal(t, util.BankBCA, resp.BankCode)
}

func TestService_GetDisbursementFailed(t *testing.T) {
	xenCli := client.New(os.Getenv("XENDIT_TOKEN"))
	testService := NewXenditClient(xenCli)
	resp, err := testService.GetDisbursement("WrongID")

	assert.Nil(t, resp)
	assert.Equal(t, ErrXenditGetDisbursement, errors.Cause(err))
}
