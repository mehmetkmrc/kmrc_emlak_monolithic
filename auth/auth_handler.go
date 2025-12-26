package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"kmrc_emlak_mono/database"
	"kmrc_emlak_mono/dto"
	"kmrc_emlak_mono/models"
	"kmrc_emlak_mono/response"
	"regexp"
	"strings"
	"sync"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	AuthHeader    = "Authorization"
	AccessToken   = "access_token"
	AccessPublic  = "access_public"
	RefreshToken  = "refresh_token"
	RefreshPublic = "refresh_public"
	UserDetail    = "UserDetail"
	AuthType      = "Bearer"
	AuthPayload   = "Payload"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")

	customValidators = map[string]func(validator.FieldLevel) bool{
		"email": func(fl validator.FieldLevel) bool {
			// Basic email regex
			email := fl.Field().String()
			regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
			return regex.MatchString(email)
		},
		"password": func(fl validator.FieldLevel) bool {
			password := fl.Field().String()
			//var isValid bool
			isValid := true
			switch {
			case regexp.MustCompile(".*[A-Z].*").MatchString(password):
				fallthrough
			case regexp.MustCompile(".*[a-z].*").MatchString(password):
				fallthrough
			case regexp.MustCompile(".*\\d.*").MatchString(password):
				fallthrough
			case regexp.MustCompile(".*[@*#$%^&+=!].*").MatchString(password):
				fallthrough
			case regexp.MustCompile(".{8,20}").MatchString(password):
				isValid = true
			}
			return isValid
		},
		"first_name": func(fl validator.FieldLevel) bool {
			name := fl.Field().String()
			regex := regexp.MustCompile("^[a-zA-Z ]+$")
			return regex.MatchString(name)
		},
		"last_name": func(fl validator.FieldLevel) bool {
			surname := fl.Field().String()
			regex := regexp.MustCompile("^[a-zA-Z ]+$")
			return regex.MatchString(surname)
		},
		"phone": func(fl validator.FieldLevel) bool {
			phone := fl.Field().String()
			regex := regexp.MustCompile("^[0-9]{10}$")
			return regex.MatchString(phone)
		},
		"image_buffer_1": func(fl validator.FieldLevel) bool {
			imageBuffer := fl.Field().String()
			return len(imageBuffer) > 0
		},
		"image_name_1": func(fl validator.FieldLevel) bool {
			imageName := fl.Field().String()
			return len(imageName) > 0
		},
	}
	
	validations *validator.Validate
	once sync.Once
)

type(
	PasetoToken struct {
		tokenTTL   time.Duration
		refreshTTL time.Duration
	} 

	Payload struct{
		ID string `json:"id"`
		IssuedAt time.Time `json:"issued_at"`
		ExpiredAt time.Time `json:"expired_at"`
	}

	UserAccess struct {
		User 			*models.User
		AccessToken 	string 			`json:"access_token"`
		AccessPublic 	string 			`json:"access_public"`
		RefreshToken 	string 			`json:"refresh_token"`
		RefreshPublic 	string 			`json:"refresh_public"`
	}

	UserRepository struct{
		db *pgxpool.Pool
	}
	
)

func IsAuthorized(c fiber.Ctx) error {

	
	if !isValidToken(c) {
		return redirectToLogin(c, fiber.StatusUnauthorized, "authorization header is not provided or invalid")
	}

	if !isValidPublicKey(c) {
		return redirectToLogin(c, fiber.StatusUnauthorized, "public key is not provided")
	}

	token := getAccessToken(c)
	publicKey := getAccessPublicKey(c)


	DecodeToken := func (pt *PasetoToken, pasetoToken, publicKeyHex string) (*Payload, error) {
		publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKeyHex)
		if err != nil {
			return nil, err
		}
	
		parser := paseto.NewParser()
		parsedToken, err := parser.ParseV4Public(publicKey, pasetoToken, nil)
		if err != nil {
			return nil, err
		}
	
		payload := new(Payload)
		expiredAt, err := parsedToken.GetExpiration()
		if err != nil {
			return nil, err
		}

		Valid := func (payload *Payload)  error {
			if !time.Now().After(payload.ExpiredAt) {
				return ErrExpiredToken
			}
			return nil
		}

		err = Valid(payload)
		if err != nil {
			return nil, err
		}
	
		issuedAt, err := parsedToken.GetIssuedAt()
		if err != nil {
			return nil, err
		}
	
		id, err := parsedToken.GetString("id")
		if err != nil {
			return nil, err
		}
	
		payload = &Payload{
			ID:        id,
			IssuedAt:  issuedAt,
			ExpiredAt: expiredAt,
		}
	
		return payload, nil
	
	}

	paseto := &PasetoToken{}
	payload, err := DecodeToken(paseto ,token, publicKey)
	if err != nil {
		return redirectToLogin(c, fiber.StatusUnauthorized, "invalid access token")
	}

	c.Locals(AuthPayload, payload)
	return c.Next()
}

func GetUserDetail(c fiber.Ctx) error {
	payload := c.Locals(AuthPayload).(*Payload)
	
	GetByID := func (ctx context.Context, r *UserRepository, id string) (*models.User, error) {
		userQuery := struct {
			UserID        string
			Name      	  string
			Surname   	  string
			Email     	  string
			Phone		  sql.NullString
			PhotoUrl	  sql.NullString
			Password  	  string
			CreatedAt 	  time.Time
		}{}
		query := `
		SELECT CAST(user_id AS VARCHAR(64)) as UserID, 
		   first_name, 
		   last_name, 
		   email, 
		   phone,
		   photo_url,
		   password, 
		   created_at 
		FROM Users 
		WHERE user_id = $1 
			  AND password IS NOT NULL 
			  AND email IS NOT NULL;
		`
		err := r.db.QueryRow(ctx, query, id).Scan(&userQuery.UserID, &userQuery.Name, &userQuery.Surname, &userQuery.Email, &userQuery.Phone, &userQuery.PhotoUrl, &userQuery.Password, &userQuery.CreatedAt)
		if err != nil {
			return nil, err
		}
	
		userData := &models.User{
			UserID:    userQuery.UserID,
			Name:      userQuery.Name,
			Surname:   userQuery.Surname,
			Email:     userQuery.Email,
			Password:  userQuery.Password,
			CreatedAt: userQuery.CreatedAt,
		}

		

		// NULL kontrolü
		if userQuery.Phone.Valid {
			userData.Phone = userQuery.Phone.String
		}

		if userQuery.PhotoUrl.Valid {
			userData.PhotoUrl = &userQuery.PhotoUrl.String
		}

		return userData, nil
	}

	GetUserByID := func (ctx context.Context, r *UserRepository,id string) (*models.User, error) {
		userModel, err := GetByID(ctx, r , id)
		if err != nil {
			return nil, err
		}
	
		return userModel, nil
	}
	repo := &UserRepository{db: database.DBPool}
	userAggregate, err := GetUserByID(c.Context(), repo, payload.ID)
	if err != nil {
		return response.Error_Response(c, "error while trying to get user detail", err, nil, fiber.StatusBadRequest)
	}

	
	GetUserModelToDto := func (userData *models.User) *dto.GetUserResponse {

		var photoUrl string
		if userData.PhotoUrl != nil {
			photoUrl = *userData.PhotoUrl
		}
		return &dto.GetUserResponse{
			UserID:    userData.UserID,
			Name:      userData.Name,
			Surname:   userData.Surname,
			Email:     userData.Email,
			Phone: 		userData.Phone,
			PhotoUrl:  photoUrl,
			CreatedAt: userData.CreatedAt,
		}
	}

	userResponse := GetUserModelToDto(userAggregate)
	
	c.Locals(UserDetail, userResponse)
	//c.Locals(AuthPayload, payload) // payload'ı context'e kaydet
	return c.Next()
}

func redirectToLogin(c fiber.Ctx, statusCode int, message string) error {
	c.Status(statusCode).JSON(fiber.Map{
		"error": message,
	})
	return c.Redirect().To("/")
}

func CreateNewValidator() *validator.Validate {
	once.Do(func() {
		validations = validator.New()
		for key, value := range customValidators {
			validations.RegisterValidation(key, value)
		}
	})
	return validations
}

func ValidateRequestByStruct[T any](s T) []*response.ValidationMessage {
	validate := CreateNewValidator()
	var allErrors []*response.ValidationMessage
	err := validate.Struct(s)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			allErrors = append(allErrors, &response.ValidationMessage{
				FailedField: "N/A",
				Tag:         "invalid",
				Message:     err.Error(),
			})
			return allErrors
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, err := range validationErrors {
				var element response.ValidationMessage
				element.FailedField = err.Field()
				element.Tag = err.Tag()
				element.Message = fmt.Sprintf("Validation failed for field '%s' with tag '%s'", err.Field(), err.Tag()) // Özel hata mesajı
				allErrors = append(allErrors, &element)
			}
		}
	}
	return allErrors
}

func LoginValidation(c fiber.Ctx) error {
	data := new(dto.UserLoginRequest)
	body := c.Body()
	err := json.Unmarshal(body, &data)
	if err != nil {
		return response.Error_Response(c, "invalid request body", err, nil, fiber.StatusBadRequest)
	}

	validationErrors := ValidateRequestByStruct(data)
	if len(validationErrors) > 0 {
		return response.Error_Response(c, "validation failed", nil, validationErrors, fiber.StatusUnprocessableEntity)
	}

	return c.Next()
}

func RegisterValidation(c fiber.Ctx) error {
	data := new(dto.UserRegisterRequest)
	body := c.Body()
	err := json.Unmarshal(body, &data)
	if err != nil {
		return response.Error_Response(c, "invalid request body", err, nil, fiber.StatusBadRequest)
	}

	// Check if Password matches ConfirmPassword
    if data.Password != data.ConfirmPassword {
        return response.Error_Response(c, "passwords do not match", nil, nil, fiber.StatusBadRequest)
    }
	validationErrors := ValidateRequestByStruct(data)
	if len(validationErrors) > 0 {
		return response.Error_Response(c, "validation failed", nil, validationErrors, fiber.StatusUnprocessableEntity)
	}

	return c.Next()
}

func RateLimiter(max int, expiration time.Duration) func(fiber.Ctx) error {
	return limiter.New(limiter.Config{Max: max, Expiration: expiration, LimitReached: limitReachedFunc, KeyGenerator: func(c fiber.Ctx) string {
		remoteIp := c.IP()
		if c.Get("X-NginX-Proxy") == "true" {
			remoteIp = c.Get("X-Real-IP")
		}

		return remoteIp
	}})
}

func limitReachedFunc(c fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(response.ErrorResponse{
		Message: fiber.ErrTooManyRequests.Message,
		Status:  fiber.StatusTooManyRequests,
	})
}

func isValidToken(c fiber.Ctx) bool {
	token := c.Cookies(AccessToken)
	if token == "" {
		return false
	}

	fields := strings.Fields(token)
	if len(fields) != 2 || fields[0] != AuthType {
		return false
	}

	return true
}

func isValidPublicKey(c fiber.Ctx) bool {
	publicKey := c.Cookies(AccessPublic)
	return publicKey != ""
}

func getAccessToken(c fiber.Ctx) string {
	fields := strings.Fields(c.Cookies(AccessToken))
	return fields[1]
}

func getAccessPublicKey(c fiber.Ctx) string {
	return c.Cookies(AccessPublic)
}

func NewPayload(userID string, duration time.Duration)(*Payload, error) {
	payload := &Payload{
		ID: userID,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func Login(c fiber.Ctx) error {
    reqBody := new(dto.UserLoginRequest)

    body := c.Body()
    if err := json.Unmarshal(body, &reqBody); err != nil {
        return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
    }

    Login := func(ctx context.Context, email, password string) (*UserAccess, error) {
        r := &UserRepository{db: database.DBPool}

        GetUserPassword := func(ctx context.Context, r *UserRepository, email string) (string, error) {
            var hashedPassword string // Veritabanından gelen hash'lenmiş şifre
            query := `
            SELECT password 
            FROM users 
            WHERE email = $1;
            `
            err := r.db.QueryRow(ctx, query, email).Scan(&hashedPassword)
            if err != nil {
                return "", err
            }
            return hashedPassword, nil
        }

        hashedPassword, err := GetUserPassword(ctx, r, email)
        if err != nil {
            return nil, err
        }

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			return nil, errors.New("Şifreler eşleşmiyor")
		}


        GetByEmail := func(ctx context.Context, r *UserRepository, email string) (*models.User, error) {
            userQuery := struct {
                UserID    string
                Name      string
                Surname   string
                Email     string
                Password  string
                CreatedAt time.Time
            }{}
            query := `
            SELECT CAST(user_id AS VARCHAR(64)) as ID, 
               first_name, 
               last_name, 
               email, 
               password, 
               created_at 
            FROM Users 
            WHERE Email = $1 
                  AND password IS NOT NULL 
                  AND email IS NOT NULL;
            `
            err := r.db.QueryRow(ctx, query, email).Scan(&userQuery.UserID, &userQuery.Name, &userQuery.Surname, &userQuery.Email, &userQuery.Password, &userQuery.CreatedAt)
            if err != nil {
                return nil, err
            }

            userData := &models.User{
                UserID:    userQuery.UserID,
                Name:      userQuery.Name,
                Surname:   userQuery.Surname,
                Email:     userQuery.Email,
                Password:  userQuery.Password,
                CreatedAt: userQuery.CreatedAt,
            }
            return userData, nil
        }

        userModel, err := GetByEmail(ctx, r, email)
        if err != nil {
            return nil, err
        }

        CreateToken := func(userID string, tokenTTL time.Duration) (string, string, *Payload, error) {
            duration := tokenTTL

           

            payload, err := NewPayload(userID, duration)
            if err != nil {
                return "", "", nil, err
            }

            tokenPaseto := paseto.NewToken()
            tokenPaseto.SetExpiration(payload.ExpiredAt)
            tokenPaseto.SetIssuedAt(payload.IssuedAt)
            tokenPaseto.SetString("id", payload.ID)
            secretKey := paseto.NewV4AsymmetricSecretKey()
            publicKey := secretKey.Public().ExportHex()
            encrypted := tokenPaseto.V4Sign(secretKey, nil)

            return encrypted, publicKey, payload, nil
        }

        accessToken, publicKey, accessPayload, err := CreateToken(userModel.UserID, time.Hour*24)
        if err != nil {
            return nil, err
        }

        CreateRefreshToken := func(refreshTTL time.Duration, payload *Payload) (string, string, error) {
            tokenPaseto := paseto.NewToken()
            payload.ExpiredAt = payload.ExpiredAt.Add(refreshTTL)
            tokenPaseto.SetExpiration(payload.ExpiredAt)
            tokenPaseto.SetIssuedAt(payload.IssuedAt)
            tokenPaseto.SetString("id", payload.ID)
            secretKey := paseto.NewV4AsymmetricSecretKey()
            publicKey := secretKey.Public().ExportHex()
            encrypted := tokenPaseto.V4Sign(secretKey, nil)
            return encrypted, publicKey, nil
        }

        refreshToken, refreshPublicKey, err := CreateRefreshToken(time.Hour*24, accessPayload)
        if err != nil {
            return nil, err
        }

        NewUserAccess := func(user *models.User, accessToken, accessPublic, refreshToken, refreshPublic string) *UserAccess {
            return &UserAccess{
                User:         user,
                AccessToken:  accessToken,
                AccessPublic: accessPublic,
                RefreshToken: refreshToken,
                RefreshPublic: refreshPublic,
            }
        }

        sessionModel := NewUserAccess(userModel, accessToken, publicKey, refreshToken, refreshPublicKey)

        return sessionModel, nil
    }




    userData, err := Login(c.Context(), reqBody.Email, reqBody.Password)
    if err != nil {
        return response.Error_Response(c, "error while trying to login", err, nil, fiber.StatusBadRequest)
    }


	CreateSession := func (ctx context.Context, user *models.User, accessToken string, ipAddress string, userAgent string) error {
		r := &UserRepository{db: database.DBPool}
		session := &models.Session{
			SessionID:  uuid.New().String(),
			UserID: user.UserID,
			Token: accessToken,
			IPAdress: ipAddress,
			UserAgent: userAgent,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(5*time.Second),//Örneğin, 24 saat geçerli bir oturum
			LastAccess: time.Now(),
			IsActive: true,
			Location: "",
		}

		query := `
			INSERT INTO sessions (
				session_id, user_id, token, ip_address, user_agent, created_at, expires_at, last_access, is_active, location
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
		`
		_, err := r.db.Exec(ctx, query, session.SessionID, session.UserID, session.Token, session.IPAdress, session.UserAgent, session.CreatedAt, session.ExpiresAt, session.LastAccess, session.IsActive, session.Location)
		if err != nil {
			return err
		}
	
		return nil

	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	err = CreateSession(c.Context(), userData.User, userData.AccessToken, ipAddress, userAgent)
	if err != nil {
		return response.Error_Response(c, "error while trying to create session", err, nil, fiber.StatusBadRequest)
	}

    GetUserModelToDto := func(userData *models.User) *dto.GetUserResponse {
        return &dto.GetUserResponse{
            UserID:    userData.UserID,
            Name:      userData.Name,
            Surname:   userData.Surname,
            Email:     userData.Email,
            CreatedAt: userData.CreatedAt,
        }
    }

    userResponse := GetUserModelToDto(userData.User)
    bearerAccess := "Bearer " + userData.AccessToken
    fmt.Println(userData.AccessToken)
    c.Cookie(&fiber.Cookie{
        Name:     "id",
        Value:    userData.User.UserID,
        Expires:  time.Now().Add(3 * time.Hour),
        HTTPOnly: true,
        Secure:   true,
    })
    c.Cookie(&fiber.Cookie{
        Name:     "name",
        Value:    userData.User.Name + " " + userData.User.Surname,
        Expires:  time.Now().Add(3 * time.Hour),
        HTTPOnly: true,
        Secure:   true,
    })
    c.Cookie(&fiber.Cookie{
        Name:     AccessToken,
        Value:    bearerAccess,
        Expires:  time.Now().Add(time.Hour * 3),
        HTTPOnly: true,
        Secure:   true,
    })
    c.Cookie(&fiber.Cookie{
        Name:     AccessPublic,
        Value:    userData.AccessPublic,
        Expires:  time.Now().Add(time.Hour * 3),
        HTTPOnly: true,
        Secure:   true,
    })

    bearerRefresh := "Bearer " + userData.RefreshToken
    c.Cookie(&fiber.Cookie{
        Name:     RefreshToken,
        Value:    bearerRefresh,
        Expires:  time.Now().Add(24 * time.Hour),
        HTTPOnly: true,
        Secure:   true,
    })

    c.Cookie(&fiber.Cookie{
        Name:     RefreshPublic,
        Value:    userData.RefreshPublic,
        Expires:  time.Now().Add(24 * time.Hour),
        HTTPOnly: true,
        Secure:   true,
    })

    return response.Success_Response(c, userResponse, "Kullanıcı başarıyla giriş yaptı.", fiber.StatusOK)
}

func Register(c fiber.Ctx) error {
	reqBody := new(dto.UserRegisterRequest)
	body := c.Body()
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	Register := func (ctx context.Context, first_name, last_name, email, phone, password string)(*models.User, error) {
		newUser := &models.User{
			UserID: uuid.New().String(),
			Name: first_name,
			Surname: last_name,
			Email: email,
			Phone: phone,
			Password: password,
		}

		Create := func(ctx context.Context, r *UserRepository, user *models.User) error {
	
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			user.CreatedAt = time.Now()
			user.Password = string(hashedPassword)
		
			query := `INSERT INTO users (user_id, first_name, last_name, email, phone, password, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7);`
			_, err = r.db.Exec(ctx, query, user.UserID, user.Name, user.Surname, user.Email, user.Phone, user.Password, user.CreatedAt)
			if err != nil {
				return err
			}
			return nil
		}
		repo := &UserRepository{db: database.DBPool}
		err := Create(ctx, repo, newUser)
		if err != nil {
			return nil, err
		}
		return newUser, nil
	}
	
	newUser, err := Register(c.Context(), reqBody.Name, reqBody.Surname, reqBody.Email, reqBody.Phone, reqBody.Password)
	if err != nil {
		return response.Error_Response(c, "Error while trying to register user", err, nil, fiber.StatusBadRequest)
	}

	GetUserModelToDto := func (userData *models.User) *dto.GetUserResponse {
		return &dto.GetUserResponse{
			UserID:    userData.UserID,
			Name:      userData.Name,
			Surname:   userData.Surname,
			Email:     userData.Email,
			CreatedAt: userData.CreatedAt,
		}
	}

	userResponse := GetUserModelToDto(newUser)
	return response.Success_Response(c, userResponse, "user registered succesfully", fiber.StatusCreated)
}

func Logout(c fiber.Ctx) error {
	userID := c.Cookies("id")
	fmt.Println(userID)
	if userID == "" {
		return response.Error_Response(c, "no user id in cookies", nil, nil, fiber.StatusBadRequest)
	}

	UpdateSessionStatus := func (ctx context.Context, userID string) error {
		r := &UserRepository{db: database.DBPool}
		fmt.Println(userID)
		query := `
			UPDATE sessions
			SET is_active = false
			WHERE user_id = $1;
		`
		_, err := r.db.Exec(ctx, query, userID)
		fmt.Println("buradaaaaaaaaaa")
		return err
	}

	err := UpdateSessionStatus(c.Context(), userID)
	if err != nil {
		return response.Error_Response(c, "session is not updated", err, nil, fiber.StatusBadRequest)
	}

	// Cookie silme işlemleri
	c.ClearCookie(AccessToken)
	c.ClearCookie(AccessPublic)
	c.ClearCookie(RefreshToken)
	c.ClearCookie(RefreshPublic)
	c.ClearCookie("id")
	c.ClearCookie("name")
	path := "home"
	return c.Render(path, fiber.Map{
		"Title": "Kömürcü Emlak - Anasayfa",
	})
}

