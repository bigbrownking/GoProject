async function createUser(e){
    e.preventDefault();
    const name = document.querySelector('#name')
    const surname = document.querySelector('#surname')
    const username = document.querySelector('#signupUsername')
    const email = document.querySelector('#email');
    const password = document.querySelector('#password')
    const formGroup = document.querySelectorAll(".form-group")

    let isNameValid = false
    let isSurnameValid  = false
    let isUsernameValid = false
    let isEmailValid = false
    let isPassValid = false

    if(!name.value && !surname.value && !username.value && !password.value && !email.value){
        alert('Fill the form')
    }else{
    if(!email.value.match(/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/)){
        notify(formGroup[5], "Invalid email address")
    }
    else{
        isEmailValid = true
        notify(formGroup[5], "")}

    if(!(/[A-Z]/.test(password.value) && /[a-z]/.test(password.value) && /[0-9]/.test(password.value) && password.value.length > 5)){
        notify(formGroup[6], 'Password length should be more than 5, contain numbers, uppercase and lowercase letters')
    }else{
        isPassValid = true
        notify(formGroup[6], "")
    }

    if(name.value.length < 6){
        notify(formGroup[2], 'Name length should be more than 5')
    }else{
        isNameValid = true
        notify(formGroup[2], "")
    }
    
    if(surname.value.length < 6){
        notify(formGroup[3], 'Surname length should be more than 5')
    }else{
        isSurnameValid = true
        notify(formGroup[3], "")
    }
    
    if(username.value.length < 6){
        notify(formGroup[4], 'Username length should be more than 5')
    }
    else{
        isUsernameValid = true
        notify(formGroup[4], "")
    }
    if(isEmailValid && isNameValid && isPassValid && isSurnameValid && isUsernameValid){
        await axios.post(
                "https://webprojass4-default-rtdb.firebaseio.com/users.json",
                {
                    "userName":name.value,
                    "surName": surname.value,
                    "login": username.value,
                    "password":password.value,
                    "email":email.value,
                }
            )
            location.reload();       
    }}
}

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


const addUser = document.querySelector('#createBtn');
addUser.addEventListener('click', createUser);