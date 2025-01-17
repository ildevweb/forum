
var button = document.querySelectorAll(".like-icon");
var deslike_button = document.querySelectorAll(".deslike-icon");


function like_saved() {
    button.forEach((btn, index) => {
      if (localStorage.getItem(`like-button${index}`) === 'true') {
        btn.classList.add('active');
      } else {
        btn.classList.remove("active");
      }
    });
}

function deslike_saved() {
    deslike_button.forEach((btn, index) => {
        if (localStorage.getItem(`deslike-button${index}`) === 'true') {
          btn.classList.add('active');
        } else {
          btn.classList.remove("active");
        }
    });
}

function handleLike(postId) {
    //button[postId-1].classList.toggle("active");
    fetch(`/like-post/${postId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            const likeCountElement = document.getElementById(`like-count-${postId}`);
            const deslikeCountElement = document.getElementById(`deslike-count-${postId}`);
            deslikeCountElement.innerText = data.newDeslikeCount;
            likeCountElement.innerText = data.newLikeCount;

            if (data.colored_like == true) {
                button[button.length-postId].classList.add("active");
                //localStorage.setItem('buttonClicked', button.length-postId);
                localStorage.setItem(`like-button${button.length-postId}`, 'true');
                localStorage.setItem(`deslike-button${button.length-postId}`, 'false');
                deslike_button[deslike_button.length-postId].classList.remove("active");
            } else {
                button[button.length-postId].classList.remove("active");
                localStorage.setItem(`like-button${button.length-postId}`, 'false');
            }
        } else {
            alert('Error liking the post!');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('An error occurred!');
    });
}

function handleDeslike(postId) {
    //deslike_button[postId-1].classList.toggle("active");

    fetch(`/deslike-post/${postId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            const likeCountElement = document.getElementById(`like-count-${postId}`);
            const deslikeCountElement = document.getElementById(`deslike-count-${postId}`);
            deslikeCountElement.innerText = data.newDeslikeCount;
            likeCountElement.innerText = data.newLikeCount;


            if (data.colored_deslike == true) {
                deslike_button[deslike_button.length-postId].classList.add("active");
                localStorage.setItem(`deslike-button${button.length-postId}`, 'true');
                localStorage.setItem(`like-button${button.length-postId}`, 'false');
                button[button.length-postId].classList.remove("active");
            } else {
                deslike_button[deslike_button.length-postId].classList.remove("active");
                localStorage.setItem(`deslike-button${button.length-postId}`, 'false');
            }
        } else {
            alert('Error liking the post!');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('An error occurred!');
    });
}

like_saved();
deslike_saved();