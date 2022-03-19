package businessadminauth

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRepo_CheckRequiredFields(t *testing.T) {
	request := &RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}

	var expectedReturn []string

	// Mock DB
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	retrievedReturn := repoMock.CheckRequiredFields(*request, expectedReturn)
	assert.Equal(t, retrievedReturn, expectedReturn)
}

func TestRepo_CheckUserFields(t *testing.T) {
	request := &RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	rows := mock.
		NewRows([]string{"phone_number"})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1")).
		WithArgs("089782828888").
		WillReturnRows(rows)

	rows = mock.
		NewRows([]string{"email"})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT email FROM users WHERE email=$1 LIMIT 1")).
		WithArgs("sebuahemail@gmail.com").
		WillReturnRows(rows)

	err = repoMock.CheckUserFields(*request)
	assert.NoError(t, err)
}

func TestRepo_CheckBusinessAdminFields(t *testing.T) {
	request := &RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	rows := mock.
		NewRows([]string{"name"})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT bank_account_number FROM business_owners WHERE bank_account_number=$1 LIMIT 1")).
		WithArgs("008-112492374950").
		WillReturnRows(rows)

	err = repoMock.CheckBusinessAdminFields(*request)
	assert.NoError(t, err)
}

func TestRepo_CheckPlaceFields(t *testing.T) {
	request := &RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	rows := mock.
		NewRows([]string{"name"})
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM places WHERE name=$1 LIMIT 1")).
		WithArgs("Kopi Kenangan").
		WillReturnRows(rows)

	err = repoMock.CheckPlaceFields(*request)
	assert.NoError(t, err)
}

func TestRepo_CheckIfPhoneNumberIsUnique(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	t.Run("phone number is not unique in database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"phone_number"}).
			AddRow("081234567890")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1")).
			WithArgs("081234567890").
			WillReturnRows(rows)

		unique, err := repoMock.CheckIfPhoneNumberIsUnique("081234567890")
		assert.NoError(t, err)
		assert.False(t, unique)
	})

	t.Run("phone number is unique in database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"phone_number"})
		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1")).
			WithArgs("081234567891").
			WillReturnRows(rows)

		unique, err := repoMock.CheckIfPhoneNumberIsUnique("081234567891")
		assert.NoError(t, err)
		assert.True(t, unique)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1")).
			WithArgs("081234567892").
			WillReturnError(sql.ErrTxDone)

		unique, err := repoMock.CheckIfPhoneNumberIsUnique("081234567892")
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.False(t, unique)
	})
}

func TestRepo_CheckIfBankAccountIsUnique(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	t.Run("bank account is not unique in database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"bank_account_number"}).
			AddRow("008-1234567890")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT bank_account_number FROM business_owners WHERE bank_account_number=$1 LIMIT 1")).
			WithArgs("008-1234567890").
			WillReturnRows(rows)

		unique, err := repoMock.CheckIfBankAccountIsUnique("008-1234567890")
		assert.NoError(t, err)
		assert.False(t, unique)
	})

	t.Run("bank account is unique in database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"bank_account_number"})
		mock.ExpectQuery(regexp.QuoteMeta("SELECT bank_account_number FROM business_owners WHERE bank_account_number=$1 LIMIT 1")).
			WithArgs("008-1234567812").
			WillReturnRows(rows)

		unique, err := repoMock.CheckIfBankAccountIsUnique("008-1234567812")
		assert.NoError(t, err)
		assert.True(t, unique)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT bank_account_number FROM business_owners WHERE bank_account_number=$1 LIMIT 1")).
			WithArgs("008-1234567834").
			WillReturnError(sql.ErrTxDone)

		unique, err := repoMock.CheckIfBankAccountIsUnique("008-1234567834")
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.False(t, unique)
	})
}

func TestRepo_CheckIfEmailIsUnique(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	t.Run("email is not unique in database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"email"}).
			AddRow("sebuahemail1@gmail.com")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT email FROM users WHERE email=$1 LIMIT 1")).
			WithArgs("sebuahemail1@gmail.com").
			WillReturnRows(rows)

		unique, err := repoMock.CheckIfEmailIsUnique("sebuahemail1@gmail.com")
		assert.NoError(t, err)
		assert.False(t, unique)
	})

	t.Run("email is unique in database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"email"})
		mock.ExpectQuery(regexp.QuoteMeta("SELECT email FROM users WHERE email=$1 LIMIT 1")).
			WithArgs("sebuahemail2@gmail.com").
			WillReturnRows(rows)

		unique, err := repoMock.CheckIfEmailIsUnique("sebuahemail2@gmail.com")
		assert.NoError(t, err)
		assert.True(t, unique)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT email FROM users WHERE email=$1 LIMIT 1")).
			WithArgs("sebuahemail3@gmail.com").
			WillReturnError(sql.ErrTxDone)

		unique, err := repoMock.CheckIfBankAccountIsUnique("sebuahemail3@gmail.com")
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.False(t, unique)
	})
}

func TestRepo_CheckIfPlaceNameIsUnique(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	t.Run("place name is not unique in database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"name"}).
			AddRow("Kopi Kenangan Rawamangun")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM places WHERE name=$1 LIMIT 1")).
			WithArgs("Kopi Kenangan Rawamangun").
			WillReturnRows(rows)

		unique, err := repoMock.CheckIfPlaceNameIsUnique("Kopi Kenangan Rawamangun")
		assert.NoError(t, err)
		assert.False(t, unique)
	})

	t.Run("place name is unique in database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"name"})
		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM places WHERE name=$1 LIMIT 1")).
			WithArgs("Kopi Kenangan Pancoran").
			WillReturnRows(rows)

		unique, err := repoMock.CheckIfPlaceNameIsUnique("Kopi Kenangan Pancoran")
		assert.NoError(t, err)
		assert.True(t, unique)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM places WHERE name=$1 LIMIT 1")).
			WithArgs("Kopi Kenangan Pasar Minggu").
			WillReturnError(sql.ErrTxDone)

		unique, err := repoMock.CheckIfPlaceNameIsUnique("Kopi Kenangan Pasar Minggu")
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.False(t, unique)
	})
}

func TestRepo_VerifyHour(t *testing.T) {
	var hour, hourName = "23:59", "openHour"

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	verified, err := repoMock.VerifyHour(hour, hourName)
	assert.NoError(t, err)
	assert.True(t, verified)
}

func TestRepo_CompareOpenAndCloseHour(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	var openHour, closeHour = "08:00", "21:00"
	compared, err := repoMock.CompareOpenAndCloseHour(openHour, closeHour)
	assert.NoError(t, err)
	assert.True(t, compared)

	openHour, closeHour = "09:00", "08:00"
	compared, err = repoMock.CompareOpenAndCloseHour(openHour, closeHour)
	assert.NoError(t, err)
	assert.False(t, compared)
}

func TestRepo_GeneratePassword(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	password := repoMock.GeneratePassword()
	passwordLength := len(password)
	assert.Equal(t, 8, passwordLength)
}

func TestRepo_CreateUser(t *testing.T) {
	request := &RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}
	password := "12345678"
	status := 1

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	_ = mock.
		NewRows([]string{"phone_number", "name", "email", "password", "status"})
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (phone_number, name, email, password, status) VALUES ($1, $2, $3, $4, $5)")).
		WithArgs(request.AdminPhoneNumber, request.AdminName, request.AdminEmail, password, status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repoMock.CreateUser(request.AdminPhoneNumber, request.AdminName, request.AdminEmail, password, status)
	assert.NoError(t, err)
}

func TestRepo_RetrieveUserID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	rows := mock.NewRows([]string{"id", "phone_number"}).
		AddRow(1, "081234567890")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE phone_number=$1 LIMIT 1")).
		WithArgs("081234567890").
		WillReturnRows(rows)

	userID, err := repoMock.RetrieveUserID("081234567890")
	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
}

func TestRepo_CreateBusinessAdmin(t *testing.T) {
	userID := 1
	request := &RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}
	balance := 0.0

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	_ = mock.
		NewRows([]string{"balance", "bank_account_number", "bank_account_name", "user_id"})
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO business_owners (balance, bank_account_number, bank_account_name, user_id) VALUES ($1, $2, $3, $4)")).
		WithArgs(balance, request.AdminBankAccount, request.AdminBankAccountName, userID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repoMock.CreateBusinessAdmin(userID, request.AdminBankAccount, request.AdminBankAccountName, float32(balance))
	assert.NoError(t, err)
}

func TestRepo_CreatePlace(t *testing.T) {
	userID := 1
	request := &RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	_ = mock.
		NewRows([]string{"name", "address", "capacity", "description", "user_id", "interval", "open_hour", "close_hour", "image",
			"min_interval_booking", "max_interval_booking", "min_slot_booking", "max_slot_booking", "lat", "long"})
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO places (
		name, address, capacity, description, user_id, interval, open_hour, close_hour, image,
		min_interval_booking, max_interval_booking, min_slot_booking, max_slot_booking, lat, long) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`)).
		WithArgs(
			request.PlaceName,
			request.PlaceAddress,
			request.PlaceCapacity,
			request.PlaceDescription,
			userID,
			request.PlaceInterval,
			request.PlaceOpenHour,
			request.PlaceCloseHour,
			request.PlaceImage,
			request.PlaceMinIntervalBooking,
			request.PlaceMaxIntervalBooking,
			request.PlaceMinSlotBooking,
			request.PlaceMaxSlotBooking,
			request.PlaceLat,
			request.PlaceLong,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repoMock.CreatePlace(request.PlaceName, request.PlaceAddress, request.PlaceCapacity,
		request.PlaceDescription, userID, request.PlaceInterval, request.PlaceOpenHour, request.PlaceCloseHour,
		request.PlaceImage, request.PlaceMinIntervalBooking, request.PlaceMaxIntervalBooking, request.PlaceMinSlotBooking,
		request.PlaceMaxSlotBooking, request.PlaceLat, request.PlaceLong)
	assert.NoError(t, err)
}

func Test_repo_GetBusinessAdminByEmail(t *testing.T) {
	expected := &BusinessAdmin{
		ID:                1,
		Name:              "Teofanus Gary",
		PhoneNumber:       "081223906674",
		Email:             "test@gmail.com",
		Password:          "testpassword",
		Status:            1,
		Balance:           1000,
		BankAccountNumber: "12321asdfasdf",
		BankAccountName:   "BCA",
	}
	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)

	t.Run("business admin exist on database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"id", "phone_number", "name", "status", "email", "password"}).
			AddRow(expected.ID, expected.PhoneNumber, expected.Name, expected.Status, expected.Email, expected.Password)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).
			WithArgs("test@gmail.com").
			WillReturnRows(rows)

		rows = mock.
			NewRows([]string{"id", "balance", "bank_account_number", "user_id", "bank_account_name"}).
			AddRow(1, expected.Balance, expected.BankAccountNumber, expected.ID, expected.BankAccountName)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM business_owners WHERE user_id = $1")).
			WithArgs(expected.ID).
			WillReturnRows(rows)

		actual, err := repoMock.GetBusinessAdminByEmail("test@gmail.com")
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})

	t.Run("error getting user data", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).
			WithArgs("test@gmail.com").
			WillReturnError(ErrInternalServerError)

		actual, err := repoMock.GetBusinessAdminByEmail("test@gmail.com")
		assert.Nil(t, actual)
		assert.True(t, errors.Is(err, ErrInternalServerError))
	})

	t.Run("error getting business admin data", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"id", "phone_number", "name", "status"}).
			AddRow(1, "081223901234", "Bambang", 1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).
			WithArgs("test@gmail.com").
			WillReturnRows(rows)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM business_owners WHERE user_id = $1")).
			WithArgs(1).
			WillReturnError(ErrInternalServerError)

		actual, err := repoMock.GetBusinessAdminByEmail("test@gmail.com")
		assert.True(t, errors.Is(err, ErrInternalServerError))
		assert.Nil(t, actual)
	})
}
