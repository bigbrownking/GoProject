const userKey = localStorage.getItem('userInfo');
console.log(userKey);   

const infoCont = document.querySelector('#info');
const adminLink = document.querySelector('#admin');
window.addEventListener('DOMContentLoaded', async ()=>{
    const user = await axios.get(`https://webprojass4-default-rtdb.firebaseio.com/users/${userKey}.json`)
    const link = document.querySelector('#loginLink');
    const userData = user.data;
    infoCont.innerHTML=`<p>${userData.userName}</p>
    <p>${userData.surName}</p>
    <p>${userData.email}</p>
    <p>${userData.password}</p>`
    link.innerHTML=`<a disabled class="nav-link">${userData.userName}</a>`
    console.log(userData.isAdmin);
    console.log(adminLink);
    if(userData.isAdmin == true){
        adminLink.style.display='block'
    }
})
const logout = document.querySelector("#logoutPage");
logout.addEventListener('click', ()=>{
    localStorage.removeItem('userInfo');
    localStorage.removeItem('key')
})