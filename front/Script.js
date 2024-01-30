var logins = []
document.querySelector('.sign_up').addEventListener('submit', async function(e) {
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
            "login": username,
            "password": password,
            "email": email,
            "number": phone,
            "address": address
        };
        console.log(data);
        try{
            await axios.post("http://localhost:8080/register", data).then(response => {
                data = response.data
                console.log(response.data);
            })
            .catch((error) => {
                console.error('Ошибка:', error);
            });
            showUsers(data)
        }
        catch(error){
            console.error("Произошла ошибка при отправке запроса", error);
        }
        console.log(domain + "/AdminPage");
        window.location.assign(domain + "/AdminPage");
        return false
    }
    return false
})

async function getData() {
    let data
    try {
        await axios.get("http://localhost:8080/admin/all").then(response => {

            data = response.data
            console.log(response.data);
        })
            .catch(error => {
                console.error('Ошибка:', error);
            });
        showUsers(data)
    }
    catch (error) {
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

function showUsers(data) {
    if (data) {
        let table = document.querySelector(".tableUsers")
        while (table.firstChild) {
            table.removeChild(table.firstChild);
        }
        table.style.border = "1px solid #000"
        var tr = document.createElement('tr');
        Object.keys(data[0]).map(key => {
            var td = document.createElement('td');
            td.innerText = key
            td.style.border = "1px solid #000"
            tr.appendChild(td)
        })
        var tbdy = document.createElement('tbody');
        tbdy.appendChild(tr);
        for (let user of data) {
            var tr = document.createElement('tr');
            Object.keys(user).map((key, id) => {
                var td = document.createElement('td');
                if(key === "login"){
                    logins.push(user[key])
                }
                td.innerText = user[key]
                td.style.border = "1px solid #000"
                td.style.padding = "3px"
                tr.appendChild(td)
            })
            //_________________delete user______________________
            var td = document.createElement('button');
            td.innerText = "Delete"
            td.id = user["Id"]
            td.addEventListener("click", () => {
                let item = document.getElementById(`${user["Id"]}`)
                deleteUser(item.id)
            })
            td.style.border = "1px solid #000"
            td.style.padding = "3px"
            tr.appendChild(td)
            //__________________update user__________________
            var tu = document.createElement('button');
            tu.innerText = "Update";
            tu.id = user["Id"]
            tu.addEventListener("click", () => {
                updateUser(user["Id"]);
            });
            tu.style.border = "1px solid #000";
            tu.style.padding = "3px";
            tr.appendChild(tu);

            tbdy.appendChild(tr);
        }
        table.appendChild(tbdy);
    }

}
async function findUserByLogin(login){
    let data
    try {
        await axios.post("http://localhost:8080/admin", { "action": "filterLogin", "login": login}).then(response => {
            data = response.data
            console.log(response.data);
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
        let loginInput = document.querySelector(".userInputLogin")
        loginInput.value = null
        let select = document.querySelector(".findUserByLogin")
        while (select.firstChild) {
            select.removeChild(select.firstChild);
        }
        showUsers(data)
    }
    catch (error) {
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

function searchUser(){
    let login = document.querySelector(".userInputLogin").value
    let select = document.querySelector(".findUserByLogin")
    while (select.firstChild) {
        select.removeChild(select.firstChild);
    }
    for(let log of logins){
        let option = document.createElement("button")
        console.log(log);
        if(log.toLowerCase().includes(login.toLowerCase())){
            option.innerText = log
            option.id = log
            option.addEventListener("click", () => {
                findUserByLogin(log);
            });
            select.appendChild(option)
        }
    }
}

async function changeSort(){
    let sortOrder = document.querySelector("#sortLogin").value
    console.log(sortOrder);
    if(sortOrder !== ""){
        try {
            const response = await axios.post("http://localhost:8080/admin", {
                "action": "sort",
                "Login": "login",
                "Id": sortOrder,
            });
            showUsers(response.data);
        } catch (error) {
            console.error("Произошла ошибка при отправке запроса", error);
        }
    }
    
}

async function findUser() {
    let id = document.querySelector(".userId")
    let data
    try {
        console.log(id.value);
        await axios.post("http://localhost:8080/admin", { "action": "filter", "id": id.value }).then(response => {
            data = response.data
            console.log(response.data);
        })
            .catch(error => {
                console.error('Ошибка:', error);
            });
        showUsers([data])
    }
    catch (error) {
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

async function deleteUser(id) {
    try {
        console.log(id);
        await axios.post("http://localhost:8080/admin", { "action": "delete", "id": id }).then(response => {
            data = response.data
            console.log(response.data);
        })
            .catch(error => {
                console.error('Ошибка:', error);
            });
        if (data.status == 200) {
            alert(data.message + " " + id)
            getData()
        }
    }
    catch (error) {
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

async function updateUser(id) {
    let data;
    try {
        await axios.post("http://localhost:8080/admin", {"action": "filter", "id": id}).then(response => {
            data = response.data
            console.log(response.data);
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
        console.log(id);
        var updatedUsername = prompt("Введите обновленное имя пользователя:");
        var updatedPassword = prompt("Введите обновленный пароль:");
        var updatedEmail = prompt("Введите обновленный адрес электронной почты:");
        var updatedPhone = prompt("Введите обновленный номер телефона:");
        var updatedAddress = prompt("Введите обновленный адрес:");
        console.log(updatedUsername.trim());
        if (updatedUsername == null || updatedUsername.trim() == ""){
            updatedUsername = data["login"]
        }
        if (updatedPassword === null || updatedPassword.trim() === ""){
            updatedPassword = data["password"]
        }
        if (updatedEmail === null || updatedEmail.trim() === ""){
            updatedEmail = data["email"]
        }
        if (updatedPhone === null || updatedPhone.trim() === ""){
            updatedPhone = data["number"]
        }
        if (updatedAddress === null || updatedAddress.trim() === ""){
            updatedAddress = data["address"]
        }
        console.log({ "action": "update", "id": id, "login": updatedUsername, "password": updatedPassword, "email": updatedEmail, "number": updatedPhone, "address": updatedAddress});
        await axios.post("http://localhost:8080/admin", { "action": "update", "id": id, "login": updatedUsername, "password": updatedPassword, "email": updatedEmail, "number": updatedPhone, "address": updatedAddress})
            .then(response => {
                data = response.data;
                console.log(response.data);
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
            getData()
    } catch (error) {
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

function leaveToBack() {
    var domain = window.location.origin;
    window.location.assign(domain);
}