const submitBtn = document.querySelector('#submit');
const phoneInput = document.querySelector('#phone');
const commentInput = document.querySelector('#comment')


submitBtn.addEventListener('click', validatorForm);

function validatorForm(event){
    const regex = /[a-zA-Z]/;
    const phoneValue = phoneInput.value;
    const commentValue = commentInput.value;
    
    if(!phoneValue || !commentValue){
        alert('Введите номер или комментарий');
    }
    if(regex.test(phoneValue)){
        alert('Номер должен состоят из чисел')
    }
    alert('Мы успешно отправили вам смс!')
    event.preventDefault();
}