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
		AboutText:   req.AboutText,
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
		updateModel.AboutText,
		updateModel.UserID,
	)

	updatedUser := &models.User{}

	err = row.Scan(
		&updatedUser.UserID,
		&updatedUser.Name,
		&updatedUser.Surname,
		&updatedUser.Email,
		&updatedUser.Phone,
		&updatedUser.AboutText,
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
	zap.S().Info("Kullanıcı başarıyla Güncellendi: ", updatedUser)

	return response.Success_Response(
		c,
		updatedUser,
		"User info updated successfully!",
		fiber.StatusOK,
	)
}

// Get social links by user_id
// Get social links by user_id (UpdateUser formatında)
func GetSocialLinksByUserID(c fiber.Ctx) error {
 
	req := new(dto.UpdateSocialLinksRequest)

	if err := json.Unmarshal(c.Body(), req); err != nil{
		return response.Error_Response(c,
		"invalid JSON body",
		err,
		nil,
		fiber.StatusBadRequest,
		)
	}
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil{
		return response.Error_Response(c, 
		"invalid user_id format",
		err,
		nil,
		fiber.StatusBadRequest,
		)
	}

	socialLinks := &models.UserSocialLinks{
		UserID: userUUID.String(),
		Facebook: req.Facebook,
		Tiktok: req.Tiktok,
		Instagram: req.Instagram,
		Twitter: req.Twitter,
		Youtube: req.Youtube,
		Linkedin: req.Linkedin,
	}

	// -------------------------
	// 4) Database Query
	// -------------------------
	query := `
		SELECT 
			id,
			user_id,
			facebook,
			tiktok,
			instagram,
			twitter,
			youtube,
			linkedin,
			updated_at
		FROM user_social_links
		WHERE user_id = $1
	`

	row := database.DBPool.QueryRow(
		c.Context(),
		query,
		socialLinks.Facebook,
		socialLinks.Tiktok,
		socialLinks.Instagram,
		socialLinks.Twitter,
		socialLinks.Youtube,
		socialLinks.Linkedin,
	)


	s := &models.UserSocialLinks{}

	err = row.Scan(
		&s.ID,
		&s.UserID,
		&s.Facebook,
		&s.Tiktok,
		&s.Instagram,
		&s.Twitter,
		&s.Youtube,
		&s.Linkedin,
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
	zap.S().Info("Kullanıcı başarıyla Güncellendi: ", s)

	return response.Success_Response(
		c,
		s,
		"User info updated successfully!",
		fiber.StatusOK,
	)
}



// Upsert social links (insert if not exists, else update)
// Upsert social links (UpdateUser formatında RETURNING ile)
func UpsertSocialLinks(c fiber.Ctx) error {

	req := new(dto.UpdateSocialLinksRequest)
	if err := json.Unmarshal(c.Body(), req); err != nil{
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
	updateModel := &models.UserSocialLinks{
		ID: int(userUUID.ID()),
		UserID: userUUID.String(),
		Facebook: req.Facebook,
		Tiktok: req.Tiktok,
		Instagram: req.Instagram,
		Twitter: req.Twitter,
		Youtube: req.Youtube,
		Linkedin: req.Linkedin,
		UpdatedAt: req.UpdatedAt,
	}

	// -------------------------
	// 4) Database Query
	// -------------------------
	query := `
		INSERT INTO user_social_links 
		(user_id, facebook, tiktok, instagram, twitter, youtube, linkedin, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			facebook  = EXCLUDED.facebook,
			tiktok    = EXCLUDED.tiktok,
			instagram = EXCLUDED.instagram,
			twitter   = EXCLUDED.twitter,
			youtube   = EXCLUDED.youtube,
			linkedin  = EXCLUDED.linkedin,
			updated_at = NOW()
		RETURNING 
			id,
			user_id,
			facebook,
			tiktok,
			instagram,
			twitter,
			youtube,
			linkedin,
			updated_at
	`

	row := database.DBPool.QueryRow(
		c.Context(),
		query,
		updateModel.UserID,
		updateModel.Facebook,
		updateModel.Tiktok,
		updateModel.Instagram,
		updateModel.Twitter,
		updateModel.Youtube,
		updateModel.Linkedin,
	)


	s := &models.UserSocialLinks{}

	err = row.Scan(
		&s.ID,
		&s.UserID,
		&s.Facebook,
		&s.Tiktok,
		&s.Instagram,
		&s.Twitter,
		&s.Youtube,
		&s.Linkedin,
		&s.UpdatedAt,
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
	zap.S().Info("Kullanıcı başarıyla Güncellendi: ", s)

	return response.Success_Response(
		c,
		s,
		"User info updated successfully!",
		fiber.StatusOK,
	)
}
