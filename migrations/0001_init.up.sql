-- USERS TABLE
CREATE TABLE public.users (
    user_id uuid PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL UNIQUE,
    phone text,
    about_text text,
    password text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at time with time zone,
    role text,
    last_login timestamp without time zone,
    property_id uuid,
    photo_url text
);

-- PROPERTY TABLE
CREATE TABLE public.property (
    property_id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    tariff_plan text,
    date timestamp without time zone NOT NULL,
    title text NOT NULL,
    property_status integer NOT NULL DEFAULT 1,
    CONSTRAINT fk_property_user
        FOREIGN KEY (user_id) REFERENCES public.users(user_id)
        ON DELETE CASCADE
);

-- IMAGES TABLE
CREATE TABLE images (
    image_id UUID PRIMARY KEY,
    property_id UUID NOT NULL,
    url TEXT NOT NULL,
    original_name TEXT,
    media_type TEXT CHECK(media_type IN ('gallery','plan','brochure','profile')),
    created_at TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_images_property
        FOREIGN KEY (property_id)
        REFERENCES property(property_id)
        ON DELETE CASCADE
);

-- PROPERTY DETAILS
CREATE TABLE public.property_details (
    property_details_id uuid PRIMARY KEY,
    property_id uuid,
    area real,
    bedrooms integer,
    bathrooms integer,
    parking integer,
    accomodation text,
    website text,
    property_message text,
    CONSTRAINT fk_prop_details_property
        FOREIGN KEY (property_id) REFERENCES public.property(property_id)
        ON DELETE CASCADE
);

-- BASIC INFOS
CREATE TABLE public.basic_infos (
    basic_info_id uuid PRIMARY KEY,
    property_id uuid NOT NULL,
    property_type text,
    category text,
    price real,
    keywords text,
    main_title text,
    CONSTRAINT fk_basic_info_property
        FOREIGN KEY (property_id) REFERENCES public.property(property_id)
        ON DELETE CASCADE
);

-- AMENITIES
CREATE TABLE public.amenities (
    amenities_id uuid PRIMARY KEY,
    property_id uuid NOT NULL,
    wifi boolean,
    pool boolean,
    security boolean,
    laundry_room boolean,
    equipped_kitchen boolean,
    air_conditioning boolean,
    parking boolean,
    garage_atached boolean,
    fireplace boolean,
    window_covering boolean,
    backyard boolean,
    fitness_gym boolean,
    elevator boolean,
    others jsonb DEFAULT '[]'::jsonb,
    CONSTRAINT fk_amenities_property
        FOREIGN KEY (property_id) REFERENCES public.property(property_id)
        ON DELETE CASCADE
);

-- LOCATION
CREATE TABLE public.location (
    location_id uuid PRIMARY KEY,
    property_id uuid NOT NULL,
    phone bigint,
    email text,
    city text,
    address text,
    longitude real,
    latitude real,
    CONSTRAINT fk_location_property
        FOREIGN KEY (property_id) REFERENCES public.property(property_id)
        ON DELETE CASCADE
);

-- NEARBY
CREATE TABLE public.nearby (
    nearby_id uuid PRIMARY KEY,
    property_id uuid NOT NULL,
    places text,
    distance integer,
    CONSTRAINT fk_nearby_property
        FOREIGN KEY (property_id) REFERENCES public.property(property_id)
        ON DELETE CASCADE
);

-- PLANS & BROCHURES
CREATE TABLE plans_brochures (
    plans_brochures_id UUID PRIMARY KEY,
    property_id UUID NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_plans_property
        FOREIGN KEY (property_id)
        REFERENCES property(property_id)
        ON DELETE CASCADE
);

-- PROPERTY MEDIA
CREATE TABLE property_media (
    property_media_id UUID PRIMARY KEY,
    property_id UUID NOT NULL,
    image_id UUID NOT NULL,
    type TEXT CHECK(type IN ('gallery','plan','brochure')),

    CONSTRAINT fk_property_media_property
        FOREIGN KEY (property_id)
        REFERENCES property(property_id)
        ON DELETE CASCADE,

    CONSTRAINT fk_property_media_image
        FOREIGN KEY (image_id)
        REFERENCES images(image_id)
        ON DELETE CASCADE
);

-- VIDEO WIDGET
CREATE TABLE public.video_widget (
    video_widget_id uuid PRIMARY KEY,
    property_id uuid NOT NULL,
    video_exist boolean,
    video_title text,
    youtube_url text,
    vimeo_url text,
    CONSTRAINT fk_video_property
        FOREIGN KEY (property_id) REFERENCES public.property(property_id)
        ON DELETE CASCADE
);

-- ACCORDION WIDGET
CREATE TABLE public.accordion_widget (
    accordion_widget_id uuid PRIMARY KEY,
    property_id uuid NOT NULL,
    accordion_exist boolean,
    accordion_title text,
    accordion_details text,
    CONSTRAINT fk_accordion_property
        FOREIGN KEY (property_id) REFERENCES public.property(property_id)
        ON DELETE CASCADE
);

-- SESSIONS TABLE
CREATE TABLE public.sessions (
    session_id uuid PRIMARY KEY,
    user_id uuid,
    token character varying,
    ip_address character varying(45),
    user_agent text,
    created_at timestamp without time zone,
    expires_at timestamp without time zone,
    last_access timestamp without time zone,
    is_active boolean,
    location character varying,
    CONSTRAINT fk_sessions_user
        FOREIGN KEY (user_id) REFERENCES public.users(user_id)
        ON DELETE SET NULL
);
