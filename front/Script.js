var logins = []
let allUsers = []
var currnetLogs = []
let maxItems = 0
let activePage = 1
var domain = window.location.href;

var currentUser = null


async function isLogin(){
    let fullpage =document.querySelector("#fullpage")
    await axios.get("isLogin").then(res =>{
        if(res.data.login !== ""){
            currentUser = res.data
            fullpage.style.display = "inline"
                let adminPage =document.querySelector("#adminPage")
                if(currentUser.isAdmin){
                    adminPage.style.display = "inline"

                }else{
                    adminPage.style.display = "none"

                }
            if(domain.includes("profile")){
    
                let table = document.querySelector(".profileInformation")
                if(table){
                    table.innerHTML = `
                    <tr>
                        <th scope="row">Login</th>
                        <td>${currentUser.login}</td>
                    </tr>
                    <tr>
                        <th scope="row">Email</th>
                        <td>${currentUser.email}</td>
                    </tr>
                    <tr>
                        <th scope="row">Address</th>
                        <td>${currentUser.address}</td>
                    </tr>
                    <tr>
                        <th scope="row">Phone number</th>
                        <td>${currentUser.number}</td>
                    </tr>`
                }
            }
            return true
        }
        fullpage.style.display = "none"
        return false
    })
}
isLogin()
if(domain.includes("main")){
    let posts = document.querySelector("#posts")
    if(posts){
        setTimeout(async () => {
            let data = null
            await axios.get("/getPosts",).then(res=>{
                data = res.data
            })
            while (posts.firstChild) {
                posts.removeChild(posts.firstChild);
            }
            console.log(data);
            if(data){
                data.forEach(element => {
                    let card = document.createElement("div")
                    card.className = "col"
                    card.id = element.Id
                    card.innerHTML = 
                    `<div style="" class="card">
                    <img style="width: 150px; height: 150px" src="${element.img}" class="card-img-top" alt="...">
                    <div class="card-body">
                    <h5 class="card-title"><p>${element.Name}</p></h5>
                    <p class="card-text">${element.Decs}</p>
                    <button class="card-text btn btn-success" id="addToCard${element.Id}">Add to cart</button>

                    </div>
                    <div class="card-footer">
                    <small class="text-muted">Price : ${element.Price}</small>
                    </div>
                    </div>
                    `
                    posts.appendChild(card)
                    document.querySelector(`#addToCard${element.Id}`).addEventListener("click",async ()=>{
                        await axios.post("/Cart", {userId: currentUser.Id, card: element})
                    })
                });
            }
        })

    }
}
const reg = document.querySelector('.sign_up')
if(reg){
    reg.addEventListener('submit', async function(e) {
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
                "address": address,
                "isAdmin": false
            };
            console.log(data);
            try{
                await axios.post("/register", data).then(response => {
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
            console.log(domain + "/loginPage");
            window.location.assign(domain + "/loginPage");
            return false
        }
        return false
    })
}

async function mailingTextSend(){
    let mailingText = document.querySelector("#mailingText")
    let mailingSubject = document.querySelector("#mailingSubject")
    if(mailingText && mailingText.value.trim().length > 0 && mailingSubject && mailingSubject.value.trim().length > 0){
        allUsers.forEach(async user =>{
            
            await axios.post("/mailing", {email: user.email, text: mailingText.value, subject: mailingSubject.value}).catch(err =>{
                console.log(err);
            })
        })
    }
}

async function getCards(){
    let posts = document.querySelector("#posts")
    if(posts){
            let data = currentUser.cards
            while (posts.firstChild) {
                posts.removeChild(posts.firstChild);
            }
            if(data){
                data.forEach(element => {
                    let card = document.createElement("div")
                    card.className = "col"
                    card.id = element.Id
                    card.innerHTML = 
                    `<div style="" class="card">
                    <img style="width: 150px; height: 150px" src="${element.img}" class="card-img-top" alt="...">
                    <div class="card-body">
                    <h5 class="card-title"><p>${element.Name}</p></h5>
                    <p class="card-text">${element.Decs}</p>
                    </div>
                    <div class="card-footer">
                    <small class="text-muted">Price : ${element.Price}</small>
                    </div>
                    </div>
                    `
                    posts.appendChild(card)
                });
            }

    }
}

async function getData() {
    let data
    try {
        await axios.get("/admin/all").then(response => {

            data = response.data
            console.log(response.data);
        })
            .catch(error => {
                console.error('Ошибка:', error);    
            });
        maxItems = data.length
        allUsers = data
        currnetLogs = data
        showUsers(data)
        changePage(1)
    }
    catch (error) {
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

function showUsers(dataOld) {
    if (dataOld) {
        logins = []
        allUsers = dataOld
        let data = currnetLogs
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
        let pages = document.querySelector(".page")
        while (pages.firstChild) {
            pages.removeChild(pages.firstChild);
        }
        let maxPage =Math.ceil( maxItems/5)
        for(let i = 1; i <= maxPage; i++){
            let page = document.createElement("button")
            page.innerText = i
            page.id = i + "page"
            page.addEventListener("click", () => {
                changePage(i)
            })
            pages.appendChild(page)
        }
    }
}

async function changePage(item){
    let maxPage =Math.ceil( maxItems/5)
    for(let i = 1; i <= maxPage; i++){
        let page = document.getElementById(`${i}page`)
        page.style.backgroundColor  = "black"
        if(item === i){
            page.style.backgroundColor  = "red"
        }
    }
    console.log(allUsers);
    currnetLogs = allUsers.slice(5*(item-1), 5*item)
    console.log(allUsers);
    console.log(currnetLogs);
    showUsers(allUsers)

}
async function findUserByLogin(login){
    let data
    try {
        await axios.post("/admin", { "action": "filterLogin", "login": login}).then(response => {
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
        changePage(1)
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
            const response = await axios.post("/admin", {
                "action": "sort",
                "Login": "login",
                "Id": sortOrder,
            });
            showUsers(response.data);
            changePage(1)
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
        await axios.post("/admin", { "action": "filter", "id": id.value }).then(response => {
            data = response.data
            console.log(response.data);
        })
            .catch(error => {
                console.error('Ошибка:', error);
            });
        showUsers([data])
        changePage(1)
    }
    catch (error) {
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

async function deleteUser(id) {
    try {
        console.log(id);
        await axios.post("/admin", { "action": "delete", "id": id }).then(response => {
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
        await axios.post("/admin", {"action": "filter", "id": id}).then(response => {
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
        await axios.post("/admin", { "action": "update", "id": id, "login": updatedUsername, "password": updatedPassword, "email": updatedEmail, "number": updatedPhone, "address": updatedAddress})
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

function redirect(page){
    console.log(123);
    let domain = window.location.origin;
    window.location.assign(domain + '/' + page);
}

async function logout(){
    await axios.get("/logOut").then(()=>{   
        window.location.assign(window.location.origin + '/loginPage');
    })
}