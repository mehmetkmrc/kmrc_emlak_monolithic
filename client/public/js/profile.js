document.getElementById("updateProfileBtn").addEventListener("click", async function (e) {
    e.preventDefault();

    const userID = document.getElementById("userIDHidden").value;

    const payload = {
        user_id: userID,
        first_name: document.querySelector("input[placeholder='First Name']").value,
        last_name: document.querySelector("input[placeholder='Second Name']").value,
        email: document.querySelector("input[placeholder='Email Address']").value,
        phone: document.querySelector("input[placeholder='Phone']").value,
        about_text: document.getElementById("comments").value,
    };

    const response = await fetch("http://127.0.0.1:8081/user/update-user-base-info", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
    });

    const data = await response.json();

    if (!response.ok) {
        showModal("error", "Hata!", "User update edilemedi: " + data.message);
        return;
    }

    showModal("success", "Başarılı", "Kullanıcı bilgileri başarıyla güncellendi!");
});


// Save socials
document.getElementById("saveSocialsBtn").addEventListener("click", async function(e) {
    e.preventDefault();
    const userID = document.getElementById("userIDHidden").value;

    const payload = {
    user_id: userID,
    facebook: document.getElementById('facebook').value,
    tiktok: document.getElementById('tiktok').value,
    instagram: document.getElementById('instagram').value,
    twitter: document.getElementById('twitter').value,
    youtube: document.getElementById('youtube').value,
    linkedin: document.getElementById('linkedin').value,
    };


    const res = await fetch('http://127.0.0.1:8081/user/social-links', {
    method: 'PUT',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify(payload)
    });


    const data = await res.json();
    if (!res.ok) {
    showModal('error','Hata','Sosyal linkler kaydedilemedi');
    return;
    }
    showModal('success','Başarılı','Sosyal linkler kaydedildi');
});