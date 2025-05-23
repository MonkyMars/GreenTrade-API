package reviews

import (
	"fmt"
	"greenvue/internal/db"
	"greenvue/lib"
	"greenvue/lib/errors"

	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

const viewName string = "review_with_username"

func GetReviews(c *fiber.Ctx) error {
	client := db.GetGlobalClient()
	if client == nil {
		return errors.InternalServerError("Failed to create client")
	}

	selectedSeller := c.Params("seller_id")

	if selectedSeller == "" {
		return errors.BadRequest("Seller ID is required")
	}

	limit := c.Query("limit", "50")
	query := fmt.Sprintf("select=*&limit=%s&seller_id=eq.%s", limit, selectedSeller)

	// Use standardized GET operation
	data, err := client.GET(viewName, query)
	if err != nil {
		return errors.DatabaseError("Failed to fetch reviews: " + err.Error())
	}

	if len(data) == 0 || string(data) == "[]" {
		return errors.SuccessResponse(c, []lib.FetchedReview{})
	}

	var reviews []lib.FetchedReview
	if err := json.Unmarshal(data, &reviews); err != nil {
		return errors.InternalServerError("Failed to parse reviews data: " + err.Error())
	}

	if reviews == nil {
		reviews = []lib.FetchedReview{}
	}

	return errors.SuccessResponse(c, reviews)
}
