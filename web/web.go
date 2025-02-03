package web

import (
	"github.com/gofiber/fiber/v3"
	
)

func LoginWeb(c fiber.Ctx) error {
	path := "login"
	return c.Render(path, fiber.Map{
		"Title": "Login",
	})
}


func HomeWeb(c fiber.Ctx) error {
	path := "home"
	return c.Render(path, fiber.Map{
		"Title": "Kömürcü Emlak - Anasayfa",
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