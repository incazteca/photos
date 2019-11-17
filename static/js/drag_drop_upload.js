"use strict";
window.onload = function () {
    const dropTarget = document.getElementById('drop-target');
    dropTarget.addEventListener('dragover', dragoverHandler);
    dropTarget.addEventListener('drop', dropHandler);
};
function dragoverHandler(event) {
    console.log('File(s) in drop zone');
    // Prevent default behavior (Prevent file from being opened)
    event.preventDefault();
}
function dropHandler(event) {
    console.log('File(s) dropped');
    // Prevent default behavior (Prevent file from being opened)
    event.preventDefault();
    const dataTransfer = event.dataTransfer;
    if (dataTransfer == null) {
    }
    else if (dataTransfer.items) {
        // Use DataTransferItemList interface to access the file(s)
        for (var i = 0; i < dataTransfer.items.length; i++) {
            // If dropped items aren't files, reject them
            if (dataTransfer.items[i].kind === 'file') {
                var file = dataTransfer.items[i].getAsFile();
                if (file != null) {
                    console.log('... file[' + i + '].name = ' + file.name);
                    uploadPhoto(file);
                }
            }
        }
    }
    else {
        // Use DataTransfer interface to access the file(s)
        for (var i = 0; i < dataTransfer.files.length; i++) {
            console.log('... file[' + i + '].name = ' + dataTransfer.files[i].name);
            uploadPhoto(dataTransfer.files[i]);
        }
    }
}
function uploadPhoto(photo) {
    const xhr = new XMLHttpRequest();
    const server = 'http://localhost:8080';
    xhr.open('POST', server + '/photo');
    xhr.send(photo);
}
