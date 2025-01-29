document.addEventListener("DOMContentLoaded", function () {
    const loginForm = document.querySelector('form[name="registerform1"]');
    const modal = document.getElementById("successModal");
    const modalClose = document.getElementById("closeModal");

    loginForm.addEventListener("submit", async function (event) {
        event.preventDefault(); // Formun otomatik olarak submit edilmesini engeller

        // Formdaki verileri almak
        const first_name = loginForm.querySelector('input[name="first_name"]').value;
        const last_name = loginForm.querySelector('input[name="last_name"]').value;
        const email = loginForm.querySelector('input[name="email"]').value;
        const phone = loginForm.querySelector('input[name="phone"]').value;
        const password = loginForm.querySelector('input[name="password"]').value;
        const confirm_password = loginForm.querySelector('input[name="confirm_password"]').value;

        // Veri objesi oluştur ve JSON formatına çevir
        const requestData = { first_name: first_name, last_name: last_name ,email: email, phone: phone, password: password, confirm_password: confirm_password};

        try {
            const response = await fetch("/auth/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"  // JSON olarak gönderildiğini belirt
                },
                body: JSON.stringify(requestData)  // Objemizi JSON formatına çeviriyoruz
            });

            if (response.ok) {
                const data = await response.json();

                // Başarılı giriş
                console.log("Giriş başarılı:", data);
                 // Modal'ı göster
                 modal.style.display = "block";

                 // Modal kapatıldıktan sonra sayfayı yenile
                 modalClose.addEventListener("click", function () {
                     modal.style.display = "none";
                     location.reload(); // Sayfayı yenile
                 });
 
                
            } else {
                // Hata mesajını göster
                const errorData = await response.json();
                console.error("Kayıt hatası:", errorData.message);
                alert("Kayıt başarısız. Lütfen bilgilerinizi kontrol edin.");
            }
        } catch (error) {
            console.error("Bir hata oluştu:", error);
            alert("Sunucuya bağlanırken bir sorun oluştu.");
        }
    });
});
