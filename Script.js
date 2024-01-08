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

async function getData(){
    let data
        try{
            await axios.get("http://localhost:8080/admin/all").then(response => {
                
                data = response.data
                console.log(response.data);
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
            showUsers(data)
        }
        catch(error){
            console.error("Произошла ошибка при отправке запроса", error);
        }
}

function showUsers(data){
    if(data){
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
            td.innerText = user[key]
            td.style.border = "1px solid #000"
            td.style.padding = "3px"
            tr.appendChild(td)
          })
          //_________________delete user______________________
          var td = document.createElement('button');
          td.innerText = "Delete"
          td.id = user["Id"]
          td.addEventListener("click",()=>{
            let item = document.getElementById(`${user["Id"]}`)
            deleteUser(item.id)
          })
          td.style.border = "1px solid #000"
          td.style.padding = "3px"
          tr.appendChild(td)
          //__________________update user__________________
          tr.id = user["Id"]
          tbdy.appendChild(tr);
        }
        table.appendChild(tbdy);

            
        }

}

async function findUser(){
    let id = document.querySelector(".userId")
    let data
    try{
        console.log( id.value);
        await axios.post("http://localhost:8080/admin", {"action": "filter", "id": id.value}).then(response => {
            data = response.data
            console.log(response.data);
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
        showUsers([data])
    }
    catch(error){
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

async function deleteUser(id){
    try{
        console.log( id);
        await axios.post("http://localhost:8080/admin", {"action": "delete", "id": id}).then(response => {
            data = response.data
            console.log(response.data);
        })
        .catch(error => {
            console.error('Ошибка:', error);
        });
        if(data.status == 200){
            alert(data.message + " " + id)
            getData()
        }
    }
    catch(error){
        console.error("Произошла ошибка при отправке запроса", error);
    }
}

function leaveToBack(){
    var domain = window.location.origin;
    window.location.assign(domain);
}