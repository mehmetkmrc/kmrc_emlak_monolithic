package user

import (
	
	"encoding/json"
	"kmrc_emlak_mono/database"
	"kmrc_emlak_mono/dto"
	"kmrc_emlak_mono/models"
	"kmrc_emlak_mono/response"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository struct {
	dbPool   *pgxpool.Pool
	validate *validator.Validate
}

/*
	UpdateUser — Base User Info Update
	name, surname, email, phone, about
*/
func UpdateUser(c fiber.Ctx) error {

	// -------------------------
	// 1) REQUEST BODY PARSE
	// -------------------------
	req := new(dto.UserUpdateRequest)

	if err := json.Unmarshal(c.Body(), req); err != nil {
		return response.Error_Response(c,
			"invalid JSON body",
			err,
			nil,
			fiber.StatusBadRequest,
		)
	}

	// -------------------------
	// 2) VALIDATE UUID
	// -------------------------
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return response.Error_Response(c,
			"invalid user_id format",
			err,
			nil,
			fiber.StatusBadRequest,
		)
	}

	// -------------------------
	// 3) Convert DTO → Model
	// -------------------------
	updateModel := &models.User{
		UserID:  userUUID.String(),
		Name:    req.Name,
		Surname: req.Surname,
		Email:   req.Email,
		Phone:   req.Phone,
		About:   req.AboutText,
	}

	// -------------------------
	// 4) Database Query
	// -------------------------
	query := `
        UPDATE users
        SET 
            first_name = $1,
            last_name = $2,
            email      = $3,
            phone      = $4,
            about_text = $5,
            updated_at = NOW()
        WHERE user_id = $6
        RETURNING user_id, first_name, last_name, email, phone, about_text
    `

	row := database.DBPool.QueryRow(
		c.Context(),
		query,
		updateModel.Name,
		updateModel.Surname,
		updateModel.Email,
		updateModel.Phone,
		updateModel.About,
		updateModel.UserID,
	)

	updatedUser := &models.User{}

	err = row.Scan(
		&updatedUser.UserID,
		&updatedUser.Name,
		&updatedUser.Surname,
		&updatedUser.Email,
		&updatedUser.Phone,
		&updatedUser.About,
	)

	if err != nil {
		return response.Error_Response(c,
			"error while updating user",
			err,
			nil,
			fiber.StatusInternalServerError,
		)
	}

	// -------------------------
	// 5) LOG & RETURN RESPONSE
	// -------------------------
	zap.S().Info("User updated successfully: ", updatedUser)

	return response.Success_Response(
		c,
		updatedUser,
		"User info updated successfully!",
		fiber.StatusOK,
	)
}
