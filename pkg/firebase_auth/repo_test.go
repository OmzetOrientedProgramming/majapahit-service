package firebaseauth

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestRepo_VerifyOTP(t *testing.T) {
	godotenv.Load("../../.env")
	firebaseRepo := NewRepo(
		os.Getenv("IDENTITY_TOOLKIT_URL"),
		os.Getenv("SECURE_TOKEN_URL"),
		os.Getenv("FIREBASE_API_KEY"),
	)

	t.Run("success", func(t *testing.T) {
		params := VerifyOTPParams{
			Code:        "111111",
			SessionInfo: "AJOnW4S-gu6XV_3bljabYFTgO6A_RlKovolps9ULQ2Rc9JBKfDZQBaj18MJfFOuH2S109P6SFav5D6jXjLsXuQiYg05ZqIKWa0CQLmXigKmU2ykJZWt5kAhdNIXscvX6s0t_DOut428zPcgg2zecP0QKN9gjr8T0Kw",
		}

		res, err := firebaseRepo.VerifyOTP(params)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("session expired", func(t *testing.T) {
		params := VerifyOTPParams{
			Code:        "111111",
			SessionInfo: "AJOnW4R97JHnsLOXnhSaqpDa2rC2k7vFC9D7ly5-u1B33RqoYKZ44BZUSe7yZszk98z6mWsoGraXcL18Imh2qDlZA8u6uVjyk244hhOqkhOENnsgtlieLUZeB0yCgejrgq5z8hQgrxkK5HKVLzMFLFGCjIS1l1KSm4aW4ALBZOM-t7saMLIHLtQeH9O-Y2-YDGl7g1qYQ4FAvQmiwbNmWjMCVmxmObIdPtmLHItXBRxAsFNvmbDxk9Q",
		}

		res, err := firebaseRepo.VerifyOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("invalid code", func(t *testing.T) {
		params := VerifyOTPParams{
			Code:        "1111111",
			SessionInfo: "AJOnW4S-gu6XV_3bljabYFTgO6A_RlKovolps9ULQ2Rc9JBKfDZQBaj18MJfFOuH2S109P6SFav5D6jXjLsXuQiYg05ZqIKWa0CQLmXigKmU2ykJZWt5kAhdNIXscvX6s0t_DOut428zPcgg2zecP0QKN9gjr8T0Kw",
		}

		res, err := firebaseRepo.VerifyOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("invalid session info", func(t *testing.T) {
		params := VerifyOTPParams{
			Code:        "1111111",
			SessionInfo: "",
		}

		res, err := firebaseRepo.VerifyOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("failed in http client post", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"localhost:8001",
			os.Getenv("test token"),
			os.Getenv("test token"),
		)

		params := VerifyOTPParams{
			Code:        "111111",
			SessionInfo: "AJOnW4S-gu6XV_3bljabYFTgO6A_RlKovolps9ULQ2Rc9JBKfDZQBaj18MJfFOuH2S109P6SFav5D6jXjLsXuQiYg05ZqIKWa0CQLmXigKmU2ykJZWt5kAhdNIXscvX6s0t_DOut428zPcgg2zecP0QKN9gjr8T0Kw",
		}

		res, err := firebaseRepo.VerifyOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("failed invalid format", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/invalid-format",
			os.Getenv("test token"),
			os.Getenv("test token"),
		)

		params := VerifyOTPParams{
			Code:        "111111",
			SessionInfo: "AJOnW4S-gu6XV_3bljabYFTgO6A_RlKovolps9ULQ2Rc9JBKfDZQBaj18MJfFOuH2S109P6SFav5D6jXjLsXuQiYg05ZqIKWa0CQLmXigKmU2ykJZWt5kAhdNIXscvX6s0t_DOut428zPcgg2zecP0QKN9gjr8T0Kw",
		}

		res, err := firebaseRepo.VerifyOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("unknown error", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/unknown-error",
			os.Getenv("test token"),
			os.Getenv("test token"),
		)

		params := VerifyOTPParams{
			Code:        "111111",
			SessionInfo: "AJOnW4S-gu6XV_3bljabYFTgO6A_RlKovolps9ULQ2Rc9JBKfDZQBaj18MJfFOuH2S109P6SFav5D6jXjLsXuQiYg05ZqIKWa0CQLmXigKmU2ykJZWt5kAhdNIXscvX6s0t_DOut428zPcgg2zecP0QKN9gjr8T0Kw",
		}

		res, err := firebaseRepo.VerifyOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("failed unmarshal error response from firebase", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/failed-unmarshal",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		params := VerifyOTPParams{
			Code:        "111111",
			SessionInfo: "AJOnW4S-gu6XV_3bljabYFTgO6A_RlKovolps9ULQ2Rc9JBKfDZQBaj18MJfFOuH2S109P6SFav5D6jXjLsXuQiYg05ZqIKWa0CQLmXigKmU2ykJZWt5kAhdNIXscvX6s0t_DOut428zPcgg2zecP0QKN9gjr8T0Kw",
		}

		res, err := firebaseRepo.VerifyOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

}

func TestRepo_SendOTP(t *testing.T) {
	godotenv.Load("../../.env")
	firebaseRepo := NewRepo(
		os.Getenv("IDENTITY_TOOLKIT_URL"),
		os.Getenv("SECURE_TOKEN_URL"),
		os.Getenv("FIREBASE_API_KEY"),
	)

	t.Run("success", func(t *testing.T) {
		params := SendOTPParams{
			PhoneNumber:    "+621111111111",
			RecaptchaToken: "testCapthcaToken",
		}

		res, err := firebaseRepo.SendOTP(params)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("failed in http client post", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"localhost:8001",
			os.Getenv("test token"),
			os.Getenv("test token"),
		)

		params := SendOTPParams{
			PhoneNumber:    "+621111111111",
			RecaptchaToken: "testCapthcaToken",
		}

		res, err := firebaseRepo.SendOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("failed captcha expired", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/session-expired",
			os.Getenv("test token"),
			os.Getenv("test token"),
		)

		params := SendOTPParams{
			PhoneNumber:    "",
			RecaptchaToken: "03AGdBq2592zaxWm6FTmba1wSFEzTTiY_LMIaVuUni5UTNAmXwXTyNxQ9k5MyLlgBbm0VzltKL2KrlzDcUYHmmmk2Sl5p0JTzyuSqphQ_tdyxK0a5gcLauknxt1XRMVTv2Q3EzsDJ8OuFEKx-yNizWpX_SN59czBhQGZEmF_m2fkt8ne4zibbOtvK9Iazb-jayM4nDaLLUUJvm6SaWXWKzac98ybHmelbqKmd3Lts6UrHBfDwWZZ6ilx1ov5BtmZye9i5UP1fjuxoCOXOE6jAoEUgQpdGoEOb0-rvN1fnKd1ZfwwNfQONW2I7Q4zFfQnm1_zsUUhPTvQ8AOZCgNTguX1IRA_oTkz63v6yOEQyD_RqFBvW3G7l6-exxCCNIFtSLAm5D4H3UFBMgdB2HjPub41PinX-eYsFvewF4JlX-L5vmb0CgDDEiTs7JdvPs5x85rPbuJgVBQ5fq9SKfEUFhPi0pGSC7zor4LA",
		}

		res, err := firebaseRepo.SendOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("failed unmarshal error response from firebase", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/failed-unmarshal",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		params := SendOTPParams{
			PhoneNumber:    "+6211111111111",
			RecaptchaToken: "GSC7zor4LA",
		}

		res, err := firebaseRepo.SendOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("failed unknown error", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/unknown-error",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		params := SendOTPParams{
			PhoneNumber:    "+6211111111111",
			RecaptchaToken: "GSC7zor4LA",
		}

		res, err := firebaseRepo.SendOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("invalid sessionInfo format", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/invalid-session-info-format",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		params := SendOTPParams{
			PhoneNumber:    "+6211111111111",
			RecaptchaToken: "GSC7zor4LA",
		}

		res, err := firebaseRepo.SendOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("session info empty", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/session-info-empty",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		params := SendOTPParams{
			PhoneNumber:    "+6211111111111",
			RecaptchaToken: "GSC7zor4LA",
		}

		res, err := firebaseRepo.SendOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("too many attempt", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/too-many-attempt",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		params := SendOTPParams{
			PhoneNumber:    "+6211111111111",
			RecaptchaToken: "GSC7zor4LA",
		}

		res, err := firebaseRepo.SendOTP(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, res)
	})
}

func TestRepo_GetUserData(t *testing.T) {
	godotenv.Load("../../.env")
	firebaseRepo := NewRepo(
		os.Getenv("IDENTITY_TOOLKIT_URL"),
		os.Getenv("SECURE_TOKEN_URL"),
		os.Getenv("FIREBASE_API_KEY"),
	)

	t.Run("success", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/success/user-data-success",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		res, err := firebaseRepo.GetUserDataFromToken("test token")
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("failed in http client post", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"localhost:8001",
			os.Getenv("test token"),
			os.Getenv("test token"),
		)

		res, err := firebaseRepo.GetUserDataFromToken("test token")
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("failed unmarshal", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/failed-unrmarshal",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		res, err := firebaseRepo.GetUserDataFromToken("test token")
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("failed unmarshal error", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/failed-unrmarshal-error",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		res, err := firebaseRepo.GetUserDataFromToken("test token")
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("invalid token", func(t *testing.T) {
		res, err := firebaseRepo.GetUserDataFromToken("eyJhbGciOiJSUzI1NiIsImtpZCI6ImYxZDU2YTI1MWU0ZGRhM2Y0NWM0MWZkNWQ0ZGEwMWQyYjlkNzJlMGQiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vYm9yZWFsLWZvcmVzdC0zNDQyMDQiLCJhdWQiOiJib3JlYWwtZm9yZXN0LTM0NDIwNCIsImF1dGhfdGltZSI6MTY0NzUyMDI1OCwidXNlcl9pZCI6ImdwQkJFZWo2UHJVQWE1dXNtQkdQRFdacFM5ajEiLCJzdWIiOiJncEJCRWVqNlByVUFhNXVzbUJHUERXWnBTOWoxIiwiaWF0IjoxNjQ3NTIwMjU4LCJleHAiOjE2NDc1MjM4NTgsInBob25lX251bWJlciI6Iis2MjExMTExMTExMTEiLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7InBob25lIjpbIis2MjExMTExMTExMTEiXX0sInNpZ25faW5fcHJvdmlkZXIiOiJwaG9uZSJ9fQ.LxSh-6hP1vGawS8hPHM3p0c795w_uQB2WxySFsDmjLY_OIOneOmHFZKLvwxLaSzCCbYB7SB1uZer736tCwDkFoVxtAwrPZTzq07sUCWY6k3aLGzZ3xud6Iib0MFDz0TAcNo0vZVdZ3nOO2Cc9fM1EOHkxneqZeGusl5f2XGbUOFb-qHoL8mBZhEin01gHD7oc251yyT7m0W7-LLndEJi59AJMxGoWggK_-iIdJvy-4My6mfkcuUxC-e1N4RNvNPbOM1CenG-wtgO5_CEtyzkeqki4QYihZiPKJ73Q16kxyzYdm1VIxKl5vBdKMhkk3huQSGu_pXlRDpNTuMMUUtYnQ")
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, res)
	})

	t.Run("unknown error", func(t *testing.T) {
		firebaseRepo := NewRepo(
			"https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase/error/unknwon-error",
			os.Getenv("SECURE_TOKEN_URL"),
			os.Getenv("FIREBASE_API_KEY"),
		)

		res, err := firebaseRepo.GetUserDataFromToken("test token")
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, res)
	})
}
