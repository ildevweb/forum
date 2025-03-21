
let span = document.querySelector(".close")

let btn = document.getElementById('new_post')


let modal = document.getElementById('createPostModal')

btn.addEventListener('click', () => {
    modal.style.display = "block"
})

span.addEventListener('click', ()=>{
    modal.style.display = "none"
})

