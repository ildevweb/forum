let  submit  = document.getElementById("login");
submit.addEventListener("click", async  function(e) {
    e.preventDefault()
    let username = document.getElementById("username").value;  let password = document.getElementById("password").value;
   
    if (username === "" || password === ""){
    let    error = document.getElementById("errorMessage");
    error.innerHTML = "Please fill all the fields ! ";
    error.style.display = "block";
    }else { 
        try {
             let  response  =  await fetch("/login", {
                 method: "POST",
                 headers: { "Content-Type": "application/json" },
                 body: JSON.stringify({
                     username: username,
                     password : password
                 })
             })
            if (response.ok) {
                window.location.href = "/home";
            }else {
                let    error = document.getElementById("errorMessage");
                const errorMessage = await response.text()
                error.innerHTML = "Something went wrong ! "+errorMessage;
                error.style.display = "block";
            }
        }catch(Eroor){
            let  error  = document.querySelector("#errorMessage");
                error.innerHTML =  Eroor;
                error.style.display = "block";
        }
    }
   
    });