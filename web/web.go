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
	path := "listing-single"
	return c.Render(path, fiber.Map{
		"Title": "Daire",
	}, "layouts/main")
}

func ListingWeb(c fiber.Ctx) error {
	path := "listing"
	return c.Render(path, fiber.Map{
		"Title": "Daireler",
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
