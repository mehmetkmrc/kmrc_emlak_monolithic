document.getElementById("updateProfileBtn").addEventListener("click", async function (e) {
    e.preventDefault();

    const userID = document.getElementById("userIDHidden").value;

    const payload = {
        user_id: userID,
        name: document.querySelector("input[placeholder='First Name']").value,
        surname: document.querySelector("input[placeholder='Second Name']").value,
        email: document.querySelector("input[placeholder='Email Address']").value,
        phone: document.querySelector("input[placeholder='Phone']").value,
        about: document.getElementById("comments").value,
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

    showModal("success", "Başarılı!", "Kullanıcı güncellendi!");
});
