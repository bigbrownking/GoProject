document.addEventListener('DOMContentLoaded', async function() {
    const key = localStorage.getItem('key');
    console.log(key);

    if(key){
        const user = await axios.get(`https://webprojass4-default-rtdb.firebaseio.com/users/${key}.json`)
        const userData = user.data;

        const loginLink = document.querySelector('#loginLink');
        loginLink.innerHTML=`<a class="nav-link" href="../user/user.html">${userData.userName}</a>`

        localStorage.setItem('userInfo', key);
    }
});