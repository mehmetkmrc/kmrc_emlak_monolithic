package web

import (
	"context"
	"fmt"
	"kmrc_emlak_mono/database"
	"kmrc_emlak_mono/models"
	"strconv"

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
	offset, _ := strconv.Atoi(c.Query("offset","0"))
	limit, _ := strconv.Atoi(c.Query("limit", "6"))

	GetPropertiesByJoin := func(ctx context.Context) ([]*models.Property, error){
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
			if err != nil{
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


	path := "home"
	return c.Render(path, fiber.Map{
		"Title": "Kömürcü Emlak - Anasayfa",
		"Properties": paginationProperties,
	})
}

func AboutWeb(c fiber.Ctx) error {
	path := "about"
	return c.Render(path, fiber.Map{
		"Title": "Hakkımızda",
	})
}

func ContactsWeb(c fiber.Ctx) error {
	path := "contacts"
	return c.Render(path, fiber.Map{
		"Title": "İletişim",
	})
}

func BlogSingleWeb(c fiber.Ctx) error {
	path := "blog-single"
	return c.Render(path, fiber.Map{
		"Title": "Tek Haberler",
	})
}
func BlogsWeb(c fiber.Ctx) error {
	path := "blogs"
	return c.Render(path, fiber.Map{
		"Title": "Haberler",
	})
}
func ListingSingle(c fiber.Ctx) error {
	path := "listing-single"
	return c.Render(path, fiber.Map{
		"Title": "Daire",
	})
}

func ListingWeb(c fiber.Ctx) error {
	path := "listing"
	return c.Render(path, fiber.Map{
		"Title": "Daireler",
	})
}
func ProjectWeb(c fiber.Ctx) error {
	path := "projects"
	return c.Render(path, fiber.Map{
		"Title": "Projeler",
	})
}

func DashboardWeb(c fiber.Ctx) error {

	//user_ID := c.Params("user_id")
	
	path := "dashboard"
	return c.Render(path, fiber.Map{
		"Title": "Dashboard",
	})
}

func AddPropertyWeb(c fiber.Ctx) error{
	path := "add-property"
	return c.Render(path, fiber.Map{
		"Title": "Mülk Ekle",
	}, "layouts/main")
}