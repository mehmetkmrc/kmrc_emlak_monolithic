package web

import (
	"context"
	"fmt"

	"kmrc_emlak_mono/database"
	"kmrc_emlak_mono/dto"
	"kmrc_emlak_mono/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	
	"go.uber.org/zap"
)



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
				p.property_id,
				bi.basic_info_id,
				bi.property_type,
				bi.category,
				bi.main_title,
				bi.price,
				loc.address,
				pd.property_message,
				pd.bedrooms,
				pd.bathrooms,
				pd.area
			FROM property p
			LEFT JOIN basic_infos bi ON p.property_id = bi.property_id
			LEFT JOIN location loc ON p.property_id = loc.property_id
			LEFT JOIN property_details pd ON p.property_id = pd.property_id
			ORDER BY p.property_id
			OFFSET $1 LIMIT $2
		`, offset, limit)
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
				&property.PropertyID,
				&basicInfos.BasicInfoID,
				&basicInfos.Type,
				&basicInfos.Category,
				&basicInfos.MainTitle,
				&basicInfos.Price,
				&location.Address,
				&propertyDetails.PropertyMessage,
				&propertyDetails.Bedrooms,
				&propertyDetails.Bathrooms,
				&propertyDetails.Area,
			)

			if err != nil {
				fmt.Println("Satır tarama hatası: ", err)
				continue // Hata durumunda sonraki satıra geç
			}
			property.BasicInfo = &basicInfos
			property.Location = &location
			property.PropertyDetails = &propertyDetails
			properties = append(properties, &property)
		}
		return properties, nil
	}


	GetImages := func(ctx context.Context) (map[uuid.UUID][]*models.Image, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT i.image_id, i.property_id, i.url, i.original_name, i.media_type
			FROM images i
			ORDER BY i.created_at
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		result := make(map[uuid.UUID][]*models.Image)

		for rows.Next() {
			var img models.Image
			if err := rows.Scan(
				&img.ImageID,
				&img.PropertyID,
				&img.Url,
				&img.OriginalName,
				&img.MediaType,
			); err != nil {
				continue
			}
			result[img.PropertyID] = append(result[img.PropertyID], &img)
		}

		return result, nil
	}

	

	ctx := context.Background()
	properties, err := GetPropertiesByJoin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Verileri alırken hata oluştu")
	}

	imagesMap, _ := GetImages(ctx)

	for _, p := range properties {
		if imgs, ok := imagesMap[p.PropertyID]; ok && len(imgs) > 0 {
			p.PropertyMedia = &models.PropertyMedia{
				PropertyID: p.PropertyID,
				Image: imgs[0], // 🔥 sadece ilk image
				Type: "gallery",
			}
		}
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

	

	path := "home"
	return c.Render(path, fiber.Map{
	"Title":      "Kömürcü Emlak - Anasayfa",
	"Properties": properties,
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
	propertyIDStr := c.Params("property_id")
	if propertyIDStr == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Geçersiz Property ID")
	}

	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Geçersiz Property ID formatı")
	}

	ctx := context.Background()

	// 1️⃣ PROPERTY ANA BİLGİLER
	row := database.DBPool.QueryRow(ctx, `
		SELECT
			p.property_id,

			bi.basic_info_id,
			bi.property_type,
			bi.category,
			bi.main_title,
			bi.price,
			bi.keywords,

			loc.address,
			loc.latitude,
			loc.longitude,

			a.amenities_id,
			a.wifi,
			a.pool,
			a.security,
			a.laundry_room,
			a.equipped_kitchen,
			a.air_conditioning,
			a.parking,
			a.garage_atached,
			a.fireplace,
			a.window_covering,
			a.backyard,
			a.fitness_gym,
			a.elevator,
			a.others_name,
			a.others_checked,

			pd.property_message,
			pd.bedrooms,
			pd.bathrooms,
			pd.area,
			pd.parking
		FROM property p
		LEFT JOIN basic_infos bi ON p.property_id = bi.property_id
		LEFT JOIN location loc ON p.property_id = loc.property_id
		LEFT JOIN amenities a ON p.property_id = a.property_id
		LEFT JOIN property_details pd ON p.property_id = pd.property_id
		WHERE p.property_id = $1
	`, propertyID)

	var property models.Property
	var bi models.BasicInfo
	var loc models.Location
	var am models.Amenities
	var pd models.PropertyDetails

	err = row.Scan(
		&property.PropertyID,

		&bi.BasicInfoID,
		&bi.Type,
		&bi.Category,
		&bi.MainTitle,
		&bi.Price,
		&bi.Keywords,

		&loc.Address,
		&loc.Latitude,
		&loc.Longitude,

		&am.AmenitiesID,
		&am.Wifi,
		&am.Pool,
		&am.Security,
		&am.LaundryRoom,
		&am.EquippedKitchen,
		&am.AirConditioning,
		&am.Parking,
		&am.GarageAtached,
		&am.Fireplace,
		&am.WindowCovering,
		&am.Backyard,
		&am.FitnessGym,
		&am.Elevator,
		&am.OthersName,
		&am.OthersChecked,

		&pd.PropertyMessage,
		&pd.Bedrooms,
		&pd.Bathrooms,
		&pd.Area,
		&pd.Parking,
	)
	if err != nil {
		return c.Status(500).SendString("Property bulunamadı")
	}

	property.BasicInfo = &bi
	property.Location = &loc
	property.Amenities = []*models.Amenities{&am}
	property.PropertyDetails = &pd

	// 2️⃣ NEARBY
	nearbyRows, err := database.DBPool.Query(ctx, `
		SELECT nearby_id, property_id, places, distance
		FROM nearby
		WHERE property_id = $1
	`, propertyID)

	if err == nil {
		defer nearbyRows.Close()
		for nearbyRows.Next() {
			var n models.Nearby
			nearbyRows.Scan(&n.NearbyID, &n.PropertyID, &n.Places, &n.Distance)
			property.Nearby = append(property.Nearby, &n)
		}
	}

	// 3️⃣ PROPERTY MEDIA + IMAGE
mediaRows, err := database.DBPool.Query(ctx, `
	SELECT
		pm.property_media_id,
		pm.property_id,
		pm.image_id,
		pm.type,
		i.image_id,
		i.property_id,
		i.url,
		i.original_name,
		i.media_type,
		i.created_at
	FROM property_media pm
	JOIN images i ON pm.image_id = i.image_id
	WHERE pm.property_id = $1
	ORDER BY i.created_at
`, propertyID)

if err == nil {
	defer mediaRows.Close()

	for mediaRows.Next() {
	var pm models.PropertyMedia
	var img models.Image

	err := mediaRows.Scan(
		&pm.PropertyMediaID,
		&pm.PropertyID,
		&pm.ImageID,
		&pm.Type,
		&img.ImageID,
		&img.PropertyID,
		&img.Url,
		&img.OriginalName,
		&img.MediaType,
		&img.CreatedAt,
	)
	if err != nil {
		continue
	}

	pm.Image = &img
	property.PropertyMediaList = append(property.PropertyMediaList, &pm)
}

}


	// 4️⃣ RENDER
	return c.Render("ilan", fiber.Map{
		"Title":    property.BasicInfo.MainTitle,
		"Property": &property,
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
			properties = append(properties, &property)
		}
		return properties, nil
	}
	GetImages := func(ctx context.Context) (map[uuid.UUID][]*models.Image, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT i.image_id, i.property_id, i.url, i.original_name, i.media_type
			FROM images i
			ORDER BY i.created_at
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		result := make(map[uuid.UUID][]*models.Image)

		for rows.Next() {
			var img models.Image
			if err := rows.Scan(
				&img.ImageID,
				&img.PropertyID,
				&img.Url,
				&img.OriginalName,
				&img.MediaType,
			); err != nil {
				continue
			}
			result[img.PropertyID] = append(result[img.PropertyID], &img)
		}

		return result, nil
	}

	ctx := context.Background()
	properties, err := GetPropertiesByJoin(ctx)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).SendString("Verileri alırken hata oluştu")
	}

	imagesMap, _ := GetImages(ctx)

	for _, p := range properties {
		if imgs, ok := imagesMap[p.PropertyID]; ok && len(imgs) > 0 {
			p.PropertyMedia = &models.PropertyMedia{
				PropertyID: p.PropertyID,
				Image: imgs[0], // 🔥 sadece ilk image
				Type: "gallery",
			}
		}
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

	 user := c.Locals("UserDetail")
    if user == nil {
        return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
    }

    userInfo := user.(*dto.GetUserResponse)
	
	path := "kullanici-panel"
	return c.Render(path, fiber.Map{
		"Title": "Kullanıcı Paneli",
		"User": userInfo,
	}, "layouts/main")
}

func AddPropertyWeb(c fiber.Ctx) error{
	path := "yeni-ilan-ekle"
	return c.Render(path, fiber.Map{
		"Title": "Mülk Ekle",
	}, "layouts/main")
}

func EditProfile(c fiber.Ctx) error {
    // 1) User login kontrolü
    user := c.Locals("UserDetail")
    if user == nil {
        return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
    }

    userData := user.(*dto.GetUserResponse)
    userID := userData.UserID

    // 2) USERS tablosundan bilgiler
    query := `
        SELECT 
            first_name,
            last_name,
            email,
            phone,
			photo_url,
            about_text
        FROM users
        WHERE user_id = $1
    `

    var profile models.User

    row := database.DBPool.QueryRow(c.Context(), query, userID)

    err := row.Scan(
        &profile.Name,
        &profile.Surname,
        &profile.Email,
        &profile.Phone,
		&profile.PhotoUrl,
        &profile.AboutText,
    )

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("User data not found")
    }

    profile.UserID = userID

    // -----------------------------------------
    // 3) user_social_links tablosundan sosyal linkleri çek
    // -----------------------------------------
    social := &models.UserSocialLinks{UserID: userID}

    socialQuery := `
        SELECT facebook, tiktok, instagram, twitter, youtube, linkedin
        FROM user_social_links
        WHERE user_id = $1
        LIMIT 1
    `

    socialRow := database.DBPool.QueryRow(c.Context(), socialQuery, userID)

    err = socialRow.Scan(
        &social.Facebook,
        &social.Tiktok,
        &social.Instagram,
        &social.Twitter,
        &social.Youtube,
        &social.Linkedin,
    )

    if err != nil {
        // Kayıt yoksa sorun değil
        zap.S().Warn("No social links found for user: ", userID)
    }

    // -----------------------------------------
    // 4) Render to View
    // -----------------------------------------
    return c.Render("profili-duzenle", fiber.Map{
        "Title":  "Profili Düzenle",
        "User":   profile,
        "Social": social,
    }, "layouts/main")
}

func ListingMyProperties(c fiber.Ctx) error {
	
    // Kullanıcı bilgilerini Locals'tan al
    user := c.Locals("UserDetail")
    if user == nil {
        return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
    }
    userData := user.(*dto.GetUserResponse)
    userID := userData.UserID

	//User ı çek
	

	GetPropertiesByJoin := func(ctx context.Context) ([]*models.Property, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT
				u.user_id as user_id,
				p.user_id as user_id,
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
				users u ON p.user_id = u.user_id
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
			//var user models.User
			var property models.Property
			var basicInfos models.BasicInfo
			var location models.Location
			var propertyDetails models.PropertyDetails

			err := rows.Scan(
				&userID, &property.UserID, &property.PropertyID, &basicInfos.PropertyID, &basicInfos.Type, &basicInfos.Category, &basicInfos.MainTitle, &basicInfos.Price, &location.PropertyID, &location.Address, &propertyDetails.PropertyID, &propertyDetails.PropertyMessage, &propertyDetails.Bedrooms, &propertyDetails.Bathrooms, &propertyDetails.Area,
			)
			if err != nil {
				fmt.Println("Satır tarama hatası: ", err)
				continue // Hata durumunda sonraki satıra geç
			}

			property.BasicInfo = &basicInfos
			property.Location = &location
			property.PropertyDetails = &propertyDetails
			properties = append(properties, &property)
		}
		return properties, nil
	}


	GetImages := func(ctx context.Context) (map[uuid.UUID][]*models.Image, error) {
		rows, err := database.DBPool.Query(ctx, `
			SELECT i.image_id, i.property_id, i.url, i.original_name, i.media_type
			FROM images i
			ORDER BY i.created_at
		`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		result := make(map[uuid.UUID][]*models.Image)

		for rows.Next() {
			var img models.Image
			if err := rows.Scan(
				&img.ImageID,
				&img.PropertyID,
				&img.Url,
				&img.OriginalName,
				&img.MediaType,
			); err != nil {
				continue
			}
			result[img.PropertyID] = append(result[img.PropertyID], &img)
		}

		return result, nil
	}


	ctx := context.Background()
	properties, err := GetPropertiesByJoin(ctx)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).SendString("Verileri alırken hata oluştu")
	}
    
	imagesMap, _ := GetImages(ctx)

	for _, p := range properties {
		if imgs, ok := imagesMap[p.PropertyID]; ok && len(imgs) > 0 {
			p.PropertyMedia = &models.PropertyMedia{
				PropertyID: p.PropertyID,
				Image: imgs[0], // 🔥 sadece ilk image
				Type: "gallery",
			}
		}
	}
	
    return c.Render("ilanlarım", fiber.Map{
        "Title":      "İlanlarım",
        "Properties": properties,
    }, "layouts/main")
}

func EditPropertyWeb(c fiber.Ctx) error{

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
				p.property_id as p_property_id,
				bi.basic_info_id as basic_info_id,
				bi.property_type as property_type,
				bi.category as category,
				bi.main_title as main_title,
				bi.price as price,
				bi.keywords as keywords,

				loc.property_id as loc_property_id,
				loc.location_id as location_id,
				loc.phone as phone,
				loc.email as email,
				loc.city as city,
				loc.address as adress,
				loc.latitude as latitude,
				loc.longitude as longitude,

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

				pd.property_id as pd_property_id,
				pd.property_details_id as pd_property_detail_id,
				pd.property_message as property_message,
				pd.bedrooms as bedrooms,
				pd.bathrooms as bathrooms,
				pd.website as website,
				pd.area as area,
				pd.accomodation as accomodation,
				pd.parking as parking_details,

				aw.accordion_widget_id as aw_id,
				aw.property_id as aw_property_id,
				aw.accordion_exist as accordion_exist,
				aw.accordion_title as accordion_title,
				aw.accordion_details as accordion_details,

				vw.video_widget_id as vw_id,
				vw.property_id as vw_propertyid,
				vw.video_exist as video_exist,
				vw.video_title as video_title,
				vw.youtube_url as youtube_url,
				vw.vimeo_url as vimeo_url

			FROM property p
			LEFT JOIN basic_infos bi ON p.property_id = bi.property_id
			LEFT JOIN location loc ON p.property_id = loc.property_id
			LEFT JOIN amenities a ON p.property_id = a.property_id
			LEFT JOIN nearby n ON p.property_id = n.property_id
			LEFT JOIN property_details pd ON p.property_id = pd.property_id
			LEFT JOIN accordion_widget aw ON p.property_id = aw.property_id
			LEFT JOIN video_widget vw ON p.property_id = vw.property_id
			WHERE p.property_id = $1
		`, propertyID)

		var property models.Property
		var basicInfos models.BasicInfo
		var location models.Location
		var amenities models.Amenities
		var nearby models.Nearby
		var propertyDetails models.PropertyDetails
		var accordionWidget models.AccordionWidget
		var videoWidget models.VideoWidget
		err := row.Scan(
			&property.PropertyID,      // p_property_id
			&basicInfos.BasicInfoID,
			&basicInfos.Type,
			&basicInfos.Category,
			&basicInfos.MainTitle,
			&basicInfos.Price,
			&basicInfos.Keywords,

			&location.PropertyID,      // loc_property_id
			&location.LocationID,
			&location.Phone,
			&location.Email,
			&location.City,
			&location.Address,
			&location.Latitude,
			&location.Longitude,

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

			&propertyDetails.PropertyID, // pd_property_id
			&propertyDetails.PropertyDetailsID,
			&propertyDetails.PropertyMessage,
			&propertyDetails.Bedrooms,
			&propertyDetails.Bathrooms,
			&propertyDetails.Website,
			&propertyDetails.Area,
			&propertyDetails.Accomodation,
			&propertyDetails.Parking,

			&accordionWidget.AccordionWidgetID, // aw_id
			&accordionWidget.PropertyID,        // aw_property_id
			&accordionWidget.AccordionExist,
			&accordionWidget.AccordionTitle,
			&accordionWidget.AccordionDetails,

			&videoWidget.VideoWidgetID,
			&videoWidget.PropertyID,
			&videoWidget.VideoExist,
			&videoWidget.VideoTitle,
			&videoWidget.YouTubeUrl,
			&videoWidget.VimeoUrl,
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
		property.AccordionWidget = []*models.AccordionWidget{&accordionWidget}
		property.VideoWidget = []*models.VideoWidget{&videoWidget}
		GetNearbyByPropertyID := func(ctx context.Context, propertyID uuid.UUID) ([]*models.Nearby, error) {
			rows, err := database.DBPool.Query(ctx, `
				SELECT nearby_id, property_id, places, distance
				FROM nearby
				WHERE property_id = $1
			`, propertyID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			var nearbyList []*models.Nearby
			for rows.Next() {
				var n models.Nearby
				err := rows.Scan(&n.NearbyID, &n.PropertyID, &n.Places, &n.Distance)
				if err != nil {
					continue
				}
				nearbyList = append(nearbyList, &n)
			}

			return nearbyList, nil
		}

		nearbyList, err := GetNearbyByPropertyID(ctx, property.PropertyID)
		if err != nil {
			property.Nearby = []*models.Nearby{}
		} else {
			property.Nearby = nearbyList
		}

		return &property, nil
	}

	ctx := context.Background()

	


	property, err := GetPropertyByID(ctx, propertyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Veri alınırken hata oluştu.")
	}

	mediaRows, err := database.DBPool.Query(ctx, `
	SELECT
		pm.property_media_id,
		pm.property_id,
		pm.image_id,
		pm.type,
		i.image_id,
		i.property_id,
		i.url,
		i.original_name,
		i.media_type,
		i.created_at
	FROM property_media pm
	JOIN images i ON pm.image_id = i.image_id
	WHERE pm.property_id = $1
	ORDER BY i.created_at
`, propertyID)

if err == nil {
	defer mediaRows.Close()

	for mediaRows.Next() {
	var pm models.PropertyMedia
	var img models.Image

	err := mediaRows.Scan(
		&pm.PropertyMediaID,
		&pm.PropertyID,
		&pm.ImageID,
		&pm.Type,
		&img.ImageID,
		&img.PropertyID,
		&img.Url,
		&img.OriginalName,
		&img.MediaType,
		&img.CreatedAt,
	)
	if err != nil {
		continue
	}

	pm.Image = &img
	property.PropertyMediaList = append(property.PropertyMediaList, &pm)
}

}
	path := "ilan-duzenle"
	return c.Render(path, fiber.Map{
		"Title":    property.BasicInfo.MainTitle, // Şablonunuza göre başlık
		"Property": &property,                      // Tüm mülk bilgilerini şablona gönderiyoruz.
	}, "layouts/main")
}
