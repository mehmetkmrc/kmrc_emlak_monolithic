package middleware

import (
	"fmt"
	"kmrc_emlak_mono/auth"
	"kmrc_emlak_mono/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func PropertyMiddleware(c fiber.Ctx) error {
	payload, ok := c.Locals(auth.AuthPayload).(*auth.Payload)
	if !ok {
		fmt.Println("payload boş döndü...")
		fmt.Println(c.Locals(auth.AuthPayload))
		return response.Error_Response(c, "payload not found in context", nil, nil, fiber.StatusInternalServerError)
	}

	userIDString := payload.ID // string olarak UserID
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return response.Error_Response(c, "invalid user ID format", err, nil, fiber.StatusBadRequest)
	}

	propertyID := uuid.New()

	// **Context'e PropertyID ve UserID'yi kaydet**
	c.Locals("propertyID", propertyID)
	c.Locals("userID", userID)

	return c.Next()
}