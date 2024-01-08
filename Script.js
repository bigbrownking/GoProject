// let table = document.querySelector(".table")


function submitForm(e) {
    e.preventDefault();
    var domain = window.location.origin;
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
    var confirm_password = document.getElementById("confirm_password").value;
    var email = document.getElementById("email").value;
    var phone = document.getElementById("phone").value;
    var address = document.getElementById("address").value;
    if (password != confirm_password){
        alert('Passwords do not match');
    }
    else{
        var data = {
            "username": username,
            "password": password,
            "email": email,
            "phone": phone,
            "address": address
        };
        fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
        .then(response => response.text())
        .then(data => {
            console.log('Успешно:', data);
        })
        .catch((error) => {
            console.error('Ошибка:', error);
        });
        console.log(domain + "/AdminPage");
        window.location.assign(domain + "/AdminPage");
    }
}

document.getElementById("log_in").onclick = function() {
    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;
    var data = {
        "login": username,
        "passwd": password
    };
    fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    })
    .then(response => {
        if (response.ok) {
            console.log('Успешный вход');
        } else {
            console.error('Неудачная попытка входа');
        }
    })
    .catch(error => {
        console.error('Ошибка:', error);
    });
    var domain = window.location.origin;
    window.location.assign(domain + "/AdminPage");
}

function send(){
    setTimeout(async()=>{
        try{
            await axios.post("http://localhost:8080/register", {
	            "login": "EsimgaliKh",
                "password":"fkmvflkvm",
                "email": "fkvmldfkmv",
                "number":"ldmokfsdmvk"
        });
        console.log("Успешно");
        }
        catch(error){
            console.error("Произошла ошибка при отправке запроса", error);
        }
    });
}

function leaveToBack(){
    var domain = window.location.origin;
    window.location.assign(domain);
}