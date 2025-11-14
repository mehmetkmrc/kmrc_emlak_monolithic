
// 📌 Sayfa yüklenince her şeyi başlat
document.addEventListener("DOMContentLoaded", function () {

    initMap();
    initNearby();
    initImageUpload();

});

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
    const uploadWrap = inputFile.closest('.upload__box').querySelector('.upload__img-wrap');
    const files = Array.from(inputFile.files);

    files.forEach(file => {
        if (!file.type.startsWith('image/')) {
            alert('Lütfen sadece resim yükleyin!');
            return;
        }

        globalImgArray.push(file);

        const reader = new FileReader();
        reader.onload = (e) => {
            const imgBox = document.createElement('div');
            imgBox.classList.add('upload__img-box');

            const img = document.createElement('img');
            img.src = e.target.result;

            img.onclick = () => {
                const win = window.open();
                win.document.write(`<img src="${e.target.result}" style="width:100%">`);
            };

            const removeBtn = document.createElement('button');
            removeBtn.textContent = "×";
            removeBtn.onclick = () => {
                globalImgArray = globalImgArray.filter(f => f !== file);
                imgBox.remove();
            };

            imgBox.appendChild(img);
            imgBox.appendChild(removeBtn);
            uploadWrap.appendChild(imgBox);
        };

        reader.readAsDataURL(file);
    });

    inputFile.value = "";
}
