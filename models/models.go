package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Role string
)

type User struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"first_name"`
	Surname   string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone 	  string 	`json:"phone"`
	AboutText string	`json:"about_text"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSocialLinks struct {
	ID int `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	Facebook string `json:"facebook" db:"facebook"`
	Tiktok string `json:"tiktok" db:"tiktok"`
	Instagram string `json:"instagram" db:"instagram"`
	Twitter string `json:"twitter" db:"twitter"`
	Youtube string `json:"youtube" db:"youtube"`
	Linkedin string `json:"linkedin" db:"linkedin"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Session struct{
	SessionID string `json:"session_id"`
	UserID 		string `json:"user_id"`
	Token 		string `json:"token"`
	IPAdress	string `json:"ip_address"`
	UserAgent	string `json:"user_agent"`
	CreatedAt	time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	LastAccess time.Time `json:"last_access"`
	IsActive bool `json:"is_active"`
	Location string `json:"location"`
}

type Property struct {
	UserID uuid.UUID `json:"user_id"`
	PropertyID uuid.UUID `json:"property_id"`
	PropertyTitle string `json:"title"`
	TariffPlan string `json:"tariff_plan"`
	Date time.Time `json:"date"`
	PropertyDetails *PropertyDetails `json:"property_details"`
	VideoWidget []*VideoWidget `json:"video_widget"`
	Location *Location `json:"location"`
	Amenities []*Amenities `json:"amenities"`
	AccordionWidget []*AccordionWidget `json:"accordion_widget"`
	PropertyMedia []*PropertyMedia `json:"property_media"`
	BasicInfo *BasicInfo `json:"basic_info"`
	Nearby 	[]*Nearby `json:"nearby"`
	PlansBrochures []*PlansBrochures `json:"plans_brochures"`
}



type PropertyType string
type PropertyCategory string

const (
	TypeSale       PropertyType = "Sale"
	TypeRent       PropertyType = "Rent"
	TypeCommercial PropertyType = "Commercial"

	CategoryHouse     PropertyCategory = "House"
	CategoryApartment PropertyCategory = "Apartment"
	CategoryHotel     PropertyCategory = "Hotel"
	CategoryVilla     PropertyCategory = "Villa"
	CategoryOffice    PropertyCategory = "Office"
)
type BasicInfo struct {
	PropertyID uuid.UUID `json:"property_id"`
	BasicInfoID uuid.UUID `json:"basic_info_id"`
	MainTitle string `json:"main_title"`
	Type 	PropertyType `json:"property_type"`
	Category PropertyCategory `json:"category"`
	Price	float32 `json:"price"`
	Keywords string `json:"keywords"`
}


type CityLocation string

const (
	Adana        CityLocation = "Adana"
	Adıyaman     CityLocation = "Adıyaman"
	Afyonkarahisar CityLocation = "Afyonkarahisar"
	Ağrı         CityLocation = "Ağrı"
	Aksaray      CityLocation = "Aksaray"
	Amasya       CityLocation = "Amasya"
	Ankara       CityLocation = "Ankara"
	Antalya      CityLocation = "Antalya"
	Ardahan      CityLocation = "Ardahan"
	Artvin       CityLocation = "Artvin"
	Aydın        CityLocation = "Aydın"
	Balıkesir    CityLocation = "Balıkesir"
	Bartın       CityLocation = "Bartın"
	Batman       CityLocation = "Batman"
	Bayburt      CityLocation = "Bayburt"
	Bilecik      CityLocation = "Bilecik"
	Bingöl       CityLocation = "Bingöl"
	Bitlis       CityLocation = "Bitlis"
	Bolu         CityLocation = "Bolu"
	Burdur       CityLocation = "Burdur"
	Bursa        CityLocation = "Bursa"
	Çanakkale    CityLocation = "Çanakkale"
	Çankırı      CityLocation = "Çankırı"
	Çorum        CityLocation = "Çorum"
	Denizli      CityLocation = "Denizli"
	Diyarbakır   CityLocation = "Diyarbakır"
	Düzce        CityLocation = "Düzce"
	Edirne       CityLocation = "Edirne"
	Elazığ       CityLocation = "Elazığ"
	Erzincan     CityLocation = "Erzincan"
	Erzurum      CityLocation = "Erzurum"
	Eskişehir   CityLocation = "Eskişehir"
	Gaziantep    CityLocation = "Gaziantep"
	Giresun      CityLocation = "Giresun"
	Gümüşhane   CityLocation = "Gümüşhane"
	Hakkari      CityLocation = "Hakkari"
	Hatay        CityLocation = "Hatay"
	Iğdır        CityLocation = "Iğdır"
	Isparta      CityLocation = "Isparta"
	İstanbul     CityLocation = "İstanbul"
	İzmir        CityLocation = "İzmir"
	Kahramanmaraş CityLocation = "Kahramanmaraş"
	Karabük      CityLocation = "Karabük"
	Karaman      CityLocation = "Karaman"
	Kars         CityLocation = "Kars"
	Kastamonu    CityLocation = "Kastamonu"
	Kayseri      CityLocation = "Kayseri"
	Kırıkkale    CityLocation = "Kırıkkale"
	Kırklareli  CityLocation = "Kırklareli"
	Kırşehir     CityLocation = "Kırşehir"
	Kilis        CityLocation = "Kilis"
	Kocaeli      CityLocation = "Kocaeli"
	Konya        CityLocation = "Konya"
	Kütahya      CityLocation = "Kütahya"
	Malatya      CityLocation = "Malatya"
	Manisa       CityLocation = "Manisa"
	Mardin       CityLocation = "Mardin"
	Mersin       CityLocation = "Mersin"
	Muğla        CityLocation = "Muğla"
	Muş          CityLocation = "Muş"
	Nevşehir    CityLocation = "Nevşehir"
	Niğde        CityLocation = "Niğde"
	Ordu         CityLocation = "Ordu"
	Osmaniye    CityLocation = "Osmaniye"
	Rize         CityLocation = "Rize"
	Sakarya      CityLocation = "Sakarya"
	Samsun       CityLocation = "Samsun"
	Siirt        CityLocation = "Siirt"
	Sinop        CityLocation = "Sinop"
	Sivas        CityLocation = "Sivas"
	Şanlıurfa   CityLocation = "Şanlıurfa"
	Şırnak       CityLocation = "Şırnak"
	Tekirdağ    CityLocation = "Tekirdağ"
	Tokat        CityLocation = "Tokat"
	Trabzon      CityLocation = "Trabzon"
	Tunceli      CityLocation = "Tunceli"
	Uşak         CityLocation = "Uşak"
	Van          CityLocation = "Van"
	Yalova       CityLocation = "Yalova"
	Yozgat       CityLocation = "Yozgat"
	Zonguldak    CityLocation = "Zonguldak"
)
type Location struct{
	PropertyID uuid.UUID `json:"property_id"`
	LocationID uuid.UUID `json:"location_id"`
	Phone 	  int 	`json:"phone"`
	Email 		string `json:"email"`
	City 		CityLocation `json:"city"`
	Address 	string `json:"address"`
	Longitude 	float32 `json:"longitude"`
	Latitude 	float32 `json:"latitude"`
}

type PropertyNearby string


const (
	AllPlaces 	PropertyNearby = "Tüm Yerler"
	School 		PropertyNearby = "Okul"
	ShoppingMall PropertyNearby = "Alışveriş Merkezi"
	PoliceStation PropertyNearby = "Polis Karakolu"
	Hospital 	PropertyNearby 	 = "Hastane"
	PlaySchool PropertyNearby = "Oyun Okulu"
	Parks 	PropertyNearby = "Parklar"
)

type Nearby struct {
	PropertyID uuid.UUID `json:"property_id"`
	NearbyID uuid.UUID `json:"nearby_id"`
	Places 		PropertyNearby `json:"places"`
	Distance	int `json:"distance"`
}

type GalleryType string

const (
	GridGallery GalleryType = "Grid tipi"
	Carousel GalleryType 	= "Carousel"
	Slider GalleryType 		= "Slider"

)
type PropertyMedia struct {
	PropertyID uuid.UUID `json:"property_id"`
	PropertyMediaID uuid.UUID `json:"property_media_id"`
	ImageID uuid.UUID `json:"image_id"`
	Type GalleryType `json:"type"`
	Image []*Image  `json:"image"`
}

type Image struct{
	PropertyID 		uuid.UUID 	`json:"property_id"`
	ImageID  	  	uuid.UUID 	`json:"image_id"`
	ImageName     	[]string    `json:"name"`
	FilePath  	  	[]string    `json:"file_path"`
}

type PropertyDetails struct {
	PropertyID uuid.UUID `json:"property_id"`
	PropertyDetailsID uuid.UUID `json:"property_details_id"`
	Area 		float32 `json:"area"`
	Bedrooms	int 	`json:"bedrooms"`
	Bathrooms	int 	`json:"bathrooms"`
	Parking 	int 	`json:"parking"`
	Accomodation	string `json:"accomodation"`
	Website		string `json:"website"`
	PropertyMessage string `json:"property_message"`
}

type Amenities struct{
	PropertyID uuid.UUID `json:"property_id"`
	AmenitiesID uuid.UUID `json:"amenities_id"`
	Wifi 		bool	`json:"wifi"`
	Pool 		bool 	`json:"pool"`
	Security 	bool 	`json:"security"`
	LaundryRoom bool 	`json:"laundry_room"`
	EquippedKitchen bool `json:"equipped_kitchen"`
	AirConditioning bool `json:"air_conditioning"`
	Parking 		bool `json:"parking"`
	GarageAtached 	bool `json:"garage_atached"`
	Fireplace		bool `json:"fireplace"`
	WindowCovering  bool `json:"window_covering"`
	Backyard 		bool `json:"backyard"`
	FitnessGym 		bool `json:"fitness_gym"`
	Elevator		bool `json:"elevator"`
	OthersName			string `json:"others_name"`
	OthersChecked		bool `json:"others_checked"`
}

type AccordionWidget struct{
	PropertyID uuid.UUID `json:"property_id"`
	AccordionWidgetID uuid.UUID `json:"accordion_widget_id"`
	AccordionExist bool `json:"accordion_exist"`
	AccordionTitle string `json:"accordion_title"`
	AccordionDetails string `json:"accordion_details"`
}

type VideoWidget struct {
	PropertyID 	uuid.UUID `json:"property_id"`
	VideoWidgetID uuid.UUID `json:"video_widget_id"`
	VideoExist 	bool `json:"video_exist"`
	VideoTitle 	string `json:"video_title"`
	YouTubeUrl 	string `json:"youtube_url"`
	VimeoUrl	string `json:"vimeo_url"`

}
type PlansBrochures struct {
	PropertyID uuid.UUID `json:"property_id"`
	PlansBrochuresID uuid.UUID `json:"plans_brochures_id"`
	FileType string `json:"file_type"`
	FilePath string `json:"file_path"`
}