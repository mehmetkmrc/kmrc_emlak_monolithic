package web

import (
	"context"
	"fmt"
	"kmrc_emlak_mono/database"
	"kmrc_emlak_mono/models"
	"strconv"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PropertyRepository struct {
	dbPool *pgxpool.Pool
	
}

func LoginWeb(c fiber.Ctx) error {
	path := "login"
	return c.Render(path, fiber.Map{
		"Title": "Login",
	})
}


func HomeWeb(c fiber.Ctx) error {
	categoryFilter := c.Query("category")
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "6"))

	GetPropertiesByJoin := func(ctx context.Context) ([]*models.Property, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT
				p.property_id as property_id,
				bi.basic_info_id as basic_info_id,
				bi.property_type as property_type,
				bi.category as category,
				bi.main_title as main_title,
				bi.price as price,
				loc.property_id as property_id,
				loc.address as address,
				pd.property_id as property_id,
				pd.property_message as property_message,
				pd.bedrooms as bedrooms,
				pd.bathrooms as bathrooms,
				pd.area as area
			FROM
				property p
			LEFT JOIN
				basic_infos bi ON p.property_id = bi.property_id
			LEFT JOIN
				location loc ON p.property_id = loc.property_id
			LEFT JOIN
				property_details pd ON p.property_id = pd.property_id
		`)
		if err != nil {
			fmt.Println("Sorgu hatası: ", err)
			return nil, err
		}
		defer rows.Close()

		var properties []*models.Property

		for rows.Next() {
			var property models.Property
			var basicInfos models.BasicInfo
			var location models.Location
			var propertyDetails models.PropertyDetails

			err := rows.Scan(
				&property.PropertyID, &basicInfos.PropertyID, &basicInfos.Type, &basicInfos.Category, &basicInfos.MainTitle, &basicInfos.Price, &location.PropertyID, &location.Address, &propertyDetails.PropertyID, &propertyDetails.PropertyMessage, &propertyDetails.Bedrooms, &propertyDetails.Bathrooms, &propertyDetails.Area,
			)
			if err != nil {
				fmt.Println("Satır tarama hatası: ", err)
				continue // Hata durumunda sonraki satıra geç
			}

			property.BasicInfo = &basicInfos
			property.Location = &location
			property.PropertyDetails = &propertyDetails



			//Resimleri getiren fonksiyon
			GetImagesByPropertyID := func(ctx context.Context, propertyID uuid.UUID) ([]*models.Image, error) {
				rows, err := database.DBPool.Query(ctx, `
				SELECT image_id, property_id, name, file_path
				FROM images
				WHERE property_id = $1
			`, propertyID)
				if err != nil {
					fmt.Println("Resim sorgulama hatası: ", err)
					return nil, err
				}
				defer rows.Close()

				var images []*models.Image
				for rows.Next() {
					var image models.Image
					err := rows.Scan(&image.ImageID, &image.PropertyID, &image.ImageName, &image.FilePath)
					if err != nil {
						fmt.Println("Resim satırı tarama hatası: ", err)
						continue
					}
					images = append(images, &image)
				}

				if err := rows.Err(); err != nil {
					fmt.Println("Resim satırları yineleme hatası: ", err)
					return nil, err
				}

				return images, nil
			}
			// Resimleri getir
			images, err := GetImagesByPropertyID(ctx, property.PropertyID)
			if err != nil {
				fmt.Println("Resim getirme hatası: ", err)
				//Hata durumunda ne yapılacağına karar verin, örneğin boş bir dilim atayın
				property.PropertyMedia = []*models.PropertyMedia{}
			} else {
				// PropertyMedia'yı doldur
				propertyMedia := &models.PropertyMedia{
					PropertyID: property.PropertyID,
					Image:      images, // Resimleri doğrudan ata
				}
				property.PropertyMedia = []*models.PropertyMedia{propertyMedia} // Slice içinde sakla
			}

			properties = append(properties, &property)
		}
		return properties, nil
	}

	

	ctx := context.Background()
	properties, err := GetPropertiesByJoin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Verileri alırken hata oluştu")
	}

	if categoryFilter != "" {
		filtered := []*models.Property{}
		for _, p := range properties {
			if p.BasicInfo.Category == models.PropertyCategory(categoryFilter) {
				filtered = append(filtered, p)
			}
		}
		properties = filtered
	}

	start := offset
	end := offset + limit
	if end > len(properties) {
		end = len(properties)
	}
	paginationProperties := properties[start:end]

	path := "home"
	return c.Render(path, fiber.Map{
		"Title":      "Kömürcü Emlak - Anasayfa",
		"Properties": paginationProperties,
	}, "layouts/main")
}

func AboutWeb(c fiber.Ctx) error {
	path := "about"
	return c.Render(path, fiber.Map{
		"Title": "Hakkımızda",
	}, "layouts/main")
}

func ContactsWeb(c fiber.Ctx) error {
	path := "contacts"
	return c.Render(path, fiber.Map{
		"Title": "İletişim",
	}, "layouts/main")
}

func BlogSingleWeb(c fiber.Ctx) error {
	path := "blog-single"
	return c.Render(path, fiber.Map{
		"Title": "Tek Haberler",
	}, "layouts/main")
}
func BlogsWeb(c fiber.Ctx) error {
	path := "blogs"
	return c.Render(path, fiber.Map{
		"Title": "Haberler",
	}, "layouts/main")
}


func ListingSingle(c fiber.Ctx) error {
	propertyIDStr := c.Params("property_id") // URL'den property ID'yi al
	if propertyIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Geçersiz Property ID")
	}

	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Geçersiz Property ID formatı")
	}

	GetPropertyByID := func(ctx context.Context, propertyID uuid.UUID) (*models.Property, error) {
		row := database.DBPool.QueryRow(ctx, `
			SELECT
				p.property_id as property_id,
				bi.basic_info_id as basic_info_id,
				bi.property_type as property_type,
				bi.category as category,
				bi.main_title as main_title,
				bi.price as price,
				bi.keywords as keywords,
				loc.property_id as property_id,
				loc.address as adress,
				a.amenities_id as amenities_id,
				a.wifi as wifi,
				a.pool as pool,
				a.security as security,
				a.laundry_room as laundry_room,
				a.equipped_kitchen as equipped_kitchen,
				a.air_conditioning as air_conditioning,
				a.parking as parking_amenities,
				a.garage_atached as garage_atached,
				a.fireplace as fireplace,
				a.window_covering as window_covering,
				a.backyard as backyard,
				a.fitness_gym as fitness_gym,
				a.elevator as elevator,
				a.others_name as others_name,
				a.others_checked as others_checked,
				n.nearby_id as nearby_id,
				n.places as places, 
				n.distance as distance,
				pd.property_id as property_id,
				pd.property_message as property_message,
				pd.bedrooms as bedrooms,
				pd.bathrooms as bathrooms,
				pd.area as area,
				pd.parking as parking_details
			FROM
				property p
			LEFT JOIN
				basic_infos bi ON p.property_id = bi.property_id
			LEFT JOIN
				location loc ON p.property_id = loc.property_id
			LEFT JOIN 
				amenities a ON p.property_id = a.property_id
			LEFT JOIN
				nearby n ON p.property_id = n.property_id
			LEFT JOIN
				property_details pd ON p.property_id = pd.property_id
			WHERE p.property_id = $1
		`, propertyID)

		var property models.Property
		var basicInfos models.BasicInfo
		var location models.Location
		var amenities models.Amenities
		var nearby models.Nearby
		var propertyDetails models.PropertyDetails
		err := row.Scan(
			&property.PropertyID,
			&basicInfos.BasicInfoID,
			&basicInfos.Type,
			&basicInfos.Category,
			&basicInfos.MainTitle,
			&basicInfos.Price,
			&basicInfos.Keywords,
			&location.PropertyID,
			&location.Address,
			&amenities.AmenitiesID,
			&amenities.Wifi,
			&amenities.Pool,
			&amenities.Security,
			&amenities.LaundryRoom,
			&amenities.EquippedKitchen,
			&amenities.AirConditioning,
			&amenities.Parking,
			&amenities.GarageAtached,
			&amenities.Fireplace,
			&amenities.WindowCovering,
			&amenities.Backyard,
			&amenities.FitnessGym,
			&amenities.Elevator,
			&amenities.OthersName,
			&amenities.OthersChecked,
			&nearby.NearbyID,
			&nearby.Places,
			&nearby.Distance,
			&propertyDetails.PropertyID,
			&propertyDetails.PropertyMessage,
			&propertyDetails.Bedrooms,
			&propertyDetails.Bathrooms,
			&propertyDetails.Area,
			&propertyDetails.Parking,
		)
		if err != nil {
			fmt.Println("Sorgu hatası: ", err)
			return nil, err
		}

		property.BasicInfo = &basicInfos
		property.Location = &location
		property.Amenities = []*models.Amenities{&amenities}
		property.Nearby = []*models.Nearby{&nearby}
		property.PropertyDetails = &propertyDetails

		//Resimleri getiren fonksiyon
		GetImagesByPropertyID := func(ctx context.Context, propertyID uuid.UUID) ([]*models.Image, error) {
			rows, err := database.DBPool.Query(ctx, `
				SELECT image_id, property_id, name, file_path
				FROM images
				WHERE property_id = $1
			`, propertyID)
			if err != nil {
				fmt.Println("Resim sorgulama hatası: ", err)
				return nil, err
			}
			defer rows.Close()

			var images []*models.Image
			for rows.Next() {
				var image models.Image
				err := rows.Scan(&image.ImageID, &image.PropertyID, &image.ImageName, &image.FilePath)
				if err != nil {
					fmt.Println("Resim satırı tarama hatası: ", err)
					continue
				}
				images = append(images, &image)
			}

			if err := rows.Err(); err != nil {
				fmt.Println("Resim satırları yineleme hatası: ", err)
				return nil, err
			}

			return images, nil
		}
		// Resimleri getir
		images, err := GetImagesByPropertyID(ctx, property.PropertyID)
		if err != nil {
			fmt.Println("Resim getirme hatası: ", err)
			//Hata durumunda ne yapılacağına karar verin, örneğin boş bir dilim atayın
			property.PropertyMedia = []*models.PropertyMedia{}
		} else {
			// PropertyMedia'yı doldur
			propertyMedia := &models.PropertyMedia{
				PropertyID: property.PropertyID,
				Image:      images, // Resimleri doğrudan ata
			}
			property.PropertyMedia = []*models.PropertyMedia{propertyMedia} // Slice içinde sakla
		}
		return &property, nil
	}

	ctx := context.Background()
	property, err := GetPropertyByID(ctx, propertyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Veri alınırken hata oluştu.")
	}

	path := "listing-single" // Şablon dosyanızın adı
	return c.Render(path, fiber.Map{
		"Title":    property.BasicInfo.MainTitle, // Şablonunuza göre başlık
		"Property": property,                      // Tüm mülk bilgilerini şablona gönderiyoruz.
	}, "layouts/main")
}

func ListingWeb(c fiber.Ctx) error {
	categoryFilter := c.Query("category")
	offset, _ := strconv.Atoi(c.Query("offset","0"))
	limit, _ := strconv.Atoi(c.Query("limit", "12"))
	GetPropertiesByJoin := func(ctx context.Context) ([]*models.Property, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT
				p.property_id as property_id,
				bi.basic_info_id as basic_info_id,
				bi.property_type as property_type,
				bi.category as category,
				bi.main_title as main_title,
				bi.price as price,
				loc.property_id as property_id,
				loc.address as address,
				pd.property_id as property_id,
				pd.property_message as property_message,
				pd.bedrooms as bedrooms,
				pd.bathrooms as bathrooms,
				pd.area as area
			FROM
				property p
			LEFT JOIN
				basic_infos bi ON p.property_id = bi.property_id
			LEFT JOIN
				location loc ON p.property_id = loc.property_id
			LEFT JOIN
				property_details pd ON p.property_id = pd.property_id
		`)
		if err != nil {
			fmt.Println("Sorgu hatası: ", err)
			return nil, err
		}
		defer rows.Close()

		var properties []*models.Property

		for rows.Next() {
			var property models.Property
			var basicInfos models.BasicInfo
			var location models.Location
			var propertyDetails models.PropertyDetails

			err := rows.Scan(
				&property.PropertyID, &basicInfos.PropertyID, &basicInfos.Type, &basicInfos.Category, &basicInfos.MainTitle, &basicInfos.Price, &location.PropertyID, &location.Address, &propertyDetails.PropertyID, &propertyDetails.PropertyMessage, &propertyDetails.Bedrooms, &propertyDetails.Bathrooms, &propertyDetails.Area,
			)
			if err != nil {
				fmt.Println("Satır tarama hatası: ", err)
				continue // Hata durumunda sonraki satıra geç
			}

			property.BasicInfo = &basicInfos
			property.Location = &location
			property.PropertyDetails = &propertyDetails



			//Resimleri getiren fonksiyon
			GetImagesByPropertyID := func(ctx context.Context, propertyID uuid.UUID) ([]*models.Image, error) {
				rows, err := database.DBPool.Query(ctx, `
				SELECT image_id, property_id, name, file_path
				FROM images
				WHERE property_id = $1
			`, propertyID)
				if err != nil {
					fmt.Println("Resim sorgulama hatası: ", err)
					return nil, err
				}
				defer rows.Close()

				var images []*models.Image
				for rows.Next() {
					var image models.Image
					err := rows.Scan(&image.ImageID, &image.PropertyID, &image.ImageName, &image.FilePath)
					if err != nil {
						fmt.Println("Resim satırı tarama hatası: ", err)
						continue
					}
					images = append(images, &image)
				}

				if err := rows.Err(); err != nil {
					fmt.Println("Resim satırları yineleme hatası: ", err)
					return nil, err
				}

				return images, nil
			}
			// Resimleri getir
			images, err := GetImagesByPropertyID(ctx, property.PropertyID)
			if err != nil {
				fmt.Println("Resim getirme hatası: ", err)
				//Hata durumunda ne yapılacağına karar verin, örneğin boş bir dilim atayın
				property.PropertyMedia = []*models.PropertyMedia{}
			} else {
				// PropertyMedia'yı doldur
				propertyMedia := &models.PropertyMedia{
					PropertyID: property.PropertyID,
					Image:      images, // Resimleri doğrudan ata
				}
				property.PropertyMedia = []*models.PropertyMedia{propertyMedia} // Slice içinde sakla
			}

			properties = append(properties, &property)
		}
		return properties, nil
	}
	ctx := context.Background()
	properties, err := GetPropertiesByJoin(ctx)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).SendString("Verileri alırken hata oluştu")
	}

	if categoryFilter != "" {
		filtered := []*models.Property{}
		for _, p := range properties {
			if p.BasicInfo.Category == models.PropertyCategory(categoryFilter) {
				filtered = append(filtered, p)
			}
		}
		properties = filtered
	}




	start := offset
	end := offset + limit
	if end > len(properties){
		end = len(properties)
	}
	paginationProperties := properties[start:end]
	path := "listing"
	return c.Render(path, fiber.Map{
		"Title": "Daireler",
		"Properties": paginationProperties,
	}, "layouts/main")
}





func ProjectWeb(c fiber.Ctx) error {
	path := "projects"
	return c.Render(path, fiber.Map{
		"Title": "Projeler",
	}, "layouts/main")
}

func DashboardWeb(c fiber.Ctx) error {

	//user_ID := c.Params("user_id")
	
	path := "dashboard"
	return c.Render(path, fiber.Map{
		"Title": "Dashboard",
	}, "layouts/main")
}

func AddPropertyWeb(c fiber.Ctx) error{
	path := "add-property"
	return c.Render(path, fiber.Map{
		"Title": "Mülk Ekle",
	}, "layouts/main")
}