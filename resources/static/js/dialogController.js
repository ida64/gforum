window.addEventListener("load", function (event) {
    var dialogHeader = document.querySelector('.dialog-header');
    var isMouseDown = false;
    var offset = [0, 0];

    dialogHeader.addEventListener('mousedown', function (e) {
        isMouseDown = true;
        offset = [
            dialogHeader.parentElement.offsetLeft - e.clientX,
            dialogHeader.parentElement.offsetTop - e.clientY
        ];
    }, true);

    document.addEventListener('mouseup', function () {
        isMouseDown = false;
    }, true);

    document.addEventListener('mousemove', function (event) {
        event.preventDefault();
        if (isMouseDown) {
            dialogHeader.parentElement.style.left = (event.clientX + offset[0]) + 'px';
            dialogHeader.parentElement.style.top = (event.clientY + offset[1]) + 'px';
        }
    }, true);

    var dialogClose = dialogHeader.querySelector('.dialog-close');
    dialogClose.addEventListener('click', function (e) {
        dialogHeader.parentElement.remove();
    });
});
