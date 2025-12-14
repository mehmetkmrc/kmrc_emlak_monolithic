

document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("documentLoader");

    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const loaderWrap = document.querySelector(".loader-wrap"); // .loader-wrap'i hedef al

        // Loader elementinin varlığını kontrol et
        if (!loaderWrap) {
            console.error("Loader elementi bulunamadı!");
            showModal("error", "Hata!", "Yükleniyor animasyonu bulunamadı!");
            return; // Form gönderme işlemini durdur
        }

        // Loader'ı görünür yap
        loaderWrap.style.display = "block"; // veya "flex", "grid" gibi uygun bir değer
        const loader = document.getElementById("loader");
        loader.classList.remove("d-none");

        try {
            // **1. Adım: BasicInfo oluştur**
            const main_title = document.querySelector('input[name="main_title"]').value;
            const type = document.querySelector('select[name="property_type"]').value;
            const category = document.querySelector('select[name="category"]').value;
            const price = document.querySelector('input[name="price"]').value;
            const keywords = document.querySelector('input[name="keywords"]').value;

            const basicInfoData = {
                main_title : main_title,
                property_type: type,
                category: category,
                price: parseFloat(price),
                keywords: keywords,
            };

            const basicInfoResponse = await fetch("http://127.0.0.1:8081/property/add-basic-info", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(basicInfoData),
            });

            if (!basicInfoResponse.ok) {
                const errorText = await basicInfoResponse.text();
                showModal("error", "Hata!", `Basic Info oluşturulurken bir hata oluştu! Hata: ${errorText}`);
                return;
            }

            const basicInfoResult = await basicInfoResponse.json(); // Yanıtı JSON olarak ayrıştır
            console.log("Basic Info Response:", basicInfoResult);

            if (basicInfoResult.status !== 200) {
                showModal("error", "Hata!", "Basic Info oluşturulamadı: " + basicInfoResult.message);
                return;
            }

            const propertyID = basicInfoResult.data; // **Backend'den gelen propertyID'yi alın.**
            //console.log("Property ID:", property_id);

            // **2. Adım: Location oluştur**
            const phone = document.querySelector('input[name="phone"]').value;
            const email = document.querySelector('input[name="email"]').value;
            const city = document.querySelector('select[name="city"]').value;
            const address = document.querySelector('input[name="address"]').value;
            const longitudeInput = document.querySelector('input[name="longitude"]');
            const longitudeValue = longitudeInput.value;
            let longitude = null; // Başlangıçta null olarak tanımla

            if (longitudeValue) {
                const parsedLongitude = parseFloat(longitudeValue);
                if (!isNaN(parsedLongitude)) { // Geçerli bir sayı mı?
                    longitude = parsedLongitude; // Geçerli ise değeri ata
                }
            }
            const latitudeInput = document.querySelector('input[name="latitude"]');
            const latitudeValue = latitudeInput.value;
            let latitude = null; // Başlangıçta null olarak tanımla

            if (latitudeValue) {
                const parsedLatitude = parseFloat(latitudeValue);
                if (!isNaN(parsedLatitude)) { // Geçerli bir sayı mı?
                    latitude = parsedLatitude; // Geçerli ise değeri ata
                }
            }

            const locationData = {
                property_id: propertyID.property_id, // **PropertyID'yi kullanın**
                phone: phone,
                email: email,
                city: city,
                address: address,
                latitude: latitudeInput.value.toString(),
                longitude: longitudeInput.value.toString(),
            };

            const locationResponse = await fetch("http://127.0.0.1:8081/property/add-location", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(locationData),
            });

            if (!locationResponse.ok) {
                const errorText = await locationResponse.text();
                showModal("error", "Hata!", `Location oluşturulurken bir hata oluştu! Hata: ${errorText}`);
                return;
            }

            const locationResult = await locationResponse.json();
            console.log("Location Response:", locationResult);

            if (locationResult.status !== 200) {
                showModal("error", "Hata!", "Location oluşturulamadı: " + locationResult.message);
                return;
            }

            // **3. Adım: Nearby oluştur**
            if (globalNearbyArray.length > 0) {

            try {
                for (let i = 0; i < globalNearbyArray.length; i++) {

                    const nearbyData = {
                        property_id: propertyID.property_id,
                        places: globalNearbyArray[i].places,
                        distance: globalNearbyArray[i].distance,
                    };

                    const nearbyResponse = await fetch("http://127.0.0.1:8081/property/add-nearby", {
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify(nearbyData),
                    });

                    if (!nearbyResponse.ok) {
                        const errorText = await nearbyResponse.text();
                        showModal("error", "Hata!", `Nearby oluşturulurken hata: ${errorText}`);
                        return;
                    }

                    const nearbyResult = await nearbyResponse.json();

                    if (nearbyResult.status !== 200) {
                        showModal("error", "Hata!", "Nearby oluşturulamadı: " + nearbyResult.message);
                        return;
                    }
                }

                } catch (error) {
                    console.error("Nearby API hatası:", error);
                    showModal("error", "Hata!", `Nearby oluşturulurken hata: ${error.message}`);
                }
        }




           
           // **4. Property Media oluştur**
           const galleryType = document.querySelector('select[name="type"]').value;
           await uploadImages(propertyID, imageState.gallery, "gallery");
           await uploadImages(propertyID, imageState.plan, "plan");


           async function uploadImages(propertyID, files, type) {
                if (!files.length) return true;

                const fd = new FormData();
                fd.append("property_id", propertyID.property_id);
                fd.append("type", type);

                files.forEach(file => fd.append("image", file));

                const res = await fetch("http://127.0.0.1:8081/property/add-image", {
                    method: "POST",
                    body: fd
                });

                if (!res.ok) {
                    const txt = await res.text();
                    showModal("error", "Hata!", txt);
                    return false;
                }

                const result = await res.json();
                return result.status === 200;
            }




            // **6. Property Details oluştur
            const area = document.querySelector('input[name="area"]').value;
            const bedrooms = document.querySelector('input[name="bedrooms"]').value;
            const bathrooms = document.querySelector('input[name="bathrooms"]').value;
            const parking = document.querySelector('input[name="parking"]').value;
            const accomodation = document.querySelector('input[name="accomodation"]').value;
            const website = document.querySelector('input[name="website"]').value;
            const property_message = document.querySelector('textarea[name="property_message"]').value;

            const propertyDetailsData ={
                property_id: propertyID.property_id,
                area : area,
                bedrooms: bedrooms,
                bathrooms: bathrooms,
                parking: parking,
                accomodation: accomodation,
                website: website,
                property_message: property_message,
            }

            const propertyDetailsResponse = await fetch("http://127.0.0.1:8081/property/add-property-details", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(propertyDetailsData)
            })

            if(!propertyDetailsResponse.ok){
                const errorText = await propertyDetailsResponse.text();
                showModal("error", "Hata!", `Property Detail oluşturulurken bir hata oluştu! Hata: ${errorText}`);
                return;
            }

            const propertyDetailsResult = await propertyDetailsResponse.json();
            console.log("Property Details Response:", propertyDetailsResult );

            if (propertyDetailsResult.status !== 200) {
                showModal("error", "Hata!", "Property Details oluşturulamadı: " + propertyDetailsResult.message);
                return;
            }

            let otherAmenitiesArray = [];

            function setupOtherAmenities() {
                const container = document.getElementById("other-amenities-container");

                function addNewOtherAmenity() {
                    const div = document.createElement("div");
                    div.classList.add("other-amenity");
                    div.style.display = "flex";
                    div.style.alignItems = "center";
                    div.style.gap = "10px";

                    const input = document.createElement("input");
                    input.type = "text";
                    input.placeholder = "Diğer Özellikleri buraya yazın...";
                    input.style.flex = "1";
                    input.style.padding = "8px 12px";
                    input.style.borderRadius = "5px";
                    input.style.border = "1px solid #ccc";

                    const checkbox = document.createElement("input");
                    checkbox.type = "checkbox";
                    checkbox.classList.add("other-amenity-checkbox");
                    checkbox.style.width = "20px";
                    checkbox.style.height = "20px";

                    const label = document.createElement("label");
                    label.style.margin = "0";
                    label.style.fontWeight = "500";
                    label.textContent = "Var";

                    div.appendChild(input);
                    div.appendChild(checkbox);
                    div.appendChild(label);

                    container.appendChild(div);

                    // Checkbox event
                    checkbox.addEventListener("change", function () {
                        if (this.checked && input.value.trim() !== "") {
                            otherAmenitiesArray.push(input.value.trim());
                            this.disabled = true;
                            input.disabled = true;

                            addNewOtherAmenity();
                        }
                    });
                }

                // İlk input event’i
                const firstCheckbox = container.querySelector(".other-amenity-checkbox");
                const firstInput = container.querySelector("input[type='text']");

                firstCheckbox.addEventListener("change", function () {
                    if (this.checked && firstInput.value.trim() !== "") {
                        otherAmenitiesArray.push(firstInput.value.trim());
                        this.disabled = true;
                        firstInput.disabled = true;

                        addNewOtherAmenity();
                    }
                });
            }

            document.addEventListener("DOMContentLoaded", setupOtherAmenities);


            // **7. Amenities oluştur
            const wifi = document.querySelector('input[name="wifi"]').checked;
            const pool = document.querySelector('input[name="pool"]').checked;
            const security = document.querySelector('input[name="security"]').checked;
            const laundryRoom = document.querySelector('input[name="laundry_room"]').checked;
            const equippedKitchen = document.querySelector('input[name="equipped_kitchen"]').checked;
            const airConditioning = document.querySelector('input[name="air_conditioning"]').checked;
            const parking1 = document.querySelector('input[name="parking_"]').checked;
            const garageAtached = document.querySelector('input[name="garage_atached"]').checked;
            const fireplace = document.querySelector('input[name="fireplace"]').checked;
            const windowCovering = document.querySelector('input[name="window_covering"]').checked;
            const backyard = document.querySelector('input[name="backyard"]').checked;
            const fitnessGym = document.querySelector('input[name="fitness_gym"]').checked;
            const elevator = document.querySelector('input[name="elevator"]').checked;
            // Yeni alanlar
           // Diğer özellikleri gönder
            const othersArray = otherAmenitiesArray; // Dinamik olarak eklenen tüm değerler


            const amenitiesData = {
                property_id: propertyID.property_id, // Daha önce aldığınız propertyID
                wifi: wifi,
                pool: pool,
                security: security,
                laundry_room: laundryRoom,
                equipped_kitchen: equippedKitchen,
                air_conditioning: airConditioning,
                parking: parking1,
                garage_atached: garageAtached,
                fireplace: fireplace,
                window_covering: windowCovering,
                backyard: backyard,
                fitness_gym: fitnessGym,
                elevator: elevator,
                others: othersArray // Burada tüm dinamik değerler gönderilecek
            };

            const addAmenitiesResponse = await fetch("http://127.0.0.1:8081/property/add-amenities", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(amenitiesData),
            });

            if (!addAmenitiesResponse.ok) {
                const errorText = await addAmenitiesResponse.text();
                showModal("error", "Hata!", `Amenities oluşturulurken bir hata oluştu! Hata: ${errorText}`);
                return;
            }

            const addAmenitiesResult = await addAmenitiesResponse.json();

            if (addAmenitiesResult.status !== 200) {
                showModal("error", "Hata!", "Amenities oluşturulamadı: " + addAmenitiesResult.message);
                return;
            }


            // **7. `Upload Plans and Brochure` oluştur**
            const fileInput1 = document.querySelector('input[type="file"][multiple]');
            await handlePlansAndBrochures(propertyID, fileInput1); // **propertyID GEÇİRİLDİ**



            // **8. Accordion Widget oluştur
            const accordion_exist = document.querySelector('input[name="accordion_exist"]').checked;
            const accordionTitle = document.querySelector('input[name="accordion_title"]').value; // Inputtan değeri al
            const accordionDetails = document.querySelector('textarea[name="accordion_details"]').value; // Inputtan değeri al

            const accordionWidgetData = {
                property_id: propertyID.property_id,
                accordion_exist: accordion_exist,
                accordion_title: accordionTitle,
                accordion_details: accordionDetails,
            }

            const accordionWidgetResponse = await fetch("http://127.0.0.1:8081/property/add-accordion-widget", {
                method: "POST",
                headers : {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(accordionWidgetData),
            });

            if(!accordionWidgetResponse.ok){
                const errorText = await accordionWidgetResponse.text();
                showModal("error", "Hata!", `Amenities oluşturulurken bir hata oluştu! Hata: ${errorText}`);
                return;
            }

            const accordionWidgetResult = await accordionWidgetResponse.json();

            if(accordionWidgetResult.status !== 200){
                showModal("error", "Hata!", "Amenities oluşturulamadı: " + accordionWidgetResult.message);
                return;
            }


            // **8. Video Widget oluştur
            const video_exist = document.querySelector('input[name="video_exist"]').checked;
            const videoTitle = document.querySelector('input[name="video_title"]').value; // Inputtan değeri al
            const youtube_url = document.querySelector('input[name="youtube_url"]').value; // Inputtan değeri al
            const vimeo_url = document.querySelector('input[name="vimeo_url"]').value; // Inputtan değeri al

            const videoWidgetData = {
                property_id: propertyID.property_id,
                video_exist: video_exist,
                video_title: videoTitle,
                youtube_url: youtube_url,
                vimeo_url: vimeo_url,
            }

            const videoWidgetResponse = await fetch("http://127.0.0.1:8081/property/add-video-widget", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(videoWidgetData),
            });

            if(!videoWidgetResponse.ok){
                const errorText = await videoWidgetResponse.text();
                showModal("error", "Hata!", `VideoWidget oluşturulurken bir hata oluştu! Hata: ${errorText}`);
                return;
            }

            const videoWidgetResult = await videoWidgetResponse.json();

            if(videoWidgetResult.status !== 200){
                showModal("error", "Hata!", "Amenities oluşturulamadı: " + videoWidgetResult.message);
                return;
            }




            showModal("success", "Başarılı!", "Emlak başarıyla eklendi!");


        } catch (error) {
            console.error("Hata oluştu: ", error);
            showModal("error", "Hata!", "Bir hata oluştu!");
        } finally {
            loaderWrap.style.display = "none";
            loader.classList.add("d-none");
        }

        const modalElement = document.getElementById("kt_modal_1");
        const modal = bootstrap.Modal.getInstance(modalElement);

        modalElement.addEventListener("hidden.bs.modal", function () {
            window.location.reload();
        });

        modal.hide();
    });
});





// **Resimleri yükleme ve Property Media oluşturma fonksiyonu**

// **PLans and brochures
async function handlePlansAndBrochures(propertyID, fileInput1) { // **propertyID EKLENDİ**
    const files = fileInput1.files;

    if (files.length > 0) {
        const imageIDs = [];

        for (let i = 0; i < files.length; i++) {
            const file = files[i];

            const formData = new FormData();
            formData.append("property_id", propertyID.property_id);
            formData.append("file_type", file.name);
            formData.append("file_path", file);

            const addPlansAndBrochuresResponse = await fetch("http://127.0.0.1:8081/property/add-plans-brochures", {
                method: "POST",
                body: formData,
            });

            if (!addPlansAndBrochuresResponse.ok) {
                const errorText = await addPlansAndBrochuresResponse.text();
                showModal("error", "Hata!", `Resim yüklenirken bir hata oluştu! Hata: ${errorText}`);
                return;
            }

            const addPlansAndBrochuresResult = await addPlansAndBrochuresResponse.json();

            if (addPlansAndBrochuresResult.status !== 200 ) {
                showModal("error", "Hata!", "Resim yüklenemedi: " + addPlansAndBrochuresResult.message);
                return;
            }
        }
    }
}




function showModal(type, title, message) {
    const modalTitle = document.getElementById("kt_modal_1").querySelector(".modal-title");
    const modalBody = document.getElementById("kt_modal_1").querySelector(".modal-body");
    const modalFooter = document.getElementById("kt_modal_1").querySelector(".modal-footer");

    if (type === 'success') {
        modalTitle.textContent = title;
        modalBody.innerHTML = `<p class="text-success">${message}</p>`;
        modalFooter.innerHTML = `<button type="button" class="btn btn-light" data-bs-dismiss="modal">Kapat</button>`;
    } else if (type === 'error') {
        modalTitle.textContent = title;
        modalBody.innerHTML = `<p class="text-danger">${message}</p>`;
        modalFooter.innerHTML = `<button type="button" class="btn btn-light" data-bs-dismiss="modal">Kapat</button>`;
    }

    const modal = new bootstrap.Modal(document.getElementById("kt_modal_1"));
    modal.show();
}