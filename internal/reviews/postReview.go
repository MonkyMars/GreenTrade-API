package reviews

import (
	"greenvue/internal/db"
	"greenvue/lib"
	"greenvue/lib/errors"

	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PostReview(c *fiber.Ctx) error {
	client := db.GetGlobalClient()
	if client == nil {
		return errors.InternalServerError("Failed to create client")
	}

	var review lib.Review
	if err := c.BodyParser(&review); err != nil {
		return errors.BadRequest("Invalid request body: " + err.Error())
	}

	// Validate required fields
	if review.UserID == uuid.Nil || review.SellerID == uuid.Nil {
		return errors.BadRequest("UserID, SellerID, and ListingID are required")
	}

	// Validate rating
	if review.Rating < 1 || review.Rating > 5 {
		return errors.BadRequest("Rating must be between 1 and 5")
	}

	// Use standardized POST operation
	data, err := client.POST("reviews", review)
	if err != nil {
		return errors.DatabaseError("Failed to post review: " + err.Error())
	}

	if len(data) == 0 || string(data) == "[]" {
		return errors.InternalServerError("Failed to create review")
	}

	// Parse the response
	var createdReview lib.Review
	if err := json.Unmarshal(data, &createdReview); err != nil {
		// If the response is an array, try parsing it as an array
		var reviewArray []lib.Review
		if err := json.Unmarshal(data, &reviewArray); err != nil {
			return errors.InternalServerError("Failed to parse review response: " + err.Error())
		}

		if len(reviewArray) == 0 {
			return errors.InternalServerError("Empty review response")
		}

		createdReview = reviewArray[0]
	}

	return errors.SuccessResponse(c, createdReview)
}
