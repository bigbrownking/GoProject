const openModalWindow = document.querySelectorAll('.cards');

openModalWindow.forEach((card) => {
    card.addEventListener('click', () => {
        const imageSrc = card.getAttribute('dataImage');
        const textContent = card.getAttribute('dataText');
        const textSpace = card.getAttribute('discription')
        console.log(textSpace);
        createModal(imageSrc, textContent, textSpace);
    });
});

function createModal(imageSrc, textContent, textSpace){
    const wrapper = document.createElement('div');
    wrapper.className = 'modalWrapper'

    const backdrop = document.createElement('div');
    backdrop.className = 'backdrop';
    wrapper.appendChild(backdrop);
    backdrop.addEventListener('click', () => {
        wrapper.remove();
    })

    const modalContainer = document.createElement('div');
    modalContainer.className = 'modall';
    wrapper.appendChild(modalContainer);

    const modalHeader = document.createElement('header');
    modalHeader.className = 'modal_header'
    modalContainer.appendChild(modalHeader);

    const closeButton = document.createElement('button');
    closeButton.innerText = 'X';
    closeButton.className = 'closeButton';
    closeButton.addEventListener('click', () => {
        wrapper.remove();
    });
    modalHeader.appendChild(closeButton);

    const modalMain = document.createElement('main');
    modalMain.className = 'modal_main'
    modalContainer.appendChild(modalMain);

    const imageElement = document.createElement('img');
    imageElement.src = imageSrc;
    modalMain.appendChild(imageElement);

    const textElement = document.createElement('p');
    textElement.innerText = textContent;
    modalMain.appendChild(textElement);

    const modalFooter = document.createElement('footer');
    modalFooter.className = 'modal_footer';
    modalContainer.appendChild(modalFooter);

    const textForm = document.createElement('p');
    textForm.innerText = textSpace;
    modalFooter.appendChild(textForm);

    document.body.appendChild(wrapper);
}