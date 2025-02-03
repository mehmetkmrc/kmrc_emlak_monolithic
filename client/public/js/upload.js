document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("load"); // Form elementini al

    form.addEventListener("submit", async function (event) {
        event.preventDefault(); // Formun varsayılan submit işlemini engelle

        // Formdaki inputları al
        const mainTitle = document.querySelector('input[placeholder="Main Title"]').value;
        const type = document.querySelector('select[data-placeholder="Categories"]').value;
        const category = document.querySelectorAll('select[data-placeholder="Categories"]')[1].value;
        const price = document.querySelector('input[placeholder="Price"]').value;
        const keywords = document.querySelector('input[placeholder="Keywords"]').value;

        // API'ye gönderilecek JSON objesi
        const requestData = {
            main_title: mainTitle,
            type: type,
            category: category,
            price: parseFloat(price), // Sayısal değeri float olarak parse et
            keywords: keywords
        };

        try {
            const response = await fetch("http://localhost:3000/api/add-basic-info", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(requestData)
            });

            const result = await response.json();

            if (response.ok) {
                alert("Basic Info başarıyla eklendi!");
                console.log("Başarılı:", result);
            } else {
                alert("Hata oluştu: " + result.message);
                console.error("Hata:", result);
            }
        } catch (error) {
            console.error("İstek gönderilirken hata oluştu:", error);
            alert("Sunucuya bağlanırken hata oluştu!");
        }
    });
});
