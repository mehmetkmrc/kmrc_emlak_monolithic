package property

import (
	"context"
	"encoding/json"
	"fmt"
	"kmrc_emlak_mono/auth"
	"kmrc_emlak_mono/database"
	"kmrc_emlak_mono/dto"
	
	"kmrc_emlak_mono/models"
	"kmrc_emlak_mono/response"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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


	PropertyDetailsCreateRequestModel := func (dto.PropertyDetailsCreateRequest) (*models.PropertyDetails, error) {
		propertyDetail := new(models.PropertyDetails)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		propertyDetail = &models.PropertyDetails{
			PropertyID: property_id,
			PropertyDetailsID: uuid.New(),
			Area: reqBody.Area,
			Bedrooms:  reqBody.Bedrooms,
			Bathrooms: reqBody.Bathrooms,
			Parking: reqBody.Parking,
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
			OthersName: reqBody.OthersName,
			OthersChecked: reqBody.OthersChecked,
		}
		return amenities, nil
	}
	amenitiesModel, err := AmenitiesCreateRequestModel(*reqBody)
	if err != nil {
		return response.Error_Response(c, "error while trying to convert amenities create request model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, amenitiesModel *models.Amenities) (*models.Amenities, error) {
		query := `INSERT INTO amenities(amenities_id, property_id, wifi, pool, security, laundry_room, equipped_kitchen, air_conditioning, parking, garage_atached, fireplace, window_covering, backyard, fitness_gym, elevator, others_name, others_checked) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING amenities_id, property_id, wifi, pool, security, laundry_room, equipped_kitchen, air_conditioning, parking, garage_atached, fireplace, window_covering, backyard, fitness_gym, elevator, others_name, others_checked`
		queryRow := q.dbPool.QueryRow(ctx, query, amenitiesModel.AmenitiesID, amenitiesModel.PropertyID, amenitiesModel.Wifi, amenitiesModel.Pool, amenitiesModel.Security, amenitiesModel.LaundryRoom, amenitiesModel.EquippedKitchen, amenitiesModel.AirConditioning, amenitiesModel.Parking, amenitiesModel.GarageAtached, amenitiesModel.Fireplace, amenitiesModel.WindowCovering, amenitiesModel.Backyard, amenitiesModel.FitnessGym, amenitiesModel.Elevator, amenitiesModel.OthersName, amenitiesModel.OthersChecked)
		err := queryRow.Scan(&amenitiesModel.AmenitiesID, &amenitiesModel.PropertyID, &amenitiesModel.Wifi, &amenitiesModel.Pool, &amenitiesModel.Security, &amenitiesModel.LaundryRoom, &amenitiesModel.EquippedKitchen, &amenitiesModel.AirConditioning, &amenitiesModel.Parking, &amenitiesModel.GarageAtached, &amenitiesModel.Fireplace, &amenitiesModel.WindowCovering, &amenitiesModel.Backyard, &amenitiesModel.FitnessGym, &amenitiesModel.Elevator, &amenitiesModel.OthersName, &amenitiesModel.OthersChecked)
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

func AddPropertyMedia(c fiber.Ctx) error{
	reqBody := new(dto.PropertyMediaCreateRequest)
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
	PropertyMediaCreateRequestModel := func (dto.PropertyMediaCreateRequest) (*models.PropertyMedia, error) {
		propertyMedia := new(models.PropertyMedia)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		propertyMedia = &models.PropertyMedia{
			PropertyID: property_id,
			PropertyMediaID: uuid.New(),
			ImageID: propertyMedia.ImageID,
			Type: models.GalleryType(reqBody.Type),
		}
		return propertyMedia, nil
	}
	propertyMediaModel, err := PropertyMediaCreateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert propertyMedia create request model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, propertyMediaModel *models.PropertyMedia) (*models.PropertyMedia, error) {
		query := `INSERT INTO property_media(property_media_id, property_id, image_id, type) VALUES($1, $2, $3, $4) RETURNING property_media_id, property_id, image_id, type`
		queryRow := q.dbPool.QueryRow(ctx, query,  propertyMediaModel.PropertyMediaID, propertyMediaModel.PropertyID, propertyMediaModel.ImageID, propertyMediaModel.Type)
		err := queryRow.Scan(&propertyMediaModel.PropertyMediaID, &propertyMediaModel.PropertyID, &propertyMediaModel.ImageID, &propertyMediaModel.Type)
		if err != nil{
			return nil, err
		}
		return propertyMediaModel, nil		
	}
	AddPropertyMedia := func (ctx context.Context, propertyMedia *models.PropertyMedia) (*models.PropertyMedia, error) {
		repo := &PropertyRepository{dbPool: database.DBPool}
		return Insert(ctx, repo, propertyMedia)
	}
	propertyMedia, err := AddPropertyMedia(c.Context(), propertyMediaModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create PropertyMedia table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("PropertyMedia table created successfully! PropertyMedia: ", propertyMedia)
	return response.Success_Response(c, propertyMediaModel, "PropertyModel Created Successfully", fiber.StatusOK)
}

func AddImage(c fiber.Ctx) error{
	reqBody := new(dto.ImageCreateRequest)
	body := c.Body()
	if err := json.Unmarshal(body, reqBody);
	err != nil{
		return response.Error_Response(c, "error while trying to parse body", err, nil, fiber.StatusBadRequest)
	}
	ImageCreateRequestModel := func (dto.ImageCreateRequest) (*models.Image, error){
		image := new(models.Image)
		property_id, err := uuid.Parse(reqBody.PropertyID)
			if err != nil {
				return nil, err
			}
		image = &models.Image{
			PropertyID: property_id,
			ImageID: uuid.New(),
			ImageName: reqBody.ImageName,
			FilePath: reqBody.FilePath,
		}		
		return image, nil
	}
	imageModel, err := ImageCreateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert Image create request model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, imageModel *models.Image) (*models.Image, error) {
		query := `INSERT INTO images(property_id, image_id, name, file_path) VALUES($1, $2, $3) RETURNING property_id, image_id, name, file_path`
		queryRow := q.dbPool.QueryRow(ctx, query, imageModel.PropertyID, imageModel.ImageID, imageModel.ImageName, imageModel.FilePath)
		err := queryRow.Scan(&imageModel.PropertyID, &imageModel.ImageID, &imageModel.ImageName, &imageModel.FilePath)
		if err != nil {
			return nil, err
		}
		return imageModel, nil
	}
	AddImage := func (ctx context.Context, image *models.Image) (*models.Image, error) {
		repo := &PropertyRepository{
			dbPool: database.DBPool,
		}
		return Insert(ctx, repo, image)
	}
	image, err := AddImage(c.Context(), imageModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create image table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("Image table created successfully! Image: ", image)
	return response.Success_Response(c, image, "Image Model Created successfully", fiber.StatusOK)
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
			Distance: reqBody.Distance,
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
	reqBody := new(dto.PlansBrochuresCreateRequest)
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
	PlansBrochuresCreateRequestModel := func (dto.PlansBrochuresCreateRequest) (*models.PlansBrochures, error) {
		plansBrochures := new(models.PlansBrochures)
		property_id, err := uuid.Parse(reqBody.PropertyID)
		if err != nil {
			return nil, err
		}
		plansBrochures = &models.PlansBrochures{
			PropertyID: property_id,
			PlansBrochuresID: uuid.New(),
			FileType: reqBody.FileType,
			FilePath: reqBody.FilePath,
		}
		return plansBrochures, nil
	}
	plansBrochuresModel, err := PlansBrochuresCreateRequestModel(*reqBody)
	if err != nil{
		return response.Error_Response(c, "error while trying to convert PlansBrochures create request model", err, nil, fiber.StatusBadRequest)
	}
	Insert := func (ctx context.Context, q *PropertyRepository, plansBrochuresModel *models.PlansBrochures) (*models.PlansBrochures, error) {
		query := `INSERT INTO accordion_widget(plans_brochures_id, property_id, file_type,file_path) VALUES($1, $2, $3, $4) RETURNING plans_brochures_id, property_id, file_type,file_path`
		queryRow := q.dbPool.QueryRow(ctx, query, plansBrochuresModel.PlansBrochuresID, plansBrochuresModel.PropertyID, plansBrochuresModel.FileType, plansBrochuresModel.FilePath)
		err := queryRow.Scan(&plansBrochuresModel.PlansBrochuresID, &plansBrochuresModel.PropertyID, &plansBrochuresModel.FileType, &plansBrochuresModel.FilePath)
		if err != nil {
			return nil, err
		}
		return plansBrochuresModel, nil
	}
	AddPlansBrochures := func (ctx context.Context, plansBrochures *models.PlansBrochures) (*models.PlansBrochures, error) {
		repo := &PropertyRepository{
			dbPool: database.DBPool,
		}
		return Insert(ctx, repo, plansBrochures)
	}
	plansBrochures, err := AddPlansBrochures(c.Context(), plansBrochuresModel)
	if err != nil{
		return response.Error_Response(c, "error while trying to create PlansBrochures table", err, nil, fiber.StatusBadRequest)
	}
	zap.S().Info("PlansBrochures table created successfully! AccordionWidget: ", plansBrochures)
	return response.Success_Response(c, plansBrochures, "PlansBrochures Model Created successfully", fiber.StatusOK)
}