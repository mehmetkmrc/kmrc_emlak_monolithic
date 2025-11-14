var globalImgArray = [];

jQuery(document).ready(function () {
    ImgUpload();
});
  
function ImgUpload() {
    var imgWrap = "";

    $('.upload__inputfile').each(function () {
        var uploadBox = $(this).closest('.upload__box');
        imgWrap = uploadBox.find('.upload__img-wrap');
        var maxLength = $(this).attr('data-max_length');

        // Input change event
        $(this).on('change', function (e) {
            handleFiles(e.target.files, imgWrap, maxLength);
            this.value = ''; // input reset
        });

        // Drag & Drop
        uploadBox.on('dragover', function(e) {
            e.preventDefault();
            e.stopPropagation();
            $(this).addClass('drag-over');
        });

        uploadBox.on('dragleave drop', function(e) {
            e.preventDefault();
            e.stopPropagation();
            $(this).removeClass('drag-over');
        });

        uploadBox.on('drop', function(e) {
            var dt = e.originalEvent.dataTransfer;
            var files = dt.files;
            handleFiles(files, imgWrap, maxLength);
        });
    });

    function handleFiles(files, wrap, maxLength) {
        var filesArr = Array.prototype.slice.call(files);
        filesArr.forEach(function(f) {
            if (!f.type.match('image.*')) return;

            if (globalImgArray.length >= maxLength) return;

            globalImgArray.push(f);
            console.log("Resim eklendi: " + f.name);

            var reader = new FileReader();
            reader.onload = function(e) {
                var html = `
                <div class='upload__img-box'>
                    <div class='img-bg' style='background-image: url(${e.target.result})' data-file='${f.name}'></div>
                    <div class='upload__img-close'>×</div>
                </div>`;
                wrap.append(html);

                // Büyük açma
                wrap.find('.img-bg').last().on('click', function() {
                    const win = window.open();
                    win.document.write(`<img src="${e.target.result}" style="width:100%">`);
                });
            };
            reader.readAsDataURL(f);
        });
    }

    $('body').on('click', ".upload__img-close", function () {
        var file = $(this).siblings('.img-bg').data("file");
        for (var i = 0; i < globalImgArray.length; i++) {
            if (globalImgArray[i].name === file) {
                globalImgArray.splice(i, 1);
                console.log("Resim silindi: " + file);
                break;
            }
        }
        $(this).parent().remove();
    });
}
