package place

import (
	"fmt"
	"math"
	"strings"

	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repo will contain all the function that can be used by repo
type Repo interface {
	GetPlacesListWithPagination(params PlacesListRequest) (*PlacesList, error)
	GetPlaceRatingAndReviewCountByPlaceID(int) (*PlacesRatingAndReviewCount, error)
	GetDetail(int) (*Detail, error)
	GetAverageRatingAndReviews(int) (*AverageRatingAndReviews, error)
	GetListReviewAndRatingWithPagination(params ListReviewRequest) (*ListReview, error)
}

type repo struct {
	db *sqlx.DB
}

// NewRepo used to initialize repo
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetDetail(placeID int) (*Detail, error) {
	var result Detail

	query := `SELECT id, name, image, address, description, open_hour, close_hour, COALESCE (booking_price,0) as booking_price, min_slot_booking, max_slot_booking, min_interval_booking, max_interval_booking, capacity
			  FROM places
			  WHERE id = $1`
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}

func (r *repo) GetAverageRatingAndReviews(placeID int) (*AverageRatingAndReviews, error) {
	var result AverageRatingAndReviews
	result.Reviews = make([]UserReview, 0)

	query := "SELECT COUNT(id) as count_review FROM reviews WHERE place_id = $1"
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	var sumRating int

	query = "SELECT COALESCE(SUM(rating), 0) as sum_rating FROM reviews WHERE place_id = $1"
	err = r.db.Get(&sumRating, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	var averageRating float64
	if result.ReviewCount != 0 {
		averageRating = float64(sumRating) / float64(result.ReviewCount)
	}
	var roundedAverageRating = math.Round(averageRating*100) / 100
	result.AverageRating = roundedAverageRating

	query = "SELECT users.name as user, reviews.rating as rating, reviews.content as content FROM reviews LEFT JOIN users ON reviews.user_id = users.id WHERE reviews.place_id = $1 LIMIT 2"
	err = r.db.Select(&result.Reviews, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}

// GetPlacesListWithPagination will do the query to database for getting list places data
func (r repo) GetPlacesListWithPagination(params PlacesListRequest) (*PlacesList, error) {
	var placeList PlacesList
	placeList.Places = make([]Place, 0)
	placeList.TotalCount = 0

	var whereQuery []string

	if params.Category != "" {
		placeIDQuery := fmt.Sprintf(
			`SELECT place_id
			FROM place_category pc JOIN categories c ON pc.category_id = c.id
			WHERE c.content='%s'`,
			params.Category)
		whereQuery = append(whereQuery, fmt.Sprintf("p.id IN (%s)", placeIDQuery))
	}

	if len(params.Price) != 0 {
		var priceQuery []string
		for _, price := range params.Price {
			switch price {
			case "16000":
				priceQuery = append(priceQuery, "(booking_price < 16000)")
			case "16000-40000":
				priceQuery = append(priceQuery, "(booking_price >= 16000 AND booking_price < 40000)")
			case "40000-100000":
				priceQuery = append(priceQuery, "(booking_price >= 40000 AND booking_price < 100000)")
			case "100000":
				priceQuery = append(priceQuery, "(booking_price >= 100000)")
			}
		}
		whereQuery = append(whereQuery, "("+strings.Join(priceQuery, " OR ")+")")
	}
	if len(params.People) != 0 {
		var peopleQuery []string
		for _, people := range params.People {
			switch people {
			case "1":
				peopleQuery = append(peopleQuery, "(capacity < 2)")
			case "2-4":
				peopleQuery = append(peopleQuery, "(capacity >= 2 AND capacity < 5)")
			case "5-10":
				peopleQuery = append(peopleQuery, "(capacity >= 5 AND capacity < 10)")
			case "10":
				peopleQuery = append(peopleQuery, "(capacity >= 10)")
			}
		}
		whereQuery = append(whereQuery, "("+strings.Join(peopleQuery, " OR ")+")")
	}

	distanceQuery := fmt.Sprintf(
		`6371 * ACOS(SIN(RADIANS(%f)) * SIN(RADIANS(p.lat)) + COS(RADIANS(%f)) * COS(RADIANS(p.lat)) * COS(RADIANS(p.long) - RADIANS(%f)))`,
		params.Latitude,
		params.Latitude,
		params.Longitude)

	reviewCountQuery := fmt.Sprintf("SELECT COUNT(r.rating) FROM reviews r WHERE p.id = r.place_id")

	ratingQuery := fmt.Sprintf("SELECT COALESCE(AVG(r.rating), 0.0) FROM reviews r WHERE r.place_id = p.id")

	query := "SELECT p.id, p.name, p.description, p.address, p.image, "
	query += fmt.Sprintf("(%s) as rating, ", ratingQuery)
	query += fmt.Sprintf("(%s) as review_count, ", reviewCountQuery)
	query += fmt.Sprintf("CAST(%s AS integer) AS distance ", distanceQuery)
	query += "FROM places p "
	countQuery := fmt.Sprintf("SELECT p.*, (%s) as rating FROM places p ", ratingQuery)
	if params.Sort == "popularity" {
		query += "LEFT JOIN bookings b on p.id = b.place_id "
		countQuery += "LEFT JOIN bookings b on p.id = b.place_id "
	}
	if len(whereQuery) != 0 {
		query += fmt.Sprintf("WHERE %s ", strings.Join(whereQuery, " AND "))
		countQuery += fmt.Sprintf("WHERE %s ", strings.Join(whereQuery, " AND "))
	}

	switch params.Sort {
	case "distance":
		query += "ORDER BY distance "
	case "popularity":
		query += "GROUP BY p.id ORDER BY COUNT(p.id) DESC "
		countQuery += "GROUP BY p.id"
	}

	if len(params.Rating) != 0 {
		query = fmt.Sprintf("SELECT * FROM (%s) AS temp WHERE ", query)
		countQuery = fmt.Sprintf("SELECT * FROM (%s) AS temp WHERE ", countQuery)
		for i, r := range params.Rating {
			query += fmt.Sprintf("(rating >= %d AND rating < %d)", r, r+1)
			countQuery += fmt.Sprintf("(rating >= %d AND rating < %d)", r, r+1)
			if i != len(params.Rating)-1 {
				query += " OR "
				countQuery += " OR "
			}
		}
		query += " LIMIT $1 OFFSET $2"
	} else {
		query += " LIMIT $1 OFFSET $2"
	}
	err := r.db.Select(&placeList.Places, query, params.Limit, (params.Page-1)*params.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			placeList.Places = make([]Place, 0)
			placeList.TotalCount = 0
			return &placeList, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	err = r.db.Get(&placeList.TotalCount, fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS temp", countQuery))
	if err != nil {
		if err == sql.ErrNoRows {
			placeList.Places = make([]Place, 0)
			placeList.TotalCount = 0
			return &placeList, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &placeList, nil
}

// GetPlaceRatingAndReviewCountByPlaceID will do the query to database for getting review and review count data
func (r repo) GetPlaceRatingAndReviewCountByPlaceID(placeID int) (*PlacesRatingAndReviewCount, error) {
	var result PlacesRatingAndReviewCount

	query := "SELECT COUNT(rating) as review_count, COALESCE (AVG(rating), 0.0) as rating FROM reviews WHERE place_id = $1"
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}

func (r repo) GetListReviewAndRatingWithPagination(params ListReviewRequest) (*ListReview, error) {
	var listReview ListReview
	listReview.Reviews = make([]Review, 0)
	listReview.TotalCount = 0
	
	mainQuery := `
	SELECT r.id, u.name, r.content, r.rating, r.created_at
	FROM reviews r, users u
	WHERE r.place_id = $1 AND u.id = r.user_id `
	branchQuery := ``

	if params.Latest && params.Rating {
		branchQuery = `ORDER BY r.created_at DESC, r.rating DESC`
	} else if params.Latest && !params.Rating {
		branchQuery = `ORDER BY r.created_at DESC`
	} else if !params.Latest && params.Rating {
		branchQuery = `ORDER BY r.rating DESC`
	}

	branchQuery += ` LIMIT $2 OFFSET $3`

	err := r.db.Select(&listReview.Reviews, mainQuery+branchQuery, params.PlaceID, params.Limit, (params.Page-1)*params.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			listReview.Reviews = make([]Review, 0)
			listReview.TotalCount = 0
			return &listReview, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	mainQuery = `SELECT COUNT(r.id) FROM reviews r, users u WHERE r.place_id = $1 AND u.id = r.user_id`
	err = r.db.Get(&listReview.TotalCount, mainQuery, params.PlaceID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &listReview, nil
}
