package main

import (
	"encoding/json"
	"errors"
	"html/template"

	"kmrc_emlak_mono/auth"
	"kmrc_emlak_mono/database"
	"kmrc_emlak_mono/middleware"
	"kmrc_emlak_mono/property"

	"kmrc_emlak_mono/web"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v2"
)

const (
	viewPath   = "./client/templates"
	publicPath = "./client/public"
	renderType = ".html"
	uploadsPath = "./"
)

func add(x, y int) int {
	return x + y
}

func main() {
	database.InitiliazeDatabaseConnection()
	engine := html.New(viewPath, renderType)
	engine.Reload(true)
	engine.AddFunc("unescape", func(s string) template.HTML {
		return template.HTML(s)
	})

	engine.AddFunc("safe", func(s string) template.HTML {
		return template.HTML(s) // HTML olarak işaretler, güvenli kabul eder
	})

	engine.AddFunc("attr", func(s string) template.HTMLAttr {
		return template.HTMLAttr(s) // Attribute olarak işaretler
	})
	engine.AddFunc("safeHTML", func(s string) template.HTML {
		return template.HTML(s) // HTML olarak işaretle
	})
	engine.AddFunc("raw", func(s string) template.HTML {
		return template.HTML(s) // Mark string as raw HTML
	})
	engine.AddFunc("add", add)
	app := fiber.New(fiber.Config{
		ReadTimeout:   time.Minute * time.Duration(5),
		StrictRouting: false,
		CaseSensitive: true,
		BodyLimit:     4 * 1024 * 1024,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		AppName:       "kmrcemlak",
		Immutable:     true,
		Views:         engine,
		//ViewsLayout: "layouts/main",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) {
				if e.Code == fiber.StatusNotFound {
					return c.Render("404", fiber.Map{
						"Title": "Page Not Found",
					})
				}
				return c.Status(e.Code).Render("error", fiber.Map{
					"Title":   "Error",
					"Message": e.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
				"Title":   "Internal Server Error",
				"Message": err,
			})
		},
	})
	app.Use(static.New(uploadsPath))
	app.Use(static.New(publicPath))
	


	app.Get("/", web.HomeWeb)
	app.Get("/about", web.AboutWeb)
	app.Get("/contacts", web.ContactsWeb)
	app.Get("/blog-single", web.BlogSingleWeb)
	app.Get("/blogs", web.BlogsWeb)
	app.Get("/ilan/:property_id", web.ListingSingle)
	app.Get("/listing", web.ListingWeb)
	app.Get("/projects", web.ProjectWeb)
	route0 := app.Group("/dashboard")
	route0.Get("/", web.DashboardWeb, auth.IsAuthorized, auth.GetUserDetail,   auth.RateLimiter(120, time.Minute))
	route0.Get("/add-new-property", web.AddPropertyWeb, auth.IsAuthorized, auth.GetUserDetail,  auth.RateLimiter(120, time.Minute))
	app.Post("/logout", auth.IsAuthorized, auth.Logout)
	//app.Get("/login", web.LoginPage, auth.RateLimiter(5, time.Minute))

	route := app.Group("/auth")
	route.Post("/login", auth.Login, auth.RateLimiter(5, time.Minute), auth.LoginValidation)
	route.Post("/register", auth.Register, auth.RateLimiter(5, time.Minute), auth.RegisterValidation)
	

	propertier := app.Group("/property")
	//propertier.Post("/add-property", property.AddProperty)
	propertier.Post("/add-property-details", property.AddPropertyDetails)
	propertier.Post("/add-video-widget", property.AddVideoWidget)
	propertier.Post("/add-location", property.AddLocation)
	propertier.Post("/add-amenities", property.AddAmenities)
	propertier.Post("/add-property-media", property.AddPropertyMedia)
	propertier.Post("/add-image", property.AddImage)
	propertier.Post("/add-basic-info", property.AddBasicInfo, auth.IsAuthorized,middleware.PropertyMiddleware, property.AddProperty,   auth.GetUserDetail)
	propertier.Post("/add-nearby", property.AddNearby)
	propertier.Post("/add-accordion-widget", property.AddAccordionWidget)
	propertier.Post("/add-plans-brochures", property.AddPlansBrochures)
	
	


	// document := app.Group("/documenter")
	// document.Post("/main", property.CreateMainDocument)
	// document.Post("/sub", property.CreateSubDocument)
	// document.Post("/content", property.CreateContentDocument)
	// document.Get("/all", property.GetAllDocuments)
	// document.Get("/all-join", property.GetAllDocumentsByJoin)

	
	
	// app.Use(web.NotFoundPage)
	
	//s.app.Get("/dashboard", s.DashboardWeb, s.authMiddleware)

	log.Fatal(app.Listen("0.0.0.0:8081"))
}
