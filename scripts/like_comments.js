

function handlecommentLike(commentId) {
    fetch(`/like-comment/${commentId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        const like_count = document.getElementById(`like-num-${commentId}`);
        const deslike_count = document.getElementById(`deslike-num-${commentId}`);
        like_count.textContent = data.newLikeCount;
        deslike_count.textContent = data.newDeslikeCount;
    })
    .catch(error => {
        console.error('Error:', error);
        alert('An error occurred!');
    });
}


function handlecommentDeslike(commentId) {
    fetch(`/deslike-comment/${commentId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        const like_count = document.getElementById(`like-num-${commentId}`);
        const deslike_count = document.getElementById(`deslike-num-${commentId}`);
        like_count.textContent = data.newLikeCount;
        deslike_count.textContent = data.newDeslikeCount;
    })
    .catch(error => {
        console.error('Error:', error);
        alert('An error occurred!');
    });
}