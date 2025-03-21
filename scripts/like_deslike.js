
var button = document.querySelectorAll(".like-icon");
var deslike_button = document.querySelectorAll(".deslike-icon");


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
            const article = document.getElementById(`article_post_${postId}`);
            article.classList.toggle("liked");

            const likeCountElement = document.getElementById(`like-count-${postId}`);
            const deslikeCountElement = document.getElementById(`deslike-count-${postId}`);
            deslikeCountElement.innerText = data.newDeslikeCount;
            likeCountElement.innerText = data.newLikeCount;
        } else {
            alert('Error liking the post!');
        }
    })
    .catch(error => {
        console.error('Error:', error);
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
            const article = document.getElementById(`article_post_${postId}`);
            article.classList.remove("liked");

            const likeCountElement = document.getElementById(`like-count-${postId}`);
            const deslikeCountElement = document.getElementById(`deslike-count-${postId}`);
            deslikeCountElement.innerText = data.newDeslikeCount;
            likeCountElement.innerText = data.newLikeCount;
        } else {
            alert('Error liking the post!');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
