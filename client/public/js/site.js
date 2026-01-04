
// 📌 Sayfa yüklenince her şeyi başlat
document.addEventListener("DOMContentLoaded", function () {

    initMap();
    initNearby();
    initImageUpload();

});
  let imageState = {
            gallery: [],
            plan: [],
            brochure: [],
            profile: []
    };

// =======================
// 🌍 MAP INIT
// =======================
function initMap() {
    const mapElement = document.getElementById("singleMap");
    if (!mapElement) return;

    const initialLat = parseFloat(mapElement.dataset.latitude) || 40.7427837;
    const initialLng = parseFloat(mapElement.dataset.longitude) || -73.11445617675781;

    const map = L.map('singleMap').setView([initialLat, initialLng], 13);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    const marker = L.marker([initialLat, initialLng], { draggable: true }).addTo(map);

    const latInput = document.getElementById("lat");
    const lngInput = document.getElementById("long");

    latInput.value = initialLat;
    lngInput.value = initialLng;

    marker.on('dragend', function (e) {
        const latlng = e.target.getLatLng();
        latInput.value = latlng.lat.toFixed(6);
        lngInput.value = latlng.lng.toFixed(6);
    });
}

// =======================
// 🏙️ NEARBY INIT
// =======================
function initNearby() {
    const addBtn = document.getElementById("add-nearby-btn");
    const listContainer = document.getElementById("nearby-list");

    if (!addBtn) return;

    addBtn.addEventListener("click", () => {
        const places = document.querySelector('select[name="places"]').value;
        const distance = document.querySelector('input[name="distance"]').value;

        if (!places || !distance) {
            showModal("error", "Hata!", "Lütfen hem yer hem mesafe giriniz!");
            return;
        }

        const nearbyObj = { places, distance };
        globalNearbyArray.push(nearbyObj);

        const el = document.createElement("div");
        el.classList.add("nearby-item");
        el.innerHTML = `
            ${places} - ${distance} m
            <button onclick="removeNearby(this)">Sil</button>
        `;
        listContainer.appendChild(el);

        document.querySelector('select[name="places"]').selectedIndex = 0;
        document.querySelector('input[name="distance"]').value = "";
    });
}

function removeNearby(btn) {
    const index = Array.from(document.querySelectorAll('#nearby-list .nearby-item'))
        .indexOf(btn.parentElement);

    if (index > -1) globalNearbyArray.splice(index, 1);

    btn.parentElement.remove();
}

async function removeNearbyForEditPage(btn) {
    const item = btn.parentElement;
    const nearbyID = item.dataset.id;

    try {
        const res = await fetch(`http://127.0.0.1:8081/update-property/nearby/${nearbyID}`, {
            method: 'DELETE'
        });

        if (!res.ok) {
            const txt = await res.text();
            showModal("error", "Hata", txt);
            return;
        }

        item.remove();

    } catch (err) {
        console.error(err);
        showModal("error", "Hata", err.message);
    }
}


// =======================
// 📸 IMAGE UPLOAD INIT
// =======================
function initImageUpload() {
    const uploadInputs = document.querySelectorAll('.upload__inputfile');
    uploadInputs.forEach(input => {
        input.addEventListener('change', () => handleUpload(input));
    });
}

function handleUpload(inputFile) {

    const key = inputFile.dataset.type;
    if (!key) return console.error("data-type eksik", inputFile);

    const uploadWrap = inputFile.closest('.upload__box').querySelector('.upload__img-wrap');
    if (!uploadWrap) return console.error("upload__img-wrap bulunamadı");

    const files = Array.from(inputFile.files);

    files.forEach(file => {

        if (!file.type.startsWith('image/')) return;

        imageState[key].push(file);

        const reader = new FileReader();
        reader.onload = e => {

            const box = document.createElement('div');
            box.className = 'upload__img-box';

            const img = document.createElement('img');
            img.src = e.target.result;

            const btn = document.createElement('button');
            btn.textContent = "×";
            btn.onclick = () => {
                imageState[key] = imageState[key].filter(f => f !== file);
                box.remove();
            };

            box.append(img, btn);
            uploadWrap.appendChild(box);
        };

        reader.readAsDataURL(file);
    });

    inputFile.value = "";
}

async function deleteProperty(propertyId) {
    if (!propertyId) {
        showModal("error", "Hata", "İlan ID bulunamadı");
        return;
    }

    const confirmed = confirm("Bu ilanı silmek istediğinize emin misiniz?");
    if (!confirmed) return;

    try {
        const res = await fetch(`http://127.0.0.1:8081/update-property/delete/${propertyId}`, {
            method: "DELETE"
        });

        if (!res.ok) {
            const txt = await res.text();
            showModal("error", "Silme Hatası", txt);
            return;
        }

        showModal("success", "Başarılı", "İlan başarıyla kaldırıldı");

        // Satırı DOM'dan kaldır (yenilemeden)
        const item = document.querySelector(`[data-property-id="${propertyId}"]`);
        if (item) item.remove();

    } catch (err) {
        console.error(err);
        showModal("error", "Sunucu Hatası", err.message);
    }
}

// =======================
// ⏸️ PROPERTY PASSIVE
// =======================
async function passiveProperty(propertyId) {
    if (!propertyId) {
        showModal("error", "Hata", "İlan ID bulunamadı");
        return;
    }

    const confirmed = confirm("Bu ilanı pasife almak istediğinize emin misiniz?");
    if (!confirmed) return;

    try {
        const res = await fetch(`http://127.0.0.1:8081/update-property/passive/${propertyId}`, {
            method: "PUT"
        });

        if (!res.ok) {
            const txt = await res.text();
            showModal("error", "Hata", txt);
            return;
        }

        showModal("success", "Başarılı", "İlan pasife alındı");

        // UI güncelle (ikon değiştir / satırı soldur)
        const item = document.querySelector(`[data-property-id="${propertyId}"]`);
        if (item) {
            item.classList.add("property-passive");
        }

    } catch (err) {
        console.error(err);
        showModal("error", "Sunucu Hatası", err.message);
    }
}
