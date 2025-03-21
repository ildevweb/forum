
let cmtbtns = document.querySelectorAll("#comment-button");
cmtbtns.forEach(button => {
    button.addEventListener("click", (e) => {
        e.preventDefault()
        addComment(e)
    })
})


async function  addComment(e) {
    let  form  =  e.target.closest("form");
    let  error  =   form.querySelector("#CommentError");
    

    fetch(`/add-comment`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            post_id:   form.querySelector("input[name=post_id]").value,
            comment: form.querySelector("textarea[name=comment]").value,
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            let comments_container = document.getElementById(`comment-container-${data.comment.PostID}`);
            comments_container.appendChild(createCommentsSection(data.comment));
            let none_comment = document.getElementById(`none-comment-${data.comment.PostID}`);
            if (none_comment != null) {
                none_comment.style.display = "none";
            }
        } else {
            alert('Error');
        }
    })
    .catch(err => {
        console.log(err)
        error.innerHTML = "Cannot Add Comment ! (Field is empty Or something went wrong or  Length  is Too long)";
        error.style.display = "block";
    });
}



function createCommentsSection(comment) {
    const commentDiv = document.createElement('div');
            commentDiv.classList.add('comment');

            const commentContent = document.createElement('div');
            commentContent.classList.add('comment-content');
            commentContent.id = "comments-content";

            // Create comment text
            const commentText = document.createElement('p');
            commentText.classList.add('comment-text');
            commentText.textContent = comment.Content;

            // Create comment author
            const commentAuthor = document.createElement('p');
            commentAuthor.classList.add('comment-author');
            commentAuthor.innerHTML = `<strong>${comment.Username}</strong>`;

            // Create comment date
            const commentDate = document.createElement('p');
            commentDate.classList.add('comment-date');
            commentDate.innerHTML = `<small>${comment.CreatedAt}</small>`;

            // Create the reaction icons
            const reactionIcons = document.createElement('div');
            reactionIcons.classList.add('reaction-icons');

            const likeIcon = document.createElement('i');
            likeIcon.classList.add('fas', 'fa-thumbs-up', 'like-button');
            likeIcon.onclick = () => handlecommentLike(comment.ID);

            const likeCount = document.createElement('span');
            likeCount.id = `like-num-${comment.ID}`;
            likeCount.classList.add('like-num');
            likeCount.textContent = comment.Likes;

            const dislikeIcon = document.createElement('i');
            dislikeIcon.classList.add('fas', 'fa-thumbs-down', 'deslike-button');
            dislikeIcon.onclick = () => handlecommentDeslike(comment.ID);

            const dislikeCount = document.createElement('span');
            dislikeCount.id = `deslike-num-${comment.ID}`;
            dislikeCount.classList.add('deslike-num');
            dislikeCount.textContent = comment.Deslikes;

            // Append icons to the reaction container
            reactionIcons.appendChild(likeIcon);
            reactionIcons.appendChild(likeCount);
            reactionIcons.appendChild(dislikeIcon);
            reactionIcons.appendChild(dislikeCount);

            // Append comment details to the comment content
            commentContent.appendChild(commentText);
            commentContent.appendChild(commentAuthor);
            commentContent.appendChild(commentDate);
            commentContent.appendChild(reactionIcons);

            // Append the comment content to the comment div
            commentDiv.appendChild(commentContent);


    return commentDiv;
}