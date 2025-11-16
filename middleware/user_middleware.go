package middleware

import (
	"fmt"
	"kmrc_emlak_mono/auth"
	"kmrc_emlak_mono/response"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func UserMiddleware(c fiber.Ctx) error {
	payload, ok := c.Locals(auth.AuthPayload).(*auth.Payload)
	if !ok{
		fmt.Println("payload boş döndü...")
		fmt.Println(c.Locals(auth.AuthPayload))
		return response.Error_Response(c, "payload not found in context", nil, nil, fiber.StatusInternalServerError)
	}
	userIDString := payload.ID
	userID, err := uuid.Parse(userIDString)
	if err != nil{
		return response.Error_Response(c, "invalid user ID format", err, nil, fiber.StatusBadRequest)
	}
	c.Locals("userID", userID)
	return c.Next()
}