package property

import (
	"context"
	"encoding/json"
	"fmt"
	"kmrc_emlak_mono/auth"
	"kmrc_emlak_mono/database"
	"kmrc_emlak_mono/dto"
	"log"
	"os"

	"path/filepath"

	"kmrc_emlak_mono/models"
	"kmrc_emlak_mono/response"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	//"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)


type PropertyRepository struct {
	dbPool *pgxpool.Pool
	validate *validator.Validate
}

//Buradan itibaren kullanıcı tabanlı property id- userid tanımlaması olacak
func AddProperty(c fiber.Ctx) error {
	reqBody := new(dto.MainPropertyCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody); 
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}
	// UserID'yi context'ten al
	 payload, ok := c.Locals(auth.AuthPayload).(*auth.Payload)
	 if !ok {
		fmt.Println("payload boş döndü...")
		fmt.Println(c.Locals(auth.AuthPayload))
	 	return response.Error_Response(c, "payload not found in context", nil, nil, fiber.StatusInternalServerError)
	 }


	userIDString := payload.ID //string türünde UserID
	fmt.Println("Kullanıcı ID:", userIDString)

	// string UserID'yi uuid.UUID'ye dönüştür
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return response.Error_Response(c, "invalid user ID format", err, nil, fiber.StatusBadRequest)
	}

	
	// propertyID, done := c.Locals("propertyID").(uuid.UUID)
	// if !done {
	// 	fmt.Println("propertyID boş döndü...")
		
	//  	return response.Error_Response(c, "payload not found in context", nil, nil, fiber.StatusInternalServerError)
	//  }

	MainPropertyCreateRequestModel := func(req *dto.MainPropertyCreateRequest)(*models.Property, error) {
		mainProperty := new(models.Property)
		propertyID, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		mainProperty = &models.Property{
			UserID: userID,
			PropertyID: propertyID,
			TariffPlan: "extended",
			Date: time.Now(),
			PropertyTitle: "default",
		}
		return mainProperty, nil
	}

	propertyModel, err := MainPropertyCreateRequestModel(reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert property create request to model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, propertyModel *models.Property) (*models.Property, error) {
		query := `INSERT INTO property(user_id, property_id, tariff_plan, date, title) VALUES ($1, $2, $3, $4, $5) RETURNING user_id, property_id, tariff_plan, date, title`
		queryRow := q.dbPool.QueryRow(ctx, query, propertyModel.UserID, propertyModel.PropertyID, propertyModel.TariffPlan, propertyModel.Date, propertyModel.PropertyTitle)
		err := queryRow.Scan(&propertyModel.UserID, &propertyModel.PropertyID, &propertyModel.TariffPlan, &propertyModel.Date, &propertyModel.PropertyTitle)
		if err != nil{
			return nil, err
		}
		return propertyModel, nil
	}

	AddMainProperty := func (ctx context.Context, property *models.Property) (*models.Property, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Insert(ctx, repo, property)
	}

	property, err := AddMainProperty(c.Context(), propertyModel)
	if err != nil {
		return response.Error_Response(c, "error while trying to create main property", err, nil, fiber.StatusBadRequest)
	} 

	propertyID := propertyModel.PropertyID

	zap.S().Info("Property Created Successfuly! Property:", property)

	// **Context'e PropertyID'yi kaydet**
	c.Locals("propertyID", propertyID)

	return c.Next()
}

func AddPropertyDetails(c fiber.Ctx) error {
	reqBody := new(dto.PropertyDetailsCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}


	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }

	area, err := strconv.ParseFloat(reqBody.Area, 32)
	if err != nil{
		return err
	}


	bedrooms, err := strconv.Atoi(reqBody.Bedrooms)
	if err != nil{
		return err
	}
	bathrooms, err := strconv.Atoi(reqBody.Bathrooms)
	if err != nil{
		return err
	}
	parkings, err := strconv.Atoi(reqBody.Parking)
	if err != nil{
		return err
	}


	PropertyDetailsCreateRequestModel := func (dto.PropertyDetailsCreateRequest) (*models.PropertyDetails, error) {
		propertyDetail := new(models.PropertyDetails)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		propertyDetail = &models.PropertyDetails{
			PropertyID: property_id,
			PropertyDetailsID: uuid.New(),
			Area: float32(area),
			Bedrooms:  bedrooms,
			Bathrooms: bathrooms,
			Parking: parkings,
			Accomodation: reqBody.Accomodation,
			Website: reqBody.Website,
			PropertyMessage: reqBody.PropertyMessage,
		}
		return propertyDetail, nil
	}
	propertyDetailModel, err := PropertyDetailsCreateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert propertyDetail create request to model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, propertyDetailModel *models.PropertyDetails) (*models.PropertyDetails, error){
		query := `INSERT INTO property_details(property_details_id, property_id, area, bedrooms, bathrooms, parking, accomodation, website, property_message) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING property_details_id, property_id, area, bedrooms, bathrooms, parking, accomodation, website, property_message`
		queryRow := q.dbPool.QueryRow(ctx, query, propertyDetailModel.PropertyDetailsID ,propertyDetailModel.PropertyID, propertyDetailModel.Area, propertyDetailModel.Bedrooms, propertyDetailModel.Bathrooms, propertyDetailModel.Parking, propertyDetailModel.Accomodation, propertyDetailModel.Website, propertyDetailModel.PropertyMessage)
		err := queryRow.Scan(&propertyDetailModel.PropertyDetailsID, &propertyDetailModel.PropertyID, &propertyDetailModel.Area, &propertyDetailModel.Bedrooms, &propertyDetailModel.Bathrooms, &propertyDetailModel.Parking, &propertyDetailModel.Accomodation, &propertyDetailModel.Website, &propertyDetailModel.PropertyMessage )
		if err != nil{
			return nil, err
		}
		return propertyDetailModel, nil
	}
	AddPropertyDetails := func (ctx context.Context, propertyDetail *models.PropertyDetails) (*models.PropertyDetails, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Insert(ctx, repo, propertyDetail)
	}

	property_detail, err := AddPropertyDetails(c.Context(), propertyDetailModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create property detail", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("PropertDetails Created Successfully! PropertyDetails:", property_detail)
	return response.Success_Response(c, propertyDetailModel, "Property Model Created Successfully", fiber.StatusOK)
}

func AddVideoWidget(c fiber.Ctx) error {
	reqBody := new(dto.VideoWidgetCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}
	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }
	VideoWidgetCreateRequestModel := func (dto.VideoWidgetCreateRequest) (*models.VideoWidget, error) {
		videoWidget := new(models.VideoWidget)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		videoWidget = &models.VideoWidget{
			PropertyID: property_id,
			VideoWidgetID: uuid.New(),
			VideoExist: reqBody.VideoExist,
			VideoTitle: reqBody.VideoTitle,
			YouTubeUrl: reqBody.YouTubeUrl,
			VimeoUrl: reqBody.VimeoUrl,
		}
		return videoWidget, nil
	}
	videoWidgetModel, err := VideoWidgetCreateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert videoWidget create request to model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, videoWidgetModel *models.VideoWidget ) (*models.VideoWidget, error) {
		query := `INSERT INTO video_widget(property_id, video_widget_id, video_exist, video_title, youtube_url, vimeo_url) VALUES($1, $2, $3, $4, $5, $6) RETURNING property_id, video_widget_id, video_exist, video_title, youtube_url, vimeo_url`
		queryRow := q.dbPool.QueryRow(ctx, query, videoWidgetModel.PropertyID, videoWidgetModel.VideoWidgetID, videoWidgetModel.VideoExist, videoWidgetModel.VideoTitle, videoWidgetModel.YouTubeUrl, videoWidgetModel.VimeoUrl)
		err := queryRow.Scan(&videoWidgetModel.PropertyID, &videoWidgetModel.VideoWidgetID, &videoWidgetModel.VideoExist, &videoWidgetModel.VideoTitle, &videoWidgetModel.YouTubeUrl, &videoWidgetModel.VimeoUrl)
		if err != nil{
			return nil, err
		}
		return videoWidgetModel, nil
	}
	AddVideoWidget := func (ctx context.Context, videoWidget *models.VideoWidget) (*models.VideoWidget, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Insert(ctx, repo, videoWidget)
	}

	video_widget, err := AddVideoWidget(c.Context(), videoWidgetModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create video_widget", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("VideoWidget Created Successfully! VideoWidget:", video_widget)
	return response.Success_Response(c, videoWidgetModel, "VideoWidget Model created successfully", fiber.StatusOK)
}

func AddLocation(c fiber.Ctx) error{
	reqBody := new(dto.LocationCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil {
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}


	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }

	phoneInt, err := strconv.Atoi(reqBody.Phone) // Atoi fonksiyonunu çağır ve değerleri al
	if err != nil {
		// Dönüşüm hatası durumunda işlemi durdur ve hatayı döndür
		return fmt.Errorf("invalid phone number format: %w", err)
	}

	// Longitude'u float32'ye dönüştürme
	var longitudeFloat32 float32 // float32 olarak tanımla
	if reqBody.Longitude != "" {    // Boş değilse dönüştür
		longitudeFloat64, err := strconv.ParseFloat(reqBody.Longitude, 32)
		if err != nil {
			fmt.Println("Longitude dönüşüm hatası:", err) // Hatayı yazdır
			return fmt.Errorf("invalid longitude format: %w", err)
		}
		longitudeFloat32 = float32(longitudeFloat64) // float32'ye dönüştür
	}

	// Latitude'u float32'ye dönüştürme
	var latitudeFloat32 float32 // float32 olarak tanımla
	if reqBody.Latitude != "" {   // Boş değilse dönüştür
		latitudeFloat64, err := strconv.ParseFloat(reqBody.Latitude, 32)
		if err != nil {
			return fmt.Errorf("invalid latitude format: %w", err)
		}
		latitudeFloat32 = float32(latitudeFloat64) // float32'ye dönüştür
	}

	LocationCreateRequestModel := func (dto.LocationCreateRequest) (*models.Location, error) {
		location := new(models.Location)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		location = &models.Location{
			PropertyID: property_id,
			LocationID: uuid.New(),
			Phone: phoneInt,
			Email: reqBody.Email,
			City: models.CityLocation(reqBody.City),
			Address: reqBody.Address,
			Longitude: longitudeFloat32,
			Latitude: latitudeFloat32,
		}
		return location, nil
	}
	locationModel, err := LocationCreateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convet location create request model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, locationModel *models.Location) (*models.Location, error) {
		query := `INSERT INTO location(location_id, property_id, phone, email, city, address, longitude, latitude) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING location_id, property_id, phone, email, city, address, longitude, latitude`
		queryRow := q.dbPool.QueryRow(ctx, query,  locationModel.LocationID, locationModel.PropertyID, locationModel.Phone, locationModel.Email, locationModel.City, locationModel.Address, locationModel.Longitude, locationModel.Latitude)
		err := queryRow.Scan(&locationModel.LocationID, &locationModel.PropertyID, &locationModel.Phone, &locationModel.Email, &locationModel.City, &locationModel.Address, &locationModel.Longitude, &locationModel.Latitude)
		if err != nil{
			return nil, err
		}
		return locationModel, nil
	}
	AddLocation := func (ctx context.Context, location *models.Location) (*models.Location, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Insert(ctx, repo, location)
	}

	location, err := AddLocation(c.Context(), locationModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create location", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("Location table created successfully! Location: ", location)
	return response.Success_Response(c, locationModel, "Location Model Created Successfully", fiber.StatusOK)
}

func AddAmenities(c fiber.Ctx) error{
	reqBody := new(dto.AmenitiesCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}



	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }


	othersJSON, err := json.Marshal(reqBody.Others)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert amenities create request model", err, nil, fiber.StatusBadRequest)
	}	

	AmenitiesCreateRequestModel := func (dto.AmenitiesCreateRequest) (*models.Amenities, error) {
		amenities := new(models.Amenities)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}

		amenities = &models.Amenities{
			PropertyID: property_id,
			AmenitiesID: uuid.New(),
			Wifi: reqBody.Wifi,
			Pool: reqBody.Pool,
			Security: reqBody.Security,
			LaundryRoom: reqBody.LaundryRoom,
			EquippedKitchen: reqBody.EquippedKitchen,
			AirConditioning: reqBody.AirConditioning,
			Parking: reqBody.Parking,
			GarageAtached: reqBody.GarageAtached,
			Fireplace: reqBody.Fireplace,
			WindowCovering: reqBody.WindowCovering,
			Backyard: reqBody.Backyard,
			FitnessGym: reqBody.FitnessGym,
			Elevator: reqBody.Elevator,
			Others: othersJSON,
		}
		return amenities, nil
	}
	amenitiesModel, err := AmenitiesCreateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert amenities create request model", err, nil, fiber.StatusBadRequest)
	}

	

	Insert := func (ctx context.Context, q *PropertyRepository, amenitiesModel *models.Amenities) (*models.Amenities, error) {
		query := `INSERT INTO amenities(amenities_id, property_id, wifi, pool, security, laundry_room, equipped_kitchen, air_conditioning, parking, garage_atached, fireplace, window_covering, backyard, fitness_gym, elevator, others) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING amenities_id, property_id, wifi, pool, security, laundry_room, equipped_kitchen, air_conditioning, parking, garage_atached, fireplace, window_covering, backyard, fitness_gym, elevator, others`
		queryRow := q.dbPool.QueryRow(ctx, query, amenitiesModel.AmenitiesID, amenitiesModel.PropertyID, amenitiesModel.Wifi, amenitiesModel.Pool, amenitiesModel.Security, amenitiesModel.LaundryRoom, amenitiesModel.EquippedKitchen, amenitiesModel.AirConditioning, amenitiesModel.Parking, amenitiesModel.GarageAtached, amenitiesModel.Fireplace, amenitiesModel.WindowCovering, amenitiesModel.Backyard, amenitiesModel.FitnessGym, amenitiesModel.Elevator, amenitiesModel.Others)
		err := queryRow.Scan(&amenitiesModel.AmenitiesID, &amenitiesModel.PropertyID, &amenitiesModel.Wifi, &amenitiesModel.Pool, &amenitiesModel.Security, &amenitiesModel.LaundryRoom, &amenitiesModel.EquippedKitchen, &amenitiesModel.AirConditioning, &amenitiesModel.Parking, &amenitiesModel.GarageAtached, &amenitiesModel.Fireplace, &amenitiesModel.WindowCovering, &amenitiesModel.Backyard, &amenitiesModel.FitnessGym, &amenitiesModel.Elevator, &amenitiesModel.Others)
		if err != nil {
			return nil, err
		}
		return amenitiesModel, nil
	}
	AddAmenities := func (ctx context.Context, amenities *models.Amenities) (*models.Amenities, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Insert(ctx, repo, amenities)
	}
	amenities, err := AddAmenities(c.Context(), amenitiesModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create Amenities table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("Amenities table created successfully! Amenities: ", amenities)
	return response.Success_Response(c, amenitiesModel, "Ameniteies Model Created successfully", fiber.StatusOK)
}

func AddAccordionWidget(c fiber.Ctx) error{
	reqBody := new(dto.AccordionWidgetCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}
	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }
	AccordionWidgetCreateRequestModel := func (dto.AccordionWidgetCreateRequest) (*models.AccordionWidget, error) {
		accordionWidget := new(models.AccordionWidget)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		accordionWidget = &models.AccordionWidget{
			PropertyID: property_id,
			AccordionWidgetID: uuid.New(),
			AccordionExist: reqBody.AccordionExist,
			AccordionTitle: reqBody.AccordionTitle,
			AccordionDetails: reqBody.AccordionDetails,
		}
		return accordionWidget, nil
	}
	accordionWidgetModel, err := AccordionWidgetCreateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert AccordionWidget create request model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, accordionWidgetModel *models.AccordionWidget) (*models.AccordionWidget, error) {
		query := `INSERT INTO accordion_widget(accordion_widget_id, property_id, accordion_exist, accordion_title, accordion_details) VALUES($1, $2, $3, $4, $5) RETURNING accordion_widget_id, property_id, accordion_exist, accordion_title, accordion_details`
		queryRow := q.dbPool.QueryRow(ctx, query, accordionWidgetModel.AccordionWidgetID, accordionWidgetModel.PropertyID, accordionWidgetModel.AccordionExist, accordionWidgetModel.AccordionTitle, accordionWidgetModel.AccordionDetails)
		err := queryRow.Scan(&accordionWidgetModel.AccordionWidgetID, &accordionWidgetModel.PropertyID, &accordionWidgetModel.AccordionExist, &accordionWidgetModel.AccordionTitle, &accordionWidgetModel.AccordionDetails)
		if err != nil {
			return nil, err
		}
		return accordionWidgetModel, nil
	}
	AddAccordionWidget := func (ctx context.Context, accordionWidget *models.AccordionWidget) (*models.AccordionWidget, error) {
		repo := &PropertyRepository{
			dbPool: database.DBPool,
		}
		return Insert(ctx, repo, accordionWidget)
	}
	accordionWidget, err := AddAccordionWidget(c.Context(), accordionWidgetModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create Amenities table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("AccordionWidget table created successfully! AccordionWidget: ", accordionWidget)
	return response.Success_Response(c, accordionWidget, "AccordionWidget Model Created successfully", fiber.StatusOK)
}




func InsertImage(c fiber.Ctx) error {

	// property_id
	propertyIDStr := c.FormValue("property_id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		return response.Error_Response(c, "Invalid property ID", err, nil, fiber.StatusBadRequest)
	}

	// type (gallery / cover / plan)
	mediaType := c.FormValue("type")
	if mediaType == "" {
		return response.Error_Response(c, "media type is required", nil, nil, fiber.StatusBadRequest)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return response.Error_Response(c, "Error retrieving form", err, nil, fiber.StatusBadRequest)
	}

	files := form.File["image"]
	if len(files) == 0 {
		return response.Error_Response(c, "No images uploaded", nil, nil, fiber.StatusBadRequest)
	}

	tx, err := database.DBPool.Begin(c.Context())
	if err != nil {
		return response.Error_Response(c, "Transaction error", err, nil, fiber.StatusInternalServerError)
	}
	defer tx.Rollback(c.Context())

	var createdImages []models.Image

	for _, file := range files {

		imageID := uuid.New()

		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%s%s", imageID, ext)
		savePath := fmt.Sprintf("uploads/%s", fileName)

		if err := c.SaveFile(file, savePath); err != nil {
			return response.Error_Response(c, "File save error", err, nil, fiber.StatusInternalServerError)
		}

		// 1️⃣ images insert
		_, err = tx.Exec(c.Context(), `
			INSERT INTO images
			(image_id, property_id, url, original_name, media_type)
			VALUES ($1,$2,$3,$4,$5)
		`,
			imageID,
			propertyID,
			savePath,
			file.Filename,
			mediaType,
		)
		if err != nil {
			return response.Error_Response(c, "Image insert error", err, nil, fiber.StatusInternalServerError)
		}

		// 2️⃣ property_media insert
		_, err = tx.Exec(c.Context(), `
			INSERT INTO property_media
			(property_media_id, property_id, image_id, type)
			VALUES ($1,$2,$3,$4)
		`,
			uuid.New(),
			propertyID,
			imageID,
			mediaType,
		)
		if err != nil {
			return response.Error_Response(c, "Property media insert error", err, nil, fiber.StatusInternalServerError)
		}

		createdImages = append(createdImages, models.Image{
			ImageID:      imageID,
			PropertyID:   propertyID,
			Url:          savePath,
			OriginalName: file.Filename,
			MediaType:    mediaType,
		})
	}

	if err := tx.Commit(c.Context()); err != nil {
		return response.Error_Response(c, "Commit error", err, nil, fiber.StatusInternalServerError)
	}

	return response.Success_Response(c, createdImages, "Images uploaded successfully", fiber.StatusOK)
}


func AddBasicInfo(c fiber.Ctx) error{
	reqBody := new(dto.BasicInfoCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	// Middleware aracılığıyla aktarılan propertyID'yi alın
	propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	if !ok {
		return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	}
	 
	BasicInfoCreateRequestModel := func (dto.BasicInfoCreateRequest) (*models.BasicInfo, error) {
		basicInfo := new(models.BasicInfo)
		
		basicInfo = &models.BasicInfo{
			BasicInfoID: uuid.New(),
			PropertyID: propertyID,
			MainTitle: reqBody.MainTitle,
			Type: models.PropertyType(reqBody.Type),
			Category: models.PropertyCategory(reqBody.Category),
			Price: reqBody.Price,
			Keywords: reqBody.Keywords,
		}
		return basicInfo, nil
	}
	basicInfoModel, err := BasicInfoCreateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert basic_info create request model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, basicInfoModel *models.BasicInfo) (*models.BasicInfo, error) {
		query := `INSERT INTO basic_infos(basic_info_id, property_id, main_title, property_type, category, price, keywords) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING basic_info_id, property_id, main_title, property_type, category, price, keywords`
		queryRow := q.dbPool.QueryRow(ctx, query, basicInfoModel.BasicInfoID, basicInfoModel.PropertyID,basicInfoModel.MainTitle,  basicInfoModel.Type, basicInfoModel.Category, basicInfoModel.Price, basicInfoModel.Keywords)
		err := queryRow.Scan(&basicInfoModel.BasicInfoID, &basicInfoModel.PropertyID,&basicInfoModel.MainTitle,  &basicInfoModel.Type, &basicInfoModel.Category, &basicInfoModel.Price, &basicInfoModel.Keywords)
		if err != nil{
			return nil, err
		}
		return basicInfoModel, nil	
	}
	AddBasicInfo := func (ctx context.Context, basicInfo *models.BasicInfo ) (*models.BasicInfo, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Insert(ctx, repo, basicInfo)
	}
	basicInfo, err := AddBasicInfo(c.Context(), basicInfoModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create BasicInfo table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("basic_info table created successfully! basic_info: ", basicInfo)
	return response.Success_Response(c, basicInfoModel, "BasicInfo Created Successfully", fiber.StatusOK)
}

func AddNearby(c fiber.Ctx) error{
	reqBody := new(dto.NearbyCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}


	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }

	distance, err := strconv.Atoi(reqBody.Distance)
	if err != nil{
		return err
	}
	

	NearbyCreateRequestModel := func (dto.NearbyCreateRequest) (*models.Nearby, error){
		nearby := new(models.Nearby)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		nearby = &models.Nearby{
			PropertyID: property_id,
			NearbyID: uuid.New(),
			Places: models.PropertyNearby(reqBody.Places),
			Distance: distance,
		}
		return nearby, nil
	}
	nearbyModel, err := NearbyCreateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert propertyMedia create request model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, nearbyModel *models.Nearby) (*models.Nearby, error) {
		query := `INSERT INTO nearby(nearby_id, property_id, places, distance) VALUES($1, $2, $3, $4) RETURNING nearby_id, property_id, places, distance`
		queryRow := q.dbPool.QueryRow(ctx, query, nearbyModel.NearbyID, nearbyModel.PropertyID, nearbyModel.Places, nearbyModel.Distance)
		err := queryRow.Scan(&nearbyModel.NearbyID, &nearbyModel.PropertyID, &nearbyModel.Places, &nearbyModel.Distance)
		if err != nil{
			return nil, err
		}
		return nearbyModel, nil
	}
	AddNearby := func (ctx context.Context, nearby *models.Nearby) (*models.Nearby, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Insert(ctx, repo, nearby)
	}
	nearby, err := AddNearby(c.Context(), nearbyModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create Nearby table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("Nearby table created successfully! Nearby: ", nearby)
	return response.Success_Response(c, nearbyModel, "NearbyModel Created Successfully", fiber.StatusOK)
}

func AddPlansBrochures(c fiber.Ctx) error{
	// Dosyayı al
	file, err := c.FormFile("file_path") // "file_path" anahtarıyla dosyayı al
	if err != nil {
		return response.Error_Response(c, "Error retrieving the file", err, nil, fiber.StatusBadRequest)
	}

	// Property ID'yi al
	propertyIDStr := c.FormValue("property_id") // "property_id" anahtarıyla Property ID'yi al
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		return response.Error_Response(c, "Invalid property ID", err, nil, fiber.StatusBadRequest)
	}

	// Dosya Adını ve Yolunu Belirle
	plansBrochuresID := uuid.New()
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", plansBrochuresID, fileExt)
	savePath := fmt.Sprintf("uploads/%s", fileName)

	// Dosyayı Kaydet
	if err := c.SaveFile(file, savePath); err != nil {
		return response.Error_Response(c, "Error saving file", err, nil, fiber.StatusInternalServerError)
	}

	// PlansBrochures modelini oluştur
	plansBrochures := &models.PlansBrochures{
		PropertyID:     propertyID,
		PlansBrochuresID: plansBrochuresID,
		FileType:       file.Filename, // Dosya adını FileType olarak kullan
		FilePath:       savePath,      // Kayıt yolunu FilePath olarak kullan
	}

	// Veritabanına kaydet
	query := `INSERT INTO plans_brochures(plans_brochures_id, property_id, file_type, file_path) VALUES($1, $2, $3, $4) RETURNING plans_brochures_id, property_id, file_type, file_path`
	row := database.DBPool.QueryRow(c.Context(), query, plansBrochures.PlansBrochuresID, plansBrochures.PropertyID, plansBrochures.FileType, plansBrochures.FilePath)

	if err := row.Scan(&plansBrochures.PlansBrochuresID, &plansBrochures.PropertyID, &plansBrochures.FileType, &plansBrochures.FilePath); err != nil {
		return response.Error_Response(c, "Error inserting into database", err, nil, fiber.StatusInternalServerError)
	}

	// Başarılı yanıt döndür
	zap.S().Info("Plans and brochures saved successfully!", plansBrochures)
	return response.Success_Response(c, plansBrochures, "Plans and brochures uploaded successfully", fiber.StatusOK)
}





func EditProperty(c fiber.Ctx) error {
    reqBody := new(dto.MainPropertyUpdateRequest)
    body := c.Body()

    if err := json.Unmarshal(body, reqBody); err != nil {
        return response.Error_Response(c, "error while parsing body", err, nil, fiber.StatusBadRequest)
    }

    // Payload al
    payload, ok := c.Locals(auth.AuthPayload).(*auth.Payload)
    if !ok {
        return response.Error_Response(c, "payload not found", nil, nil, fiber.StatusUnauthorized)
    }

    userID, err := uuid.Parse(payload.ID)
    if err != nil {
        return response.Error_Response(c, "invalid user id", err, nil, fiber.StatusBadRequest)
    }

    // PropertyID body'den al
    propertyID, err := uuid.Parse(reqBody.PropertyID)
    if err != nil {
        return response.Error_Response(c, "invalid property id", err, nil, fiber.StatusBadRequest)
    }

    // Model
    property := &models.Property{
        PropertyID:    propertyID,
        UserID:        userID,
        PropertyTitle:    reqBody.PropertyTitle,
        
    }

	UpdateMainProperty :=func(ctx context.Context, property *models.Property) (*models.Property, error) {
    repo := &PropertyRepository{dbPool: database.DBPool}

    query := `
        UPDATE property
        SET 
            tariff_plan = $1,
            title = $2
        WHERE 
            property_id = $3
            AND user_id = $4
        RETURNING user_id, property_id, tariff_plan, date, title
    `

    row := repo.dbPool.QueryRow(
        ctx,
        query,
        property.TariffPlan,
        property.PropertyTitle,
        property.PropertyID,
        property.UserID,
    )

    err := row.Scan(
        &property.UserID,
        &property.PropertyID,
        &property.TariffPlan,
        &property.Date,
        &property.PropertyTitle,
    )

    if err != nil {
        return nil, err
    }

    return property, nil
}


    // Update repository
    updatedProperty, err := UpdateMainProperty(c.Context(), property)
    if err != nil {
        return response.Error_Response(c, "error while updating property", err, nil, fiber.StatusInternalServerError)
    }

    zap.S().Info("Property Updated Successfully:", updatedProperty)

    return c.JSON(updatedProperty)
}

func EditPropertyDetails(c fiber.Ctx) error {
	reqBody := new(dto.PropertyDetailsUpdateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}


	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }

	area, err := strconv.ParseFloat(reqBody.Area, 32)
	if err != nil{
		return err
	}


	bedrooms, err := strconv.Atoi(reqBody.Bedrooms)
	if err != nil{
		return err
	}
	bathrooms, err := strconv.Atoi(reqBody.Bathrooms)
	if err != nil{
		return err
	}
	parkings, err := strconv.Atoi(reqBody.Parking)
	if err != nil{
		return err
	}


	PropertyDetailsUpdateRequestModel := func (dto.PropertyDetailsUpdateRequest) (*models.PropertyDetails, error) {
		propertyDetail := new(models.PropertyDetails)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		propertyDetail = &models.PropertyDetails{
			PropertyID: property_id,
			PropertyDetailsID: uuid.New(),
			Area: float32(area),
			Bedrooms:  bedrooms,
			Bathrooms: bathrooms,
			Parking: parkings,
			Accomodation: reqBody.Accomodation,
			Website: reqBody.Website,
			PropertyMessage: reqBody.PropertyMessage,
		}
		return propertyDetail, nil
	}
	propertyDetailModel, err := PropertyDetailsUpdateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert propertyDetail create request to model", err, nil, fiber.StatusBadRequest)
	}
	Update := func (ctx context.Context, q *PropertyRepository, propertyDetailModel *models.PropertyDetails) (*models.PropertyDetails, error){
		query := `
					UPDATE property_details
					SET
						area = $1,
						bedrooms = $2,
						bathrooms = $3,
						parking = $4,
						accomodation = $5,
						website = $6,
						property_message = $7 
					WHERE
						property_id = $8
					RETURNING property_details_id, property_id, area, bedrooms, bathrooms, parking, accomodation, website, property_message`
		queryRow := q.dbPool.QueryRow(ctx, query, propertyDetailModel.Area, propertyDetailModel.Bedrooms, propertyDetailModel.Bathrooms, propertyDetailModel.Parking, propertyDetailModel.Accomodation, propertyDetailModel.Website, propertyDetailModel.PropertyMessage, propertyDetailModel.PropertyID)
		err := queryRow.Scan(&propertyDetailModel.PropertyDetailsID, &propertyDetailModel.PropertyID, &propertyDetailModel.Area, &propertyDetailModel.Bedrooms, &propertyDetailModel.Bathrooms, &propertyDetailModel.Parking, &propertyDetailModel.Accomodation, &propertyDetailModel.Website, &propertyDetailModel.PropertyMessage )
		if err != nil{
			return nil, err
		}
		return propertyDetailModel, nil
	}
	AddPropertyDetails := func (ctx context.Context, propertyDetail *models.PropertyDetails) (*models.PropertyDetails, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Update(ctx, repo, propertyDetail)
	}

	property_detail, err := AddPropertyDetails(c.Context(), propertyDetailModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create property detail", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("PropertDetails Updated Successfully! PropertyDetails:", property_detail)
	return response.Success_Response(c, propertyDetailModel, "Property Model Updated Successfully", fiber.StatusOK)
}

func EditVideoWidget(c fiber.Ctx) error {
	reqBody := new(dto.VideoWidgetUpdateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}
	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }
	VideoWidgetUpdateRequestModel := func (dto.VideoWidgetUpdateRequest) (*models.VideoWidget, error) {
		videoWidget := new(models.VideoWidget)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		videoWidget = &models.VideoWidget{
			PropertyID: property_id,
			VideoExist: reqBody.VideoExist,
			VideoTitle: reqBody.VideoTitle,
			YouTubeUrl: reqBody.YouTubeUrl,
			VimeoUrl: reqBody.VimeoUrl,
		}
		return videoWidget, nil
	}
	videoWidgetModel, err := VideoWidgetUpdateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert videoWidget create request to model", err, nil, fiber.StatusBadRequest)
	}
	Update := func (ctx context.Context, q *PropertyRepository, videoWidgetModel *models.VideoWidget ) (*models.VideoWidget, error) {
		query := `
			UPDATE video_widget
			SET		
				video_exist = $1, 
				video_title = $2, 
				youtube_url = $3, 
				vimeo_url = $4
			WHERE 
				property_id = $5
			RETURNING video_widget_id,  property_id, video_exist, video_title, youtube_url, vimeo_url`
		queryRow := q.dbPool.QueryRow(ctx, query, videoWidgetModel.VideoExist, videoWidgetModel.VideoTitle, videoWidgetModel.YouTubeUrl, videoWidgetModel.VimeoUrl, videoWidgetModel.PropertyID)
		err := queryRow.Scan(&videoWidgetModel.VideoWidgetID, &videoWidgetModel.PropertyID, &videoWidgetModel.VideoExist, &videoWidgetModel.VideoTitle, &videoWidgetModel.YouTubeUrl, &videoWidgetModel.VimeoUrl)
		if err != nil{
			return nil, err
		}
		return videoWidgetModel, nil
	}
	UpdateVideoWidget := func (ctx context.Context, videoWidget *models.VideoWidget) (*models.VideoWidget, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Update(ctx, repo, videoWidget)
	}

	video_widget, err := UpdateVideoWidget(c.Context(), videoWidgetModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create video_widget", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("VideoWidget Created Successfully! VideoWidget:", video_widget)
	return response.Success_Response(c, videoWidgetModel, "VideoWidget Model created successfully", fiber.StatusOK)
}

func EditLocation(c fiber.Ctx) error{
	reqBody := new(dto.LocationUpdateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil {
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}


	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }

	phoneInt, err := strconv.Atoi(reqBody.Phone) // Atoi fonksiyonunu çağır ve değerleri al
	if err != nil {
		// Dönüşüm hatası durumunda işlemi durdur ve hatayı döndür
		return fmt.Errorf("invalid phone number format: %w", err)
	}

	// Longitude'u float32'ye dönüştürme
	var longitudeFloat32 float32 // float32 olarak tanımla
	if reqBody.Longitude != "" {    // Boş değilse dönüştür
		longitudeFloat64, err := strconv.ParseFloat(reqBody.Longitude, 32)
		if err != nil {
			fmt.Println("Longitude dönüşüm hatası:", err) // Hatayı yazdır
			return fmt.Errorf("invalid longitude format: %w", err)
		}
		longitudeFloat32 = float32(longitudeFloat64) // float32'ye dönüştür
	}

	// Latitude'u float32'ye dönüştürme
	var latitudeFloat32 float32 // float32 olarak tanımla
	if reqBody.Latitude != "" {   // Boş değilse dönüştür
		latitudeFloat64, err := strconv.ParseFloat(reqBody.Latitude, 32)
		if err != nil {
			return fmt.Errorf("invalid latitude format: %w", err)
		}
		latitudeFloat32 = float32(latitudeFloat64) // float32'ye dönüştür
	}

	LocationUpdateRequestModel := func (dto.LocationUpdateRequest) (*models.Location, error) {
		location := new(models.Location)
		parsedID, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		location = &models.Location{
			PropertyID: parsedID,
			Phone: phoneInt,
			Email: reqBody.Email,
			City: models.CityLocation(reqBody.City),
			Address: reqBody.Address,
			Longitude: longitudeFloat32,
			Latitude: latitudeFloat32,
		}
		return location, nil
	}
	locationModel, err := LocationUpdateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convet location create request model", err, nil, fiber.StatusBadRequest)
	}
	Update := func (ctx context.Context, q *PropertyRepository, locationModel *models.Location) (*models.Location, error) {
		query := `
			UPDATE location
			SET
				phone = $1,
				email = $2,
				city = $3,
				address = $4,
				longitude = $5,
				latitude = $6
			WHERE
				property_id = $7
			RETURNING location_id, property_id, phone, email, city, address, longitude, latitude
		`

			queryRow := q.dbPool.QueryRow(ctx, query, locationModel.Phone, locationModel.Email, locationModel.City, locationModel.Address, locationModel.Longitude, locationModel.Latitude, locationModel.PropertyID)
			err := queryRow.Scan(&locationModel.LocationID, &locationModel.PropertyID, &locationModel.Phone, &locationModel.Email, &locationModel.City, &locationModel.Address, &locationModel.Longitude, &locationModel.Latitude)
			if err != nil{
				return nil, err
			}
			return locationModel, nil
	}
	UpdateLocation := func (ctx context.Context, location *models.Location) (*models.Location, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Update(ctx, repo, location)
	}

	location, err := UpdateLocation(c.Context(), locationModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create location", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("Location table created successfully! Location: ", location)
	return response.Success_Response(c, locationModel, "Location Model Created Successfully", fiber.StatusOK)
}

func EditAmenities(c fiber.Ctx) error{
	reqBody := new(dto.AmenitiesUpdateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}



	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }


	

	AmenitiesUpdateRequestModel := func (dto.AmenitiesUpdateRequest) (*models.Amenities, error) {
		amenities := new(models.Amenities)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		var othersJSON []byte

		if reqBody.Others != nil {
			othersJSON, err = json.Marshal(reqBody.Others)
			if err != nil {
				return nil, err
			}
		}

		amenities = &models.Amenities{
			PropertyID: property_id,
			AmenitiesID: uuid.Nil,
			Wifi: reqBody.Wifi,
			Pool: reqBody.Pool,
			Security: reqBody.Security,
			LaundryRoom: reqBody.LaundryRoom,
			EquippedKitchen: reqBody.EquippedKitchen,
			AirConditioning: reqBody.AirConditioning,
			Parking: reqBody.Parking,
			GarageAtached: reqBody.GarageAtached,
			Fireplace: reqBody.Fireplace,
			WindowCovering: reqBody.WindowCovering,
			Backyard: reqBody.Backyard,
			FitnessGym: reqBody.FitnessGym,
			Elevator: reqBody.Elevator,
			Others: othersJSON,
		}
		if reqBody.Others != nil {
			amenities.Others = othersJSON
		}
		return amenities, nil
	}
	amenitiesModel, err := AmenitiesUpdateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert amenities create request model", err, nil, fiber.StatusBadRequest)
	}
	Update := func (ctx context.Context, q *PropertyRepository, amenitiesModel *models.Amenities) (*models.Amenities, error) {
		query := `
			UPDATE amenities SET
				wifi = $1,
				pool = $2,
				security = $3,
				laundry_room = $4,
				equipped_kitchen = $5,
				air_conditioning = $6,
				parking = $7,
				garage_atached = $8,
				fireplace = $9,
				window_covering = $10,
				backyard = $11,
				fitness_gym = $12,
				elevator = $13
			`

			args := []any{
				amenitiesModel.Wifi,
				amenitiesModel.Pool,
				amenitiesModel.Security,
				amenitiesModel.LaundryRoom,
				amenitiesModel.EquippedKitchen,
				amenitiesModel.AirConditioning,
				amenitiesModel.Parking,
				amenitiesModel.GarageAtached,
				amenitiesModel.Fireplace,
				amenitiesModel.WindowCovering,
				amenitiesModel.Backyard,
				amenitiesModel.FitnessGym,
				amenitiesModel.Elevator,
			}

			argIndex := 14

			if reqBody.Others != nil {
				query += fmt.Sprintf(", others = $%d", argIndex)
				args = append(args, amenitiesModel.Others)
				argIndex++
			}

			query += fmt.Sprintf(`
			WHERE property_id = $%d
			RETURNING amenities_id, property_id, wifi, pool, security, laundry_room,
			equipped_kitchen, air_conditioning, parking, garage_atached, fireplace,
			window_covering, backyard, fitness_gym, elevator, others
			`, argIndex)

			args = append(args, amenitiesModel.PropertyID)

		queryRow := q.dbPool.QueryRow(ctx, query, amenitiesModel.Wifi, amenitiesModel.Pool, amenitiesModel.Security, amenitiesModel.LaundryRoom, amenitiesModel.EquippedKitchen, amenitiesModel.AirConditioning, amenitiesModel.Parking, amenitiesModel.GarageAtached, amenitiesModel.Fireplace, amenitiesModel.WindowCovering, amenitiesModel.Backyard, amenitiesModel.FitnessGym, amenitiesModel.Elevator, amenitiesModel.Others, amenitiesModel.PropertyID)
		err := queryRow.Scan(&amenitiesModel.AmenitiesID, &amenitiesModel.PropertyID, &amenitiesModel.Wifi, &amenitiesModel.Pool, &amenitiesModel.Security, &amenitiesModel.LaundryRoom, &amenitiesModel.EquippedKitchen, &amenitiesModel.AirConditioning, &amenitiesModel.Parking, &amenitiesModel.GarageAtached, &amenitiesModel.Fireplace, &amenitiesModel.WindowCovering, &amenitiesModel.Backyard, &amenitiesModel.FitnessGym, &amenitiesModel.Elevator, &amenitiesModel.Others)
		if err != nil {
			return nil, err
		}
		return amenitiesModel, nil
	}
	UpdateAmenities := func (ctx context.Context, amenities *models.Amenities) (*models.Amenities, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Update(ctx, repo, amenities)
	}
	amenities, err := UpdateAmenities(c.Context(), amenitiesModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create Amenities table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("Amenities table updated successfully! Amenities: ", amenities)
	return response.Success_Response(c, amenitiesModel, "Ameniteies Model Updated successfully", fiber.StatusOK)
}

func EditAccordionWidget(c fiber.Ctx) error{
	reqBody := new(dto.AccordionWidgetUpdateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}
	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }
	AccordionWidgetUpdateRequestModel := func (dto.AccordionWidgetUpdateRequest) (*models.AccordionWidget, error) {
		accordionWidget := new(models.AccordionWidget)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		accordionWidget = &models.AccordionWidget{
			PropertyID: property_id,
			AccordionWidgetID: uuid.New(),
			AccordionExist: reqBody.AccordionExist,
			AccordionTitle: reqBody.AccordionTitle,
			AccordionDetails: reqBody.AccordionDetails,
		}
		return accordionWidget, nil
	}
	accordionWidgetModel, err := AccordionWidgetUpdateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert AccordionWidget create request model", err, nil, fiber.StatusBadRequest)
	}
	Update := func (ctx context.Context, q *PropertyRepository, accordionWidgetModel *models.AccordionWidget) (*models.AccordionWidget, error) {
		query := `
			UPDATE accordion_widget
			SET
				accordion_exist = $1, 
				accordion_title = $2, 
				accordion_details = $3
			WHERE
				property_id = $4
			RETURNING accordion_widget_id, property_id, accordion_exist, accordion_title, accordion_details`
		queryRow := q.dbPool.QueryRow(ctx, query, accordionWidgetModel.AccordionExist, accordionWidgetModel.AccordionTitle, accordionWidgetModel.AccordionDetails, accordionWidgetModel.PropertyID)
		err := queryRow.Scan(&accordionWidgetModel.AccordionWidgetID, &accordionWidgetModel.PropertyID, &accordionWidgetModel.AccordionExist, &accordionWidgetModel.AccordionTitle, &accordionWidgetModel.AccordionDetails)
		if err != nil {
			return nil, err
		}
		return accordionWidgetModel, nil
	}
	UpdateAccordionWidget := func (ctx context.Context, accordionWidget *models.AccordionWidget) (*models.AccordionWidget, error) {
		repo := &PropertyRepository{
			dbPool: database.DBPool,
		}
		return Update(ctx, repo, accordionWidget)
	}
	accordionWidget, err := UpdateAccordionWidget(c.Context(), accordionWidgetModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create Amenities table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("AccordionWidget table created successfully! AccordionWidget: ", accordionWidget)
	return response.Success_Response(c, accordionWidget, "AccordionWidget Model Created successfully", fiber.StatusOK)
}





func EditBasicInfo(c fiber.Ctx) error{
	reqBody := new(dto.BasicInfoUpdateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}

	// Middleware aracılığıyla aktarılan propertyID'yi alın
	parsedID, err := uuid.Parse(reqBody.PropertyID)
if err != nil {
    return response.Error_Response(c, "invalid property_id", err, nil, fiber.StatusBadRequest)
}

	 
	BasicInfoUpdateRequestModel := func (dto.BasicInfoUpdateRequest) (*models.BasicInfo, error) {
		basicInfo := new(models.BasicInfo)
		
		basicInfo = &models.BasicInfo{
			PropertyID: parsedID,
			MainTitle: reqBody.MainTitle,
			Type: models.PropertyType(reqBody.Type),
			Category: models.PropertyCategory(reqBody.Category),
			Price: reqBody.Price,
			Keywords: reqBody.Keywords,
		}

		return basicInfo, nil
	}
	basicInfoModel, err := BasicInfoUpdateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert basic_info create request model", err, nil, fiber.StatusBadRequest)
	}
	
	Update := func (ctx context.Context, q *PropertyRepository, basicInfoModel *models.BasicInfo) (*models.BasicInfo, error) {
		query := `
			UPDATE basic_infos
			SET
				main_title = $1,
				property_type = $2,
				category = $3,
				price = $4,
				keywords = $5
			WHERE
				property_id = $6
			RETURNING basic_info_id, property_id, main_title, property_type, category, price, keywords
		`

		row := q.dbPool.QueryRow(ctx, query,
			basicInfoModel.MainTitle,   // $1
			basicInfoModel.Type,        // $2
			basicInfoModel.Category,    // $3
			basicInfoModel.Price,       // $4
			basicInfoModel.Keywords,    // $5
			basicInfoModel.PropertyID,  // $6
		)


		err := row.Scan(
			&basicInfoModel.BasicInfoID,
			&basicInfoModel.PropertyID,
			&basicInfoModel.MainTitle,
			&basicInfoModel.Type,
			&basicInfoModel.Category,
			&basicInfoModel.Price,
			&basicInfoModel.Keywords,
		)

		if err != nil {
			return nil, err
		}

		return basicInfoModel, nil
	}

	UpdateBasicInfo := func (ctx context.Context, basicInfo *models.BasicInfo ) (*models.BasicInfo, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Update(ctx, repo, basicInfo)
	}
	basicInfo, err := UpdateBasicInfo(c.Context(), basicInfoModel)
	if err != nil{
		zap.S().Info("UUID STRING:", fmt.Sprintf("%s", basicInfoModel.PropertyID))
		return response.Error_Response(c, "error while trying to update BasicInfo table", err, nil, fiber.StatusBadRequest)
	}
	
	zap.S().Info("basic_info table updated successfully! basic_info: ", basicInfo)
	return response.Success_Response(c, basicInfoModel, "BasicInfo Updated Successfully", fiber.StatusOK)
}

func EditNearby(c fiber.Ctx) error{
	reqBody := new(dto.NearbyUpdateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}


	// Middleware aracılığıyla aktarılan propertyID'yi alın
	// propertyID, ok := c.Locals("propertyID").(uuid.UUID)
	// if !ok {
	// 	return response.Error_Response(c, "propertyID not found in context", nil, nil, fiber.StatusBadRequest)
	// }

	distance, err := strconv.Atoi(reqBody.Distance)
	if err != nil{
		return err
	}
	

	NearbyUpdateRequestModel := func (dto.NearbyUpdateRequest) (*models.Nearby, error){
		nearby := new(models.Nearby)
		parsedID, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		nearby = &models.Nearby{
			PropertyID: parsedID,
			NearbyID: uuid.New(),
			Places: models.PropertyNearby(reqBody.Places),
			Distance: distance,
		}
		return nearby, nil
	}
	nearbyModel, err := NearbyUpdateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert propertyMedia create request model", err, nil, fiber.StatusBadRequest)
	}
	Update := func (ctx context.Context, q *PropertyRepository, nearbyModel *models.Nearby) (*models.Nearby, error) {
		query := `INSERT INTO nearby(nearby_id, property_id, places, distance) VALUES($1, $2, $3, $4) RETURNING nearby_id, property_id, places, distance`
		queryRow := q.dbPool.QueryRow(ctx, query, nearbyModel.NearbyID, nearbyModel.PropertyID, nearbyModel.Places, nearbyModel.Distance)
		err := queryRow.Scan(&nearbyModel.NearbyID, &nearbyModel.PropertyID, &nearbyModel.Places, &nearbyModel.Distance)
		if err != nil{
			return nil, err
		}
		return nearbyModel, nil
	}
	UpdateNearby := func (ctx context.Context, nearby *models.Nearby) (*models.Nearby, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Update(ctx, repo, nearby)
	}
	nearby, err := UpdateNearby(c.Context(), nearbyModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create Nearby table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("Nearby table updated successfully! Nearby: ", nearby)
	return response.Success_Response(c, nearbyModel, "NearbyModel Updated Successfully", fiber.StatusOK)
}

func EditPlansBrochures(c fiber.Ctx) error{
	// Dosyayı al
	file, err := c.FormFile("file_path") // "file_path" anahtarıyla dosyayı al
	if err != nil {
		return response.Error_Response(c, "Error retrieving the file", err, nil, fiber.StatusBadRequest)
	}

	// Property ID'yi al
	propertyIDStr := c.FormValue("property_id") // "property_id" anahtarıyla Property ID'yi al
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		return response.Error_Response(c, "Invalid property ID", err, nil, fiber.StatusBadRequest)
	}

	// Dosya Adını ve Yolunu Belirle
	plansBrochuresID := uuid.New()
	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s%s", plansBrochuresID, fileExt)
	savePath := fmt.Sprintf("uploads/%s", fileName)

	// Dosyayı Kaydet
	if err := c.SaveFile(file, savePath); err != nil {
		return response.Error_Response(c, "Error saving file", err, nil, fiber.StatusInternalServerError)
	}

	// PlansBrochures modelini oluştur
	plansBrochures := &models.PlansBrochures{
		PropertyID:     propertyID,
		PlansBrochuresID: plansBrochuresID,
		FileType:       file.Filename, // Dosya adını FileType olarak kullan
		FilePath:       savePath,      // Kayıt yolunu FilePath olarak kullan
	}

	// Veritabanına kaydet
	query := `
		UPDATE plans_brochures
		SET
			file_type = $1, 
			file_path = $2 
		WHERE
			plans_brochures_id = $3 
			AND property_id = $4 
		RETURNING plans_brochures_id, property_id, file_type, file_path`
	row := database.DBPool.QueryRow(c.Context(), query, plansBrochures.PlansBrochuresID, plansBrochures.PropertyID, plansBrochures.FileType, plansBrochures.FilePath)

	if err := row.Scan(&plansBrochures.PlansBrochuresID, &plansBrochures.PropertyID, &plansBrochures.FileType, &plansBrochures.FilePath); err != nil {
		return response.Error_Response(c, "Error updating into database", err, nil, fiber.StatusInternalServerError)
	}

	// Başarılı yanıt döndür
	zap.S().Info("Plans and brochures updated successfully!", plansBrochures)
	return response.Success_Response(c, plansBrochures, "Plans and brochures updated successfully", fiber.StatusOK)
}

func DeleteNearby(c fiber.Ctx) error {
    id := c.Params("nearbyID")
    parsedID, err := uuid.Parse(id)
    if err != nil {
        return c.Status(400).JSON("Invalid ID")
    }

    _, err = database.DBPool.Exec(context.Background(),
        "DELETE FROM nearby WHERE nearby_id = $1", parsedID)

    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.JSON(fiber.Map{
        "status": 200,
        "message": "Deleted",
    })
}

func DeleteImage(c fiber.Ctx) error {

	propertyMediaIDStr := c.Params("property_media_id")
	propertyMediaID, err := uuid.Parse(propertyMediaIDStr)
	if err != nil {
		return response.Error_Response(c, "Invalid property_media_id", err, nil, fiber.StatusBadRequest)
	}

	ctx := c.Context()

	tx, err := database.DBPool.Begin(ctx)
	if err != nil {
		return response.Error_Response(c, "Transaction error", err, nil, fiber.StatusInternalServerError)
	}
	defer tx.Rollback(ctx)

	var imageID uuid.UUID
	var imagePath string

	// 1️⃣ İlgili image bilgisi
	err = tx.QueryRow(ctx, `
		SELECT i.image_id, i.url
		FROM property_media pm
		JOIN images i ON pm.image_id = i.image_id
		WHERE pm.property_media_id = $1
	`, propertyMediaID).Scan(&imageID, &imagePath)

	if err != nil {
		return response.Error_Response(c, "Image not found", err, nil, fiber.StatusNotFound)
	}

	// 2️⃣ property_media sil
	_, err = tx.Exec(ctx, `
		DELETE FROM property_media
		WHERE property_media_id = $1
	`, propertyMediaID)
	if err != nil {
		return response.Error_Response(c, "property_media delete error", err, nil, fiber.StatusInternalServerError)
	}

	// 3️⃣ images sil
	_, err = tx.Exec(ctx, `
		DELETE FROM images
		WHERE image_id = $1
	`, imageID)
	if err != nil {
		return response.Error_Response(c, "image delete error", err, nil, fiber.StatusInternalServerError)
	}

	// 4️⃣ Diskten dosyayı sil
	if err := os.Remove(imagePath); err != nil {
		// dosya yoksa panic yapma
		log.Println("file delete warning:", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return response.Error_Response(c, "Commit error", err, nil, fiber.StatusInternalServerError)
	}

	return response.Success_Response(c, nil, "Image deleted successfully", fiber.StatusOK)
}


func DeleteProperty(c fiber.Ctx) error {
    propertyIDParam := c.Params("property_id")
    if propertyIDParam == "" {
        return response.Error_Response(
            c,
            "property_id is required",
            nil,
            nil,
            fiber.StatusBadRequest,
        )
    }

    propertyID, err := uuid.Parse(propertyIDParam)
    if err != nil {
        return response.Error_Response(
            c,
            "invalid property_id",
            err,
            nil,
            fiber.StatusBadRequest,
        )
    }

    Delete := func(ctx context.Context, q *PropertyRepository, propertyID uuid.UUID) error {
        query := `
            DELETE FROM property
            WHERE property_id = $1
        `
        _, err := q.dbPool.Exec(ctx, query, propertyID)
        return err
    }

    DeleteProperty := func(ctx context.Context, propertyID uuid.UUID) error {
        repo := &PropertyRepository{dbPool: database.DBPool}
        return Delete(ctx, repo, propertyID)
    }

    if err := DeleteProperty(c.Context(), propertyID); err != nil {
        return response.Error_Response(
            c,
            "error while trying to delete property",
            err,
            nil,
            fiber.StatusInternalServerError,
        )
    }

    zap.S().Info("Property deleted successfully:", propertyID)

    return response.Success_Response(
        c,
        nil,
        "Property deleted successfully",
        fiber.StatusOK,
    )
}


func PassiveProperty(c fiber.Ctx) error {
    propertyIDParam := c.Params("property_id")
    if propertyIDParam == "" {
        return response.Error_Response(c, "property_id is required", nil, nil, fiber.StatusBadRequest)
    }

    propertyID, err := uuid.Parse(propertyIDParam)
    if err != nil {
        return response.Error_Response(c, "invalid property_id", err, nil, fiber.StatusBadRequest)
    }

    Update := func(ctx context.Context, q *PropertyRepository, propertyID uuid.UUID) error {
        query := `
            UPDATE property
            SET property_status = 0
            WHERE property_id = $1
        `
        _, err := q.dbPool.Exec(ctx, query, propertyID)
        return err
    }

    PassiveProperty := func(ctx context.Context, propertyID uuid.UUID) error {
        repo := &PropertyRepository{dbPool: database.DBPool}
        return Update(ctx, repo, propertyID)
    }

    if err := PassiveProperty(c.Context(), propertyID); err != nil {
        return response.Error_Response(
            c,
            "error while trying to passive property",
            err,
            nil,
            fiber.StatusInternalServerError,
        )
    }

    zap.S().Info("Property set to passive:", propertyID)

    return response.Success_Response(
        c,
        nil,
        "Property set to passive successfully",
        fiber.StatusOK,
    )
}
