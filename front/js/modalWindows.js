const parent = document.querySelector("#projects");

setTimeout(async() => {
    const response = await axios.get("https://webprojass4-default-rtdb.firebaseio.com/posts.json")
    const data = response.data
    console.log(data);
    Object.keys(data).map((key) => {
        const child =  document.createElement("a");
        child.className = "cards flex-shrink-1 text-decoration-none"

        child.id = data[key].state
        child.setAttribute('data', data[key].data)
        child.setAttribute('dataImage', data[key].img)
        child.setAttribute('discription', data[key].discription)
        child.setAttribute('dataText',data[key].title)
        const childImg = document.createElement("img");
        childImg.style = "height:405px"
        childImg.src = data[key].img
        childImg.className = "img-fluid"

        const childH = document.createElement("h5");
        childH.innerHTML = data[key].title
        childH.className = "mt-3 text-primary-emphasis"

        child.appendChild(childImg)
        child.appendChild(childH)
        parent.appendChild(child)
    })

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

    const modalMain =document.createElement('div');
    modalMain.className = 'modalMain row container-fluid'
    modalContainer.appendChild(modalMain);

    const modalMainDiv =document.createElement('div');
    modalMainDiv.className='first_content col-md-4 d-flex flex-column';
    modalMainDiv.style = "height: 405px"
    modalMain.appendChild(modalMainDiv);

    const imageElement = document.createElement('img');
    imageElement.src = imageSrc;
    imageElement.className = "img-fluid"
    modalMainDiv.appendChild(imageElement);

    const textElement = document.createElement('p');
    textElement.innerText = textContent;
    textElement.style = "position: relative"
    modalMainDiv.appendChild(textElement);
    //--------------------------------------------------------------

    const modalMainDiv2 = document.createElement('div');
    modalMainDiv2.className = 'second_content col-md-8';
    modalMain.appendChild(modalMainDiv2);

    const textForm = document.createElement('p');
    textForm.innerText = textSpace;
    modalMainDiv2.appendChild(textForm);

    const buttonOk = document.createElement('button');
    buttonOk.className = 'modalOkButton';
    buttonOk.innerText = 'X';
    buttonOk.addEventListener('click', () => {
        wrapper.remove()
    })

    modalMainDiv2.appendChild(buttonOk);

    document.body.appendChild(wrapper);
}

function openSignUpModal() {
    document.getElementById('signupModal').style.display = 'block';
}

// Функция для закрытия модального окна "Sign Up"
function closeSignUpModal() {
    document.getElementById('signupModal').style.display = 'none';
}
}
)


/* <a href="#" class="cards flex-shrink-1 text-decoration-none" id="finishedProj" data-image="./img/projects-pngs.png" data-text="Ramada plaza & Dar Tower" text="Ramada Plaza & Dar Tower: Ваша Оазис в Туристической Зоне



Добро пожаловать в Ramada Plaza & Dar Tower – мультифункциональный комплекс, идеально расположенный в самой сердце туристически привлекательной зоны нашего города. Наш отель и апартаменты предлагают незабываемое сочетание комфорта, роскоши и удобства, чтобы сделать ваше пребывание незабываемым.

Совершенное расположение в туристической зоне города делает наш комплекс идеальным местом для тех, кто желает насладиться всем, что наш город может предложить. Наши гости находятся в нескольких шагах от самых известных туристических достопримечательностей, ресторанов и магазинов, что позволяет им наслаждаться всей культурой и развлечениями, которые город предлагает.">
                    <img src="./img/projects-pngs.png" class='img-fluid' alt="123" >
                    <h5 class="mt-3 text-primary-emphasis">Ramada plaza & Dar Tower</h5>
                </a>

<a href="#" class="cards flex-shrink-1 text-decoration-none" id="finishedProj" data-image="./img/projects-pngs-1.png" data-text="Hampton by Hilton" text="Hampton by Hilton: Ваш Идеальный Выбор для Комфортного Проживания

Добро пожаловать в Hampton by Hilton – ваш ключ к комфортному и незабываемому проживанию, где каждая деталь заботливо продумана для вашего удовольствия. Наш отель предоставляет идеальное сочетание уютных номеров, превосходного обслуживания и удобного расположения, чтобы сделать ваше пребывание незабываемым.

Расположенный в самом сердце города, Hampton by Hilton является идеальным местом для путешественников всех типов. Независимо от того, приезжаете ли вы в командировку или на отдых, наш отель находится в нескольких шагах от основных достопримечательностей, бизнес-центров и множества ресторанов и магазинов.">
                    <img src="./img/projects-pngs-1.png" class='img-fluid' alt="">
                    <h5 class="mt-3 text-primary-emphasis">Hampton by Hilton</h5>
                </a>
                <a href="#" class="cards flex-shrink-1 text-decoration-none" id="inProcess" data-image="./img/projects-pngs-1%20(2).png" data-text="Sunrise" text="Sunrise: Где Каждое Утро Начинается с Яркого Сияния

Добро пожаловать в Sunrise – ваш уголок спокойствия и красоты, где каждое утро начинается с мягких лучей солнца, напоминая о новом дне и бескрайних возможностях. Наша локация предоставляет идеальное сочетание природной красоты, комфорта и удобств, чтобы сделать ваше пребывание незабываемым.

Расположенный в живописном уголке природы, Sunrise предоставляет уединение и покой для тех, кто ищет отдых и релаксацию. Наши номера и апартаменты предлагают уютное пристанище, откуда вы сможете наслаждаться потрясающими видами и уникальными впечатлениями.">
                    <img src="./img/projects-pngs-1%20(2).png" class='img-fluid' alt="">
                    <h5 class="mt-3 text-primary-emphasis">Sunrise</h5>
                </a>
                <a href="#" class="cards flex-shrink-1 text-decoration-none" id="designProcess" data-image="./img/projects-pngs-1%20(3).png" data-text="Era Palace" text="Era Palace: Ваша Империя Комфорта и Роскоши

Добро пожаловать в Era Palace – ваша королевская резиденция, где каждая минута пребывания приносит невероятные впечатления и бесконечную роскошь. Наш отель предлагает вам уникальное сочетание элегантности, комфорта и выдающегося обслуживания, чтобы сделать ваше пребывание незабываемым.

Расположенный в самом центре города, Era Palace стоит в одном ряду с культурными и историческими достопримечательностями. Независимо от того, приезжаете ли вы в деловую поездку или отдых, наш отель находится в непосредственной близости от многочисленных ресторанов, магазинов и мест для развлечений.">
                    <img src="./img/projects-pngs-1%20(3).png" class='img-fluid' alt="" >
                    <h5 class="mt-3 text-primary-emphasis">Era Palace</h5>
                </a>
                <a href="#" class="cards flex-shrink-1 text-decoration-none" id="inProcess" data-image="./img/projects-pngs-1%20(4).png" data-text="DS Group" text="DS Group: Формируем Лучшее Завтра

Добро пожаловать в DS Group, где наше стремление к инновациям и превосходству формирует светлое будущее. Как ведущий конгломерат, DS Group предан созданию положительных изменений через наше разнообразное портфолио бизнесов с целью улучшения жизни людей по всему миру.

Наши достижения определяются постоянным качеством и устойчивостью. С нашим наследием успехов в различных отраслях, включая сельское хозяйство, производство и технологии, DS Group постоянно исследует новые горизонты. Наше стремление к передовым исследованиям и разработкам гарантирует, что наши продукты и услуги всегда соответствуют высшим стандартам отрасли.">
                    <img src="./img/projects-pngs-1%20(4).png" class='img-fluid' alt="" >
                    <h5 class="mt-3 text-primary-emphasis">DS Group</h5>
                </a>
                <a href="#" class="cards flex-shrink-1 text-decoration-none" id="inProcess" data-image="./img/projects-pngs-1%20(5).png" data-text="Vox" text="Vox: Ваш Голос в Мире Инноваций

Добро пожаловать в мир Vox, где ваш голос имеет значение, и инновации определяют наше будущее. Наша компания предоставляет платформу, на которой каждый может выразить свои идеи, создавать, и развиваться.

Мы стремимся к тому, чтобы каждый, кто приходит к нам, чувствовал себя вдохновленным и поддержанным в своих усилиях. Мы поощряем креативность, и исследования, что позволяет нам оставаться на передовой технологических и культурных тенденций. Мы поддерживаем стартапы, исследования и проекты, которые могут привести к революционным изменениям в области технологий, бизнеса и общества.">
                    <img src="./img/projects-pngs-1%20(5).png" class='img-fluid' alt="">
                    <h5 class="mt-3 text-primary-emphasis">Vox</h5>
                </a>
                <a href="#" class="cards flex-shrink-1 text-decoration-none" id="designProcess" data-image="./img/projects-pngs-1%20(6).png" data-text="Navi" text="Navi: Ваша Элегантная Оазис в Мире Жилых Комплексов

Добро пожаловать в Navi, вашу идеальную резиденцию в мире жилых комплексов. Мы предлагаем вам не только дом, но и уникальное пространство, где стиль, комфорт и удобство сочетаются в идеальном балансе.

Navi располагается в живописной и удобной локации, предлагая своим жителям роскошные апартаменты и уникальные возможности для комфортного проживания. Наши апартаменты сочетают в себе современный дизайн, качественные материалы и внимание к каждой детали, чтобы создать идеальное жилье. Мы гордимся обширными инфраструктурными возможностями и сервисами, включая фитнес-центр, спа, бассейн и многое другое, чтобы удовлетворить запросы наших жителей.">
                    <img src="./img/projects-pngs-1%20(6).png" class='img-fluid' alt="">
                    <h5 class="mt-3 text-primary-emphasis">Navi</h5>
                </a>
                <a href="#" class="cards flex-shrink-1 text-decoration-none" id="designProcess" data-image="./img/projects-pngs-1%20(7).png" data-text="New Wave" text="New Wave: Ваша Перспектива на Будущее

Добро пожаловать в мир New Wave, где каждый день принесет новые горизонты и возможности. Мы создаем атмосферу, где инновации, творчество и вдохновение идут рука об руку.

В New Wave мы стремимся быть на переднем крае перемен. Наше название говорит само за себя - мы поддерживаем свежие идеи, новаторские проекты и смелые решения. Мы предлагаем платформу для тех, кто хочет изменить мир, начиная с собственных начинаний. Мы создаем сообщество, где люди с разными бакграундами и интересами могут обмениваться идеями и сотрудничать над проектами, которые могут изменить будущее. ">
                    <img src="./img/projects-pngs-1%20(7).png" class='img-fluid' alt="">
                    <h5 class="mt-3 text-primary-emphasis">New Wave</h5>
                </a> */