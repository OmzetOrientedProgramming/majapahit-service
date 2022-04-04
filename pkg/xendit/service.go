package xendit

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/client"
	"github.com/xendit/xendit-go/disbursement"
	"github.com/xendit/xendit-go/invoice"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type service struct {
	Client *client.API
}

// Service for calling xendit
type Service interface {
	CreateInvoice(params CreateInvoiceParams) (*xendit.Invoice, error)
	CreateDisbursement(params CreateDisbursementParams) (*xendit.Disbursement, error)
	GetInvoice(ID string) (*xendit.Invoice, error)
	GetDisbursement(ID string) (*xendit.Disbursement, error)
}

// NewXenditClient for initialize xendit service
func NewXenditClient(api *client.API) Service {
	return &service{Client: api}
}

func (x service) CreateInvoice(params CreateInvoiceParams) (*xendit.Invoice, error) {
	var xenditItems []xendit.InvoiceItem
	for _, item := range params.Items {
		xenditItems = append(xenditItems, xendit.InvoiceItem{
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Qty,
		})
	}

	withBookingFee := util.XenditFeesDefault
	withBookingFee = append(withBookingFee, xendit.InvoiceFee{
		Type:  "Booking Fee",
		Value: params.BookingFee,
	})

	invoiceParams := &invoice.CreateParams{
		ExternalID:  strconv.Itoa(params.PlaceID),
		Description: params.Description,
		Customer: xendit.InvoiceCustomer{
			GivenNames:   params.CustomerName,
			MobileNumber: params.CustomerPhoneNumber,
		},
		PaymentMethods:  util.DefaultPaymentMethod,
		Items:           xenditItems,
		Fees:            withBookingFee,
		InvoiceDuration: util.InvoiceDuration,
		CustomerNotificationPreference: xendit.InvoiceCustomerNotificationPreference{
			InvoiceCreated:  util.SMSNotification,
			InvoiceReminder: util.SMSNotification,
			InvoicePaid:     util.SMSNotification,
			InvoiceExpired:  util.SMSNotification,
		},
	}

	totalAmount := 0.0
	for _, xenditItem := range invoiceParams.Items {
		totalAmount += xenditItem.Price * float64(xenditItem.Quantity)
	}
	for _, xenditFee := range invoiceParams.Fees {
		totalAmount += xenditFee.Value
	}
	invoiceParams.Amount = totalAmount

	resp, err := x.Client.Invoice.Create(invoiceParams)
	if err != nil {
		return nil, errors.Wrap(ErrXenditCreateInvoice, err.Error())
	}

	return resp, nil
}

func (x service) CreateDisbursement(params CreateDisbursementParams) (*xendit.Disbursement, error) {
	disbursementParams := disbursement.CreateParams{
		IdempotencyKey:    time.Now().String(),
		ExternalID:        strconv.Itoa(params.ID),
		BankCode:          util.BankBCA,
		AccountHolderName: params.BankAccountName,
		AccountNumber:     params.BankAccountNumber,
		Description:       params.Description,
		Amount:            params.Amount,
		EmailTo:           params.Email,
		EmailBCC:          []string{util.OOPEmail},
	}

	resp, err := x.Client.Disbursement.Create(&disbursementParams)
	if err != nil {
		return nil, errors.Wrap(ErrXenditCreateDisbursement, err.Error())
	}

	return resp, nil
}

func (x service) GetInvoice(ID string) (*xendit.Invoice, error) {
	params := invoice.GetParams{
		ID: ID,
	}

	resp, err := x.Client.Invoice.Get(&params)
	if err != nil {
		return nil, errors.Wrap(ErrXenditGetInvoice, err.Error())
	}

	return resp, nil
}

func (x service) GetDisbursement(ID string) (*xendit.Disbursement, error) {
	params := disbursement.GetByIDParams{
		DisbursementID: ID,
	}

	resp, err := x.Client.Disbursement.GetByID(&params)
	if err != nil {
		return nil, errors.Wrap(ErrXenditGetDisbursement, err.Error())
	}

	return resp, nil
}
