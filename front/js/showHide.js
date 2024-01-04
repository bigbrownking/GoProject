let arr = null

function showFinished(){
    var finished = document.getElementById("projects").childNodes
    var buttons = document.getElementById("sort").childNodes
    for(let i = 0; i < buttons.length; i++){
        if(buttons[i].id !== "finished"){
            buttons[i].style = "color: #052c65;"
        }
        else{
            buttons[i].style = "background-color: #ffac38;  color: white;"
        }
    }
    for(let i = 0; i < finished.length; i++){
        if(finished[i].id !== "finished"){
            finished[i].style = "display: none"
        }
        else{
            finished[i].style = "display: "
        }
    }
}
function showInProcess(){
    var finished = document.getElementById("projects").childNodes
    for(let i = 0; i < finished.length; i++){
        if(finished[i].id !== "inProcess"){
            finished[i].style = "display: none"
        }
        else{
            finished[i].style = "display: "
        }
    }
    var buttons = document.getElementById("sort").childNodes
    for(let i = 0; i < buttons.length; i++){
        if(buttons[i].id !== "inProcess"){
            buttons[i].style = "color: #052c65;"
        }
        else{
            buttons[i].style = "background-color: #ffac38;  color: white;"
        }
    }
}
function showDesignProcess(){
    var finished = document.getElementById("projects").childNodes
    for(let i = 0; i < finished.length; i++){
        if(finished[i].id !== "designProcess"){
            finished[i].style = "display: none"
        }
        else{
            finished[i].style = "display: "
        }
    }
    var buttons = document.getElementById("sort").childNodes
    for(let i = 0; i < buttons.length; i++){
        if(buttons[i].id !== "designProcess"){
            buttons[i].style = "color: #052c65;"
        }
        else{
            buttons[i].style = "background-color: #ffac38;  color: white;"
        }
    }
}
function showAll(){
    var finished = document.getElementById("projects").childNodes
    for(let i = 0; i < finished.length; i++){
            finished[i].style = "display: "
    }
    var buttons = document.getElementById("sort").childNodes
    for(let i = 0; i < buttons.length; i++){
        if(buttons[i].id !== "all"){
            buttons[i].style = "color: #052c65;"
        }
        else{
            buttons[i].style = "background-color: #ffac38;  color: white;"
        }
    }
}