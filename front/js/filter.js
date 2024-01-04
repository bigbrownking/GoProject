const dates = document.getElementById("filter")
let arrDates = []
setTimeout(async() => {
    const response = await axios.get("https://webprojass4-default-rtdb.firebaseio.com/posts.json")
    const data = response.data
    console.log(!arrDates.includes(data[0].data));
    
    const child =  document.createElement("option");
    child.innerHTML = "Показать все"
    child.className = "all"
    dates.appendChild(child)


    Object.keys(data).map((key) => {
        if(!arrDates.includes(data[key].data)){
            const child =  document.createElement("option");
            arrDates.push(data[key].data)
            child.innerHTML = data[key].data
            child.className = data[key].data
            dates.appendChild(child)
        }
    })
})

dates.addEventListener("change", filter)

function filter(e){
    console.log(e.target.value);
    let projects = document.getElementById("projects").childNodes
    for(let i = 0; i < projects.length; i++){
        console.log(projects[i].getAttribute('data'));
        if(e.target.value == "Показать все"){
            projects[i].style = "display: "
        }
        else{
            if(projects[i].getAttribute('data') !== e.target.value){
                projects[i].style = "display: none"
            }
            else{
                projects[i].style = "display: "
            }
        }
    }
}