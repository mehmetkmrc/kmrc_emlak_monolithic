package dto


type PropertyType string
type PropertyCategory string

type CityLocation string

type PropertyNearby string

type GalleryType string



const (
	TypeSale       PropertyType = "Sale"
	TypeRent       PropertyType = "Rent"
	TypeCommercial PropertyType = "Commercial"

	CategoryHouse     PropertyCategory = "House"
	CategoryApartment PropertyCategory = "Apartment"
	CategoryHotel     PropertyCategory = "Hotel"
	CategoryVilla     PropertyCategory = "Villa"
	CategoryOffice    PropertyCategory = "Office"

	//Şehirler
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

	AllPlaces 	PropertyNearby = "Tüm Yerler"
	School 		PropertyNearby = "Okul"
	ShoppingMall PropertyNearby = "Alışveriş Merkezi"
	PoliceStation PropertyNearby = "Polis Karakolu"
	Hospital 	PropertyNearby 	 = "Hastane"
	PlaySchool PropertyNearby = "Oyun Okulu"
	Parks 	PropertyNearby = "Parklar"

	GridGallery GalleryType = "Grid tipi"
	Carousel GalleryType 	= "Carousel"
	Slider GalleryType 		= "Slider"
)

type (

	BasicInfoCreateRequest struct {
		MainTitle string `json:"main_title"`
		Type     PropertyType  `json:"property_type"`
		Category PropertyCategory  `json:"category"`
		Price    float32 `json:"price"`
		Keywords string  `json:"keywords"`
	}


	LocationCreateRequest struct {
		PropertyID string `json:"property_id"`
		Phone     string     `json:"phone"`
		Email     string  `json:"email"`
		City      CityLocation  `json:"city"`
		Address   string  `json:"address"`
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	}

	NearbyCreateRequest struct {
		PropertyID string `json:"property_id"`
		Places   PropertyNearby `json:"places"`
		Distance string    `json:"distance"`
	}

	PropertyMediaCreateRequest struct {
		PropertyID string `json:"property_id"`
		ImageID string `json:"image_id"`
		Type string `json:"type"`
	}

	ImageCreateRequest struct {
		PropertyID string `json:"property_id"`
		ImageName string `json:"name"`
		FilePath  string `json:"file_path"`
	}

	PropertyDetailsCreateRequest struct {
		PropertyID string `json:"property_id"`
		Area            string `json:"area"`
		Bedrooms        string     `json:"bedrooms"`
		Bathrooms       string     `json:"bathrooms"`
		Parking         string     `json:"parking"`
		Accomodation    string  `json:"accomodation"`
		Website         string  `json:"website"`
		PropertyMessage string  `json:"property_message"`
	}

	AmenitiesCreateRequest struct {
		PropertyID string `json:"property_id"`
		Wifi            bool   `json:"wifi"`
		Pool            bool   `json:"pool"`
		Security        bool   `json:"security"`
		LaundryRoom     bool   `json:"laundry_room"`
		EquippedKitchen bool   `json:"equipped_kitchen"`
		AirConditioning bool   `json:"air_conditioning"`
		Parking         bool   `json:"parking"`
		GarageAtached   bool   `json:"garage_atached"`
		Fireplace       bool   `json:"fireplace"`
		WindowCovering  bool   `json:"window_covering"`
		Backyard        bool   `json:"backyard"`
		FitnessGym      bool   `json:"fitness_gym"`
		Elevator        bool   `json:"elevator"`
		OthersName      string `json:"others_name"`
		OthersChecked   bool   `json:"others_checked"`
	}

	AccordionWidgetCreateRequest struct {
		PropertyID string `json:"property_id"`
		AccordionExist   bool   `json:"accordion_exist"`
		AccordionTitle   string `json:"accordion_title"`
		AccordionDetails string `json:"accordion_details"`
	}

	VideoWidgetCreateRequest struct {
		PropertyID string `json:"property_id"`
		VideoExist bool   `json:"video_exist"`
		VideoTitle string `json:"video_title"`
		YouTubeUrl string `json:"youtube_url"`
		VimeoUrl   string `json:"vimeo_url"`
	}

	PlansBrochuresCreateRequest struct {
		PropertyID string `json:"property_id"`
		FileType string `json:"file_type"`
		FilePath string `json:"file_path"`
	}

	MainPropertyCreateRequest struct {
		PropertyID string `json:"property_id"`
		PropertyTitle string `json:"title"`
	}	
)