var modal = document.getElementById("myModal");
var btn = document.getElementById("createBtn");
var span = document.getElementsByClassName("close")[0];

btn.onclick = function() {
  modal.style.display = "block";
}

span.onclick = function() {
  modal.style.display = "none";
}

window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}



async function createUser(e){
    e.preventDefault();
    const isLogin = document.querySelector('#login')
    const isPassword = document.querySelector('#password')
    const isUserName = document.querySelector('#username')
    const isUserSurname = document.querySelector('#surname')
    const email = document.querySelector('#email')
    const isAdmin = document.querySelector('#radioButton');
    await axios.post(
        "https://webprojass4-default-rtdb.firebaseio.com/users.json",
        {
            "userName":isUserName.value,
            "surName": isUserSurname.value,
            "login": isLogin.value,
            "password":isPassword.value,
            "email":email.value,
            "isAdmin":Boolean(isAdmin.value),
            "isLike": false,
            "isLogin":false
        }
    )
    location.reload()
  }

const addUser = document.querySelector('#addUser');
addUser.addEventListener('click', createUser);

let data = {};

window.addEventListener('DOMContentLoaded', async ()=>{
    const response = await axios.get("https://webprojass4-default-rtdb.firebaseio.com/users.json")
    data = response.data;
    console.log(data);
})


const usersContainer = document.querySelector('#users');


setTimeout(() => {
    console.log('some');
    const deleteBtn = document.querySelectorAll('#delete')
    console.log(deleteBtn);

    console.log(usersContainer);

    usersContainer.innerHTML += Object.entries(data).map(([key, value]) => {

        return `
        <tr>
            <td>${value.userName}</td>
            <td>${value.surName}</td>
            <td>${value.login}</td>
            <td>${value.email}</td>
            <th><button class="edit" data-target="${key}">Edit</button></th>
            <th><button class="delete" data-target="${key}">Delete</button></th>
        </tr>`;
    }).join('');
}, 1500);





usersContainer.addEventListener('click', async (event) => {
    
    
    if (event.target.classList.contains('delete')) {
        const userId = event.target.getAttribute('data-target');
        await axios.delete(`https://webprojass4-default-rtdb.firebaseio.com/users/${userId}.json`)
        location.reload()
    } else if (event.target.classList.contains('edit')) {
        const userId = event.target.getAttribute('data-target');

        const user = await axios.get(`https://webprojass4-default-rtdb.firebaseio.com/users/${userId}.json`)
        const userData = user.data;
        console.log(userData);
        displayEditForm(userData);

        const editBtn = document.querySelector('#editBtn');
        editBtn.addEventListener('click', async ()=>{
            await axios.put(
                `https://webprojass4-default-rtdb.firebaseio.com/users/${userId}.json`,  
                {
                    "userName":document.getElementById('editUserName').value,
                    "surName": document.getElementById('editSurName').value,
                    "login": document.getElementById('editLogin').value,
                    "password":document.getElementById('editPassword').value,
                    "email":document.getElementById('editEmail').value,
                    "isAdmin":false,
                    "isLike": false,
                    "isLogin":false
                }
            )
            location.reload()
        })
    }
});


function displayEditForm(userData) {
    document.getElementById('editForm').style.display = 'block';

    var editFormModal = document.getElementById("editForm");
    var close = document.querySelector("#closeBtn");


    close.onclick = function() {
        editFormModal.style.display = "none";
    }
    
    window.onclick = function(event) {
        console.log(event.target);
        if (event.target == editFormModal) {
            editFormModal.style.display = "none";
        }
    }



    document.getElementById('editUserName').value = userData.userName;
    document.getElementById('editSurName').value = userData.surName;
    document.getElementById('editLogin').value = userData.login;
    document.getElementById('editEmail').value = userData.email;
    document.getElementById('editPassword').value = userData.password;

    
}


var modalPost = document.getElementById("modalPost");
var modalPostCloseBtn = document.getElementById("modalPost-close");
const createPost = document.querySelector('#createPost')

createPost.onclick = function() {
    modalPost.style.display = "block";
}
modalPostCloseBtn.onclick = function() {
    modalPost.style.display = "none";
}   
  
  window.onclick = function(event) {
    if (event.target == modalPost) {
        modalPost.style.display = "none";
    }
  }

const createPostApi = document.querySelector('#createPostApi');
createPostApi.addEventListener('click', createPostFunc)

async function createPostFunc(e){
    e.preventDefault();
    const title = document.querySelector('#title').value    
    const desc = document.querySelector('#description').value
    const state = document.getElementById("status").value;
    const img = document.getElementById("img").value
    const currentDate = new Date();
    await axios.post("https://webprojass4-default-rtdb.firebaseio.com/posts.json", {
        data: currentDate.getDate() + "." + (currentDate.getMonth() + 1) + "." + currentDate.getFullYear(),
        discription: desc,
        title: title,
        state: state,
        img: img
    })
    modalPost.style.display = "none";
}
