package booking

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/xendit"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Service interface for define function in service
type Service interface {
	GetListCustomerBookingWithPagination(params ListRequest) (*ListBooking, *util.Pagination, error)
	GetAvailableTime(params GetAvailableTimeParams) (*[]AvailableTimeResponse, error)
	GetAvailableDate(params GetAvailableDateParams) (*[]AvailableDateResponse, error)
	CreateBooking(params CreateBookingServiceRequest) (*CreateBookingServiceResponse, error)
	GetTimeSlots(placeID int, selectedDate time.Time) (*[]TimeSlot, error)
	GetDetail(bookingID int) (*Detail, error)
	UpdateBookingStatus(bookingID int, newStatus int) error
	GetMyBookingsOngoing(localID string) (*[]Booking, error)
	GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, *util.Pagination, error)
	UpdateBookingStatusByXendit(callback XenditInvoicesCallback) error
}

type service struct {
	repo   Repo
	xendit xendit.Service
}

// NewService for initialize service
func NewService(repo Repo, xendit xendit.Service) Service {
	return &service{
		repo:   repo,
		xendit: xendit,
	}
}

func (s service) GetListCustomerBookingWithPagination(params ListRequest) (*ListBooking, *util.Pagination, error) {
	var errorList []string

	if params.State < 0 || params.State > 5 {
		params.State = 0
	}

	if params.Page == 0 {
		params.Page = util.DefaultPage
	}

	if params.Limit == 0 {
		params.Limit = util.DefaultLimit
	}

	if params.Limit > util.MaxLimit {
		errorList = append(errorList, "limit should be 1 - 100")
	}

	if params.Path == "" {
		errorList = append(errorList, "path is required for pagination")
	}

	if len(errorList) > 0 {
		return nil, nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	listCustomerBooking, err := s.repo.GetListCustomerBookingWithPagination(params)
	if err != nil {
		return nil, nil, err
	}
	pagination := util.GeneratePagination(listCustomerBooking.TotalCount, params.Limit, params.Page, params.Path)

	return listCustomerBooking, &pagination, err
}

func (s service) CreateBooking(params CreateBookingServiceRequest) (*CreateBookingServiceResponse, error) {
	var (
		errorList []string
		err       error
	)

	if params.Count <= 0 {
		errorList = append(errorList, "count should be positive integer")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	var items []CheckedItemParams
	for _, item := range params.Items {
		items = append(items, CheckedItemParams{
			ID:      item.ID,
			PlaceID: params.PlaceID,
		})
	}

	// item validation
	var checkedItems *[]CheckedItemParams
	var isMatch bool
	if len(items) > 0 {
		checkedItems, isMatch, err = s.repo.CheckedItem(items)
		if !isMatch && errors.Cause(err) == ErrInputValidationError {
			diff := s.difference(items, *checkedItems)
			errorMessage := make([]string, 0)

			for _, i := range diff {
				errorMessage = append(errorMessage, strconv.Itoa(i.ID))
			}

			return nil, errors.Wrap(ErrInputValidationError, fmt.Sprintf("item with id %s is not found", strings.Join(errorMessage, ", ")))
		}

		if err != nil {
			return nil, err
		}
	}

	// Selected date time validation
	getAvaialableTimeParams := GetAvailableTimeParams{
		PlaceID:      params.PlaceID,
		SelectedDate: params.Date,
		StartTime:    params.StartTime,
		BookedSlot:   params.Count,
	}

	availableTime, err := s.GetAvailableTime(getAvaialableTimeParams)
	if err != nil {
		return nil, err
	}

	isExist := false
	for _, i := range *availableTime {
		if i.Time == params.EndTime.Format(util.TimeLayout) {
			isExist = true
		}
	}

	if !isExist {
		return nil, errors.Wrap(ErrInputValidationError, "selected date time is not available for booking")
	}

	// Create booking
	bookingParams := CreateBookingParams{
		UserID:     params.UserID,
		PlaceID:    params.PlaceID,
		Date:       params.Date,
		StartTime:  params.StartTime,
		EndTime:    params.EndTime,
		Capacity:   params.Count,
		Status:     util.BookingMenungguKonfirmasi,
		TotalPrice: 0,
	}

	// create booking instance
	bookingID, err := s.repo.CreateBooking(bookingParams)
	if err != nil {
		return nil, err
	}

	bookingPrice, err := s.repo.GetPlaceBookingPrice(params.PlaceID)
	if err != nil {
		return nil, err
	}

	if checkedItems != nil && isMatch {
		// convert items to booking items & xendit items instance
		var bookingItems []CreateBookingItemsParams
		var xenditItems []xendit.Item

		for _, i := range params.Items {
			bookingItems = append(bookingItems, CreateBookingItemsParams{
				BookingID:  bookingID.ID,
				ItemID:     i.ID,
				TotalPrice: i.Price * float64(i.Qty),
				Qty:        i.Qty,
			})

			xenditItems = append(xenditItems, xendit.Item{
				Name:  i.Name,
				Price: i.Price,
				Qty:   i.Qty,
			})
		}

		// create booking items
		totalPrice, err := s.repo.CreateBookingItems(bookingItems)
		if err != nil {
			return nil, err
		}

		// update total price
		_, err = s.repo.UpdateTotalPrice(UpdateTotalPriceParams{
			BookingID:  bookingID.ID,
			TotalPrice: totalPrice.TotalPrice,
		})
		if err != nil {
			return nil, err
		}

		// create xendit invoices
		invoiceParams := xendit.CreateInvoiceParams{
			PlaceID:             params.PlaceID,
			Items:               xenditItems,
			Description:         fmt.Sprintf("order from %s", params.CustomerName),
			CustomerName:        params.CustomerName,
			CustomerPhoneNumber: params.CustomerPhoneNumber,
			BookingFee:          bookingPrice,
		}

		invoice, err := s.xendit.CreateInvoice(invoiceParams)
		if err != nil {
			return nil, err
		}

		xenditInformationParams := XenditInformation{
			XenditID:    invoice.ID,
			InvoicesURL: invoice.InvoiceURL,
			BookingID:   bookingID.ID,
		}
		_, err = s.repo.InsertXenditInformation(xenditInformationParams)
		if err != nil {
			return nil, err
		}

		return &CreateBookingServiceResponse{
			XenditID:   invoice.ID,
			BookingID:  bookingID.ID,
			PaymentURL: invoice.InvoiceURL,
		}, nil
	}

	// create xendit invoices
	invoiceParams := xendit.CreateInvoiceParams{
		PlaceID:             params.PlaceID,
		Items:               nil,
		Description:         fmt.Sprintf("order from %s", params.CustomerName),
		CustomerName:        params.CustomerName,
		CustomerPhoneNumber: params.CustomerPhoneNumber,
		BookingFee:          bookingPrice,
	}

	invoice, err := s.xendit.CreateInvoice(invoiceParams)
	if err != nil {
		return nil, err
	}

	xenditInformationParams := XenditInformation{
		XenditID:    invoice.ID,
		InvoicesURL: invoice.InvoiceURL,
		BookingID:   bookingID.ID,
	}
	_, err = s.repo.InsertXenditInformation(xenditInformationParams)
	if err != nil {
		return nil, err
	}

	return &CreateBookingServiceResponse{
		XenditID:   invoice.ID,
		BookingID:  bookingID.ID,
		PaymentURL: invoice.InvoiceURL,
	}, nil
}

func (s service) GetTimeSlots(placeID int, selectedDate time.Time) (*[]TimeSlot, error) {
	errorList := make([]string, 0)

	if placeID <= 0 {
		errorList = append(errorList, "placeID must positive integer")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	timeSlot, err := s.repo.GetTimeSlotsData(placeID, selectedDate)
	if err != nil {
		return nil, err
	}

	return timeSlot, nil
}

func (s service) GetAvailableDate(params GetAvailableDateParams) (*[]AvailableDateResponse, error) {
	errorList := make([]string, 0)

	if params.PlaceID <= 0 {
		errorList = append(errorList, "placeID must positive integer")
	}

	if params.BookedSlot <= 0 {
		errorList = append(errorList, "booking count must positive integer")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	if params.Interval == 0 {
		params.Interval = 7
	}

	var checkedDate []time.Time
	for i := 0; i <= params.Interval; i++ {
		checkedDate = append(checkedDate, params.StartDate.Add(time.Duration(i*24)*time.Hour))
	}
	endDate := checkedDate[len(checkedDate)-1]

	repoParams := GetBookingDataParams{
		PlaceID:   params.PlaceID,
		StartDate: params.StartDate,
		EndDate:   endDate,
	}

	bookingData, err := s.repo.GetBookingData(repoParams)
	if err != nil {
		return nil, err
	}

	timeSlot, err := s.repo.GetTimeSlotsData(params.PlaceID, checkedDate...)
	if err != nil {
		return nil, err
	}

	place, err := s.repo.GetPlaceCapacity(params.PlaceID)
	if err != nil {
		return nil, err
	}

	mapTimeSlot := s.makeTimeSlotsAsMap(*timeSlot)
	dividedBooking := s.divideBookings(*bookingData, mapTimeSlot, params.StartDate, params.Interval)
	availableDate := s.checkAvailableSchedule(dividedBooking, place.OpenHour, place.Capacity, params.BookedSlot, *timeSlot, true)
	formattedData := s.formatAvailableDateData(availableDate)

	return &formattedData, nil
}

func (s service) GetAvailableTime(params GetAvailableTimeParams) (*[]AvailableTimeResponse, error) {
	errorList := make([]string, 0)

	if params.PlaceID <= 0 {
		errorList = append(errorList, "placeID must positive integer")
	}

	if params.BookedSlot <= 0 {
		errorList = append(errorList, "booking count must positive integer")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	midnight := time.Date(params.SelectedDate.Year(), params.SelectedDate.Month(), params.SelectedDate.Day(), 0, 0, 0, 0, params.SelectedDate.Location())
	midnight = midnight.Add(time.Duration(1*24) * time.Hour)

	repoParams := GetBookingDataParams{
		PlaceID:   params.PlaceID,
		StartDate: params.SelectedDate,
		EndDate:   midnight,
		StartTime: params.StartTime,
	}

	bookingData, err := s.repo.GetBookingData(repoParams)
	if err != nil {
		return nil, err
	}

	timeSlot, err := s.repo.GetTimeSlotsData(params.PlaceID, params.SelectedDate)
	if err != nil {
		return nil, err
	}

	place, err := s.repo.GetPlaceCapacity(params.PlaceID)
	if err != nil {
		return nil, err
	}

	mapTimeSlot := s.makeTimeSlotsAsMap(*timeSlot)
	dividedBooking := s.divideBookings(*bookingData, mapTimeSlot, params.SelectedDate, 1)
	availableTime := s.checkAvailableSchedule(dividedBooking, params.StartTime, place.Capacity, params.BookedSlot, *timeSlot, false)

	availableTimesFormatted := s.formatAvailableTimeData(availableTime, params.SelectedDate)

	return &availableTimesFormatted, nil
}

func (s service) makeTimeSlotsAsMap(timeSlot []TimeSlot) map[int]map[time.Time]time.Time {
	mapTimeSlot := make(map[int]map[time.Time]time.Time)
	for i := 0; i < 7; i++ {
		mapTimeSlot[i] = make(map[time.Time]time.Time)
	}

	for _, i := range timeSlot {
		mapTimeSlot[i.Day][i.StartTime] = i.EndTime
	}

	return mapTimeSlot
}

func (s service) divideBookings(bookings []DataForCheckAvailableSchedule, timeSlot map[int]map[time.Time]time.Time, fromDate time.Time, checkedInterval int) map[string]map[string]int {
	result := make(map[string]map[string]int)
	for i := 0; i < checkedInterval; i++ {
		dateFormatting := fmt.Sprintf("%s", fromDate.Add(time.Duration(24*i)*time.Hour).Format(util.DateLayout))
		result[dateFormatting] = make(map[string]int)
	}

	for _, i := range bookings {
		tempTime := i
		isZero := false

		for isZero == false {
			endTimeFromTimeSlot := timeSlot[int(i.Date.Weekday())][tempTime.StartTime]
			timeFormatting := fmt.Sprintf("%s - %s", tempTime.StartTime.Format(util.TimeLayout), endTimeFromTimeSlot.Format(util.TimeLayout))
			dateFormat := i.Date.Format(util.DateLayout)
			if val, ok := result[dateFormat]; ok {
				val[timeFormatting] += tempTime.Capacity
			}

			tempTime.StartTime = tempTime.StartTime.Add(endTimeFromTimeSlot.Sub(tempTime.StartTime))
			if tempTime.StartTime.Equal(tempTime.EndTime) || tempTime.StartTime.After(tempTime.EndTime) {
				isZero = true
			}
		}
	}

	return result
}

func (s service) checkAvailableSchedule(bookings map[string]map[string]int, startTime time.Time, placeCapacity int, bookingSlot int, timeSlots []TimeSlot, isCheckedDate bool) map[string]map[string]int {
	availableBookings := make(map[string]map[string]int)
	for key, val := range bookings {
		availableBookings[key] = make(map[string]int)

		for index, timeSlot := range timeSlots {
			timeKey := fmt.Sprintf("%s - %s", timeSlot.StartTime.Format(util.TimeLayout), timeSlot.EndTime.Format(util.TimeLayout))
			oneSecond := time.Duration(1) * time.Second
			date, _ := time.Parse(util.DateLayout, key)
			day := int(date.Weekday())

			if day != timeSlot.Day {
				continue
			}

			capacityOfBookingGivenTime := val[timeKey]
			endTimeOnly := fmt.Sprintf("%s", timeSlot.EndTime.Format(util.TimeLayout))
			if index == 0 && timeSlot.StartTime.Add(oneSecond).After(startTime) && capacityOfBookingGivenTime+bookingSlot <= placeCapacity {
				availableBookings[key][endTimeOnly] = val[timeKey]
				continue
			}

			if timeSlot.StartTime.Add(oneSecond).After(startTime) {
				if capacityOfBookingGivenTime+bookingSlot <= placeCapacity && (timeSlots[index-1].EndTime.Equal(timeSlot.StartTime) || timeSlot.StartTime.Equal(startTime)) {
					availableBookings[key][endTimeOnly] = val[timeKey]
				} else {
					if isCheckedDate {
						continue
					} else {
						break
					}
				}
			}
		}
	}

	return availableBookings
}

func (s service) formatAvailableTimeData(data map[string]map[string]int, selectedDate time.Time) []AvailableTimeResponse {
	availableTimes := make([]AvailableTimeResponse, 0)
	perDayData := data[selectedDate.Format(util.DateLayout)]

	for key, val := range perDayData {
		availableTime := AvailableTimeResponse{
			Time:  key,
			Total: val,
		}

		availableTimes = append(availableTimes, availableTime)
	}

	sort.Slice(availableTimes[:], func(i, j int) bool {
		return availableTimes[i].Time < availableTimes[j].Time
	})

	return availableTimes
}

func (s service) formatAvailableDateData(data map[string]map[string]int) []AvailableDateResponse {
	availableDates := make([]AvailableDateResponse, 0)

	for key := range data {
		availableDate := AvailableDateResponse{
			Date: key,
		}

		if len(data[key]) != 0 {
			availableDate.Status = "available"
		} else {
			availableDate.Status = "fully book"
		}

		availableDates = append(availableDates, availableDate)
	}

	sort.Slice(availableDates[:], func(i, j int) bool {
		return availableDates[i].Date < availableDates[j].Date
	})

	return availableDates
}

func (s service) difference(s1 []CheckedItemParams, s2 []CheckedItemParams) []CheckedItemParams {
	mb := make(map[CheckedItemParams]struct{}, len(s2))
	for _, x := range s2 {
		mb[x] = struct{}{}
	}

	var diff []CheckedItemParams
	for _, x := range s1 {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}

func (s *service) GetDetail(bookingID int) (*Detail, error) {
	errorList := []string{}

	if bookingID <= 0 {
		errorList = append(errorList, "bookingID must be above 0")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	bookingDetail, err := s.repo.GetDetail(bookingID)
	if err != nil {
		return nil, err
	}

	ticketPriceWrapper, err := s.repo.GetTicketPriceWrapper(bookingID)
	if err != nil {
		return nil, err
	}

	itemsWrapper, err := s.repo.GetItemWrapper(bookingID)
	if err != nil {
		return nil, err
	}

	totalPriceTicket := ticketPriceWrapper.Price
	totalPrice := totalPriceTicket + bookingDetail.TotalPriceItem

	bookingDetail.TotalPriceTicket = totalPriceTicket
	bookingDetail.TotalPrice = totalPrice

	bookingDetail.Items = make([]ItemDetail, 0)
	for _, item := range itemsWrapper.Items {
		bookingDetail.Items = append(bookingDetail.Items, item)
	}

	return bookingDetail, nil
}

func (s *service) UpdateBookingStatus(bookingID int, newStatus int) error {
	errorList := []string{}

	if bookingID <= 0 {
		errorList = append(errorList, "bookingID must be above 0")
	}

	if newStatus < 0 {
		errorMessage := fmt.Sprintf("there are no status: %d", newStatus)
		errorList = append(errorList, errorMessage)
	}

	if len(errorList) > 0 {
		return errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	err := s.repo.UpdateBookingStatus(bookingID, newStatus)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateBookingStatusByXendit(callback XenditInvoicesCallback) error {
	var errorList []string

	var status int
	switch callback.Status {
	case util.XenditStatusPaid:
		status = util.BookingBerhasil
	case util.XenditStatusExpired:
		status = util.BookingGagal
	default:
		errorList = append(errorList, fmt.Sprintf("callback status %s is unknown", callback.Status))
	}

	if len(errorList) > 0 {
		return errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	err := s.repo.UpdateBookingStatusByXenditID(callback.ID, status)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	errorList := []string{}

	if localID == "" {
		errorList = append(errorList, "localID cannot be empty")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	myBookingsOngoing, err := s.repo.GetMyBookingsOngoing(localID)
	if err != nil {
		return nil, err
	}

	return myBookingsOngoing, nil
}

func (s service) GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, *util.Pagination, error) {
	var errorList []string

	if params.Page == 0 {
		params.Page = util.DefaultPage
	}

	if params.Limit == 0 {
		params.Limit = util.DefaultLimit
	}

	if params.Limit > util.MaxLimit {
		errorList = append(errorList, "limit should be 1 - 100")
	}

	if params.Path == "" {
		errorList = append(errorList, "path is required for pagination")
	}

	if len(errorList) > 0 {
		return nil, nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	myBookingsPrevious, err := s.repo.GetMyBookingsPreviousWithPagination(localID, params)
	if err != nil {
		return nil, nil, err
	}

	pagination := util.GeneratePagination(myBookingsPrevious.TotalCount, params.Limit, params.Page, params.Path)

	return myBookingsPrevious, &pagination, err
}
