const loginBtn = document.querySelector('#getLogin');
let isNameValid = false
let isPassValid = false
let user={};
const formGroup = document.querySelectorAll(".form-group")

const uname = document.querySelector('#usernameLogin')
const nameValid = document.createElement("p")
uname.addEventListener("input", ()=>{
    if(uname.value.length < 6){
        notify(formGroup[0], "Username length should be more than 5")
    }
    else{
        isNameValid = true
        notify(formGroup[0], "")
    }
})

const pass = document.querySelector('#passwordLogin')
const passValid =document.createElement("p")
pass.addEventListener("input", ()=>{
    if(!(/[A-Z]/.test(pass.value) &&
    /[a-z]/ .test(pass.value) &&
    /[0-9]/.test(pass.value) &&
    pass.value.length > 5)){
        notify(formGroup[1], "Password length should be more than 5 and contain uppercase and lowercase letters")
    }
    else{
        isPassValid = true
        notify(formGroup[1],"")
    }
})


function notify(tag, text){
    let add = true
    for(const child of tag.children){
        if(child.tagName == 'P'){
            add = false
        }
    }
    if(add){
    const notif = document.createElement("p")
    notif.innerHTML = text
    tag.appendChild(notif)
    }
    else{
        tag.children[2].innerHTML = text
    }
}


window.addEventListener('DOMContentLoaded', async () => {
    user = await axios.get(`https://webprojass4-default-rtdb.firebaseio.com/users.json`)
})

loginBtn.addEventListener('click', async (e)=>{
    e.preventDefault();

    
    const incorrectLOgin = document.createElement("p")

    const username = document.querySelector('#usernameLogin').value;
    const password = document.querySelector('#passwordLogin').value;
    let isCorrect = false
    if(isNameValid && isPassValid){
        const userData = user.data;
        Object.keys(userData).map((key)=>{
        if (userData[key].login ==  username && userData[key].password == password) {
            isCorrect = true
            console.log('nice');
            localStorage.setItem('key', key);
            location.reload();
            notify(formGroup[1],"")
        }})
        if(!isCorrect){
            notify(formGroup[1], "Вы не правильно ввели пароль или имя пользователя")
            
        }
    }
    else{
        alert("Форма не валидна");
    }
})

