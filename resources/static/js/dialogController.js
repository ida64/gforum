function getHighestZIndex() {
    const elements = document.getElementsByTagName('*');
    let highestZIndex = 0;

    for (let i = 0; i < elements.length; i++) {
        const zIndex = parseInt(window.getComputedStyle(elements[i]).getPropertyValue('z-index'));
        if (zIndex > highestZIndex) {
            highestZIndex = zIndex;
        }
    }

    return highestZIndex;
}

function initializeDialogController() {
    const dialog = document.querySelector('.dialog-header');
    if(!dialog)
    {
        return;
    }

    let isMouseDown = false;
    let offset = [0, 0];

    const centerX = (window.innerWidth - dialog.parentElement.offsetWidth);
    const centerY = (window.innerHeight / 2) - (dialog.parentElement.offsetHeight / 2);

    dialog.parentElement.style.left = centerX + 'px';
    dialog.parentElement.style.top = centerY + 'px';

    dialog.addEventListener('mousedown', (e) => {
        isMouseDown = true;
        offset = [
            dialog.parentElement.offsetLeft - e.clientX,
            dialog.parentElement.offsetTop - e.clientY
        ];
    }, true);

    document.addEventListener('mouseup', () => {
        isMouseDown = false;
    }, true);

    dialog.parentElement.addEventListener('click', () => {
        dialog.parentElement.style.zIndex = getHighestZIndex() + 1;
    });

    document.addEventListener('mousemove', (event) => {
        event.preventDefault();
        if (isMouseDown) {
            dialog.parentElement.style.left = (event.clientX + offset[0]) + 'px';
            dialog.parentElement.style.top = (event.clientY + offset[1]) + 'px';
        }
    }, true);

    const dialogClose = dialog.querySelector('.dialog-close');
    dialogClose.addEventListener('click', () => {
        dialog.parentElement.remove();
    });
}

document.addEventListener("htmx:load", () => {
    initializeDialogController();
});
