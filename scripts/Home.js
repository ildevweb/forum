let creatpostbtn =  document.getElementById("create-post");
 creatpostbtn.addEventListener("click", async  function(e) {
    e.preventDefault()
    let  Title  = document.getElementById("title").value;  let Content = document.getElementById("content").value;
    
        try {
             let  response  =  await fetch("/make-post", {
                 method: "POST",
                 headers: { "Content-Type": "application/json" },
                 body: JSON.stringify({
                     title: Title,
                     content : Content
                 })
             })
            if (response.ok) {
               window.location.href = "/home";
            }else {
                let  error  = document.getElementById("error");
                const errorMessage = await response.text()
                console.log(errorMessage)
                error.innerHTML = "Something went wrong ! "+errorMessage;
                error.style.display = "block";
            }
        }catch(Eroor){
            let  error  = document.querySelector("#error");
                error.innerHTML =  Eroor;
                error.style.display = "block";
                console.log(Eroor)
        }
 })


