var   Allfields  = {
    Email : document.getElementById("email"),
    Username : document.getElementById("username"),
    Password : document.getElementById("password")
}
let  submit  = document.getElementById("register");
submit.addEventListener("click", async  function(e) {
    e.preventDefault()
    if (Allfields.Email.value === "" || Allfields.Username.value === "" || Allfields.Password.value === ""){
        let    error = document.getElementById("errorMessage");
        error.innerHTML = "Please fill all the fields ! ";
        error.style.display = "block";
    }else { 
        try {
             let  response  =  await fetch("/register", {
                 method: "POST",
                 headers: { "Content-Type": "application/json" },
                 body: JSON.stringify({
                     email: Allfields.Email.value,
                     username: Allfields.Username.value,
                     password : Allfields.Password.value
                 })
             })
            if (response.ok) {
                window.location.href = "/login";
            }else {
                const errorMessage = await response.text()
                let    error = document.getElementById("errorMessage");
                error.innerHTML = "Something went wrong ! " + errorMessage;
                error.style.display = "block";
            }
        }catch(Eroor){
            let  error  = document.querySelector("#errorMessage");
                error.innerHTML =  Eroor;
                error.style.display = "block";
        }
    }
})