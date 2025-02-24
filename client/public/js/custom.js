var globalImgArray = [];

jQuery(document).ready(function () {
    ImgUpload();
});
  
function ImgUpload() {
    var imgWrap = "";
  
    $('.upload__inputfile').each(function () {
        $(this).on('change', function (e) {
            imgWrap = $(this).closest('.upload__box').find('.upload__img-wrap');
            var maxLength = $(this).attr('data-max_length');
  
            var files = e.target.files;
            var filesArr = Array.prototype.slice.call(files);
            var iterator = 0;
            
            filesArr.forEach(function (f, index) {
                if (!f.type.match('image.*')) {
                    return;
                }
  
                if (globalImgArray.length > maxLength) {
                    return false
                } else {
                    var len = 0;
                    for (var i = 0; i < globalImgArray.length; i++) {
                        if (globalImgArray[i] !== undefined) {
                            len++;
                        }
                    }
                    if (len > maxLength) {
                        return false;
                    } else {
                        globalImgArray.push(f);
                        console.log("Resim eklendi: " + f.name + ", Toplam: " + globalImgArray.length);
  
                        var reader = new FileReader();
                        reader.onload = function (e) {
                            var html = "<div class='upload__img-box'><div style='background-image: url(" + e.target.result + ")' data-number='" + $(".upload__img-close").length + "' data-file='" + f.name + "' class='img-bg'><div class='upload__img-close'></div></div></div>";
                            imgWrap.append(html);
                            iterator++;
                        }
                        reader.readAsDataURL(f);
                    }
                }
            });
        });
    });
  
    $('body').on('click', ".upload__img-close", function (e) {
        var file = $(this).parent().data("file");
        for (var i = 0; i < globalImgArray.length; i++) {
            if (globalImgArray[i].name === file) {
                globalImgArray.splice(i, 1);
                console.log("Resim silindi: " + file + ", Kalan: " + globalImgArray.length);
                break;
            }
        }
        $(this).parent().parent().remove();
    });
}