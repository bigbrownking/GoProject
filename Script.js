// let table = document.querySelector(".table")


function submitForm(e) {
    e.preventDefault();
    var domain = window.location.origin;
    console.log(domain + "/AdminPAge");
    window.location.assign(domain + "/AdminPage");
}