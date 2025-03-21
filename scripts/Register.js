var Allfields = {
    Email : document.getElementById("register_email"),
    Username : document.getElementById("register_username"),
    Password : document.getElementById("register_password")
}
let  register_submit  = document.getElementById("register");
register_submit.addEventListener("click", async  function(e) {
    e.preventDefault()
    if (Allfields.Email.value === "" || Allfields.Username.value === "" || Allfields.Password.value === ""){
        let error = document.getElementById("register_errorMessage");
        error.innerHTML = "Please fill all the fields !";
        error.style.display = "block";
    } else {
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
                window.location.href = "/";
            } else {
                const errorMessage = await response.text()
                let error = document.getElementById("register_errorMessage");
                error.innerHTML = "Something went wrong ! " + errorMessage;
                error.style.display = "block";
            }
        }catch(Eroor){
            let  error  = document.querySelector("#register_errorMessage");
            error.innerHTML =  Eroor;
            error.style.display = "block";
        }
    }
})