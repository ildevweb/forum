
function loginClicked() {
    const loginWrapper = document.getElementById("login-wrapper");
    const modalOverlay = document.getElementById("login-modal-overlay");

    loginWrapper.classList.toggle("display_none");

    const registerWrapper = document.getElementById("register-wrapper");
    const registermodalOverlay = document.getElementById("register-modal-overlay");
    registerWrapper.classList.add("display_none");
    registermodalOverlay.style.display = 'none';

    if (loginWrapper.classList.contains("display_none")) {
        modalOverlay.style.display = 'none';
    } else {
        modalOverlay.style.display = 'flex';
        loginWrapper.classList.add("show");
    }
}

function registerClicked() {
    const registerWrapper = document.getElementById("register-wrapper");
    const modalOverlay = document.getElementById("register-modal-overlay");

    registerWrapper.classList.toggle("display_none");

    const loginWrapper = document.getElementById("login-wrapper");
    const loginmodalOverlay = document.getElementById("login-modal-overlay");
    loginWrapper.classList.add("display_none");
    loginmodalOverlay.style.display = 'none';

    if (registerWrapper.classList.contains("display_none")) {
        modalOverlay.style.display = 'none';
    } else {
        modalOverlay.style.display = 'flex';
        registerWrapper.classList.add("show");
    }
}

document.getElementById("login-modal-overlay").addEventListener("click", function(event) {
    if (event.target === this) {
        const loginWrapper = document.getElementById("login-wrapper");
        const modalOverlay = document.getElementById("login-modal-overlay");

        loginWrapper.classList.add("display_none");
        modalOverlay.style.display = 'none';
    }
});

document.getElementById("register-modal-overlay").addEventListener("click", function(event) {
    if (event.target === this) {
        const registerWrapper = document.getElementById("register-wrapper");
        const modalOverlay = document.getElementById("register-modal-overlay");

        registerWrapper.classList.add("display_none");
        modalOverlay.style.display = 'none';
    }
});


let creatpostbtn =  document.getElementById("create-post");
 creatpostbtn.addEventListener("click", async  function(e) {
    e.preventDefault()
    let  Title  = document.getElementById("title").value;
    let Content = document.getElementById("content").value;
    let Category = getSelectedCategories()
    
    
    fetch(`/make-post`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            title: Title,
            content : Content,
            category : Category
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            error.innerHTML = "Post  Has  Created Sucessfuly !"
            error.style.color =  "green"
            error.style.display = "block";
            
            let posts_section = document.querySelector(".posts-section");
            let firstElement = posts_section.firstChild ;
            let article = createPost(data.post);
            posts_section.insertBefore(article, firstElement);
            //posts_section.appendChild(article);

            let no_posts = document.getElementById(`no-posts`);
            if (no_posts != null) {
                no_posts.style.display = "none";
            }
            //handle adding comments for the added posts without reloading page
            let buttons = document.querySelectorAll("#comment-button-created");
            buttons.forEach((button) => {
                button.addEventListener("click", (e) => {
                    e.preventDefault()
                    addComment(e)
                })
            })


            //handle filter if adding posts without reloading the page
            const posts = document.querySelectorAll('article[id^="article_post_"]');
            let  radiosbtn  =  document.querySelectorAll("input[name='filter']");
            radiosbtn.forEach(radio => {
                radio.addEventListener("change", () => {
                    let  filter  =  radio.value;
                    filterPosts(filter, posts);
                })
            })
            
            search_field(posts);
        } else {
            alert('Error');
        }
    })
    .catch(err => {
        let  error  = document.querySelector("#error");
        error.innerHTML =  "Cannot create post ! (Field is empty Or something went wrong)";
        error.style.color =  "red"
        error.style.display = "block";
        //console.log(err)
    });
 })







 function createPost(post) {
    // Create the article element
    const article = document.createElement('article');
    article.id = `article_post_${post.ID}`;
    article.classList.add("created");

    // Create the category section
    const categorySection = document.createElement('h4');
    categorySection.innerHTML = `Categories: <span class="posts_categories" style="font-size:15px;">${post.Category}</span>`;

    // Create the title
    const title = document.createElement('h3');
    title.textContent = post.Title;

    // Create the content paragraph
    const content = document.createElement('p');
    content.textContent = post.Content;
    content.id = "content-post";

    // Create the posted info
    const postedInfo = document.createElement('p');
    postedInfo.innerHTML = `<small>Posted by <strong>${post.Username}</strong> on ${post.CreatedAt}</small>`;

    // Create the reaction section
    const reactionIcons = document.createElement('div');
    reactionIcons.classList.add('reaction-icons');

    const likeIcon = document.createElement('i');
    likeIcon.classList.add('fas', 'fa-thumbs-up', 'like-icon');
    likeIcon.onclick = () => handleLike(post.ID);

    const likeCount = document.createElement('span');
    likeCount.id = `like-count-${post.ID}`;
    likeCount.classList.add('like-count');
    likeCount.textContent = post.Likes;

    const dislikeIcon = document.createElement('i');
    dislikeIcon.classList.add('fas', 'fa-thumbs-down', 'deslike-icon');
    dislikeIcon.onclick = () => handleDeslike(post.ID);

    const dislikeCount = document.createElement('span');
    dislikeCount.id = `deslike-count-${post.ID}`;
    dislikeCount.classList.add('deslike-count');
    dislikeCount.textContent = post.Deslikes;

    // Append elements to the reactionIcons div
    reactionIcons.appendChild(likeIcon);
    reactionIcons.appendChild(likeCount);
    reactionIcons.appendChild(dislikeIcon);
    reactionIcons.appendChild(dislikeCount);

    // Append all sections to the article
    article.appendChild(categorySection);
    article.appendChild(title);
    article.appendChild(content);
    article.appendChild(postedInfo);
    article.appendChild(reactionIcons);

    // Create comments section if there are comments
    const commentsSection = document.createElement('div');
    commentsSection.classList.add('comments-section');

    const commentsContainer = document.createElement('section');
    commentsContainer.classList.add('comments-container-created');
    commentsContainer.setAttribute("id", `comment-container-${post.ID}`)
    if (post.Comments && post.Comments.length > 0) {
        post.Comments.forEach(comment => {
            const commentDiv = document.createElement('div');
            commentDiv.classList.add('comment');

            const commentContent = document.createElement('div');
            commentContent.classList.add('comment-content');

            const commentText = document.createElement('p');
            commentText.classList.add('comment-text');
            commentText.textContent = comment.Content;

            const commentAuthor = document.createElement('p');
            commentAuthor.classList.add('comment-author');
            commentAuthor.innerHTML = `<strong>${comment.Created_by}</strong>`;

            const commentDate = document.createElement('p');
            commentDate.classList.add('comment-date');
            commentDate.innerHTML = `<small>${comment.CreatedAt}</small>`;

            const commentReactionIcons = document.createElement('div');
            commentReactionIcons.classList.add('reaction-icons');

            const commentLikeIcon = document.createElement('i');
            commentLikeIcon.classList.add('fas', 'fa-thumbs-up', 'like-button');
            commentLikeIcon.onclick = () => handlecommentLike(comment.CommentID);

            const commentLikeCount = document.createElement('span');
            commentLikeCount.id = `like-num-${comment.CommentID}`;
            commentLikeCount.classList.add('like-num');
            commentLikeCount.textContent = comment.Likes;

            const commentDislikeIcon = document.createElement('i');
            commentDislikeIcon.classList.add('fas', 'fa-thumbs-down', 'deslike-button');
            commentDislikeIcon.onclick = () => handlecommentDeslike(comment.CommentID);

            const commentDislikeCount = document.createElement('span');
            commentDislikeCount.id = `deslike-num-${comment.CommentID}`;
            commentDislikeCount.classList.add('deslike-num');
            commentDislikeCount.textContent = comment.Deslikes;

            // Append reaction icons
            commentReactionIcons.appendChild(commentLikeIcon);
            commentReactionIcons.appendChild(commentLikeCount);
            commentReactionIcons.appendChild(commentDislikeIcon);
            commentReactionIcons.appendChild(commentDislikeCount);

            // Append all comment content to the comment div
            commentContent.appendChild(commentText);
            commentContent.appendChild(commentAuthor);
            commentContent.appendChild(commentDate);
            commentContent.appendChild(commentReactionIcons);

            // Append the comment to the comments container
            commentDiv.appendChild(commentContent);
            commentsContainer.appendChild(commentDiv);
        });

        
    } else {
        const noCommentsMessage = document.createElement('p');
        noCommentsMessage.textContent = 'No comments yet.';
        commentsSection.appendChild(noCommentsMessage);
    }
    commentsSection.appendChild(commentsContainer);

    // Create comment form
    const commentForm = document.createElement('form');
    commentForm.action = '/add-comment';
    commentForm.method = 'POST';
    commentForm.classList.add('comment-form');

    const hiddenInput = document.createElement('input');
    hiddenInput.type = 'hidden';
    hiddenInput.name = 'post_id';
    hiddenInput.value = post.ID;

    const commentField = document.createElement('textarea');
    commentField.classList.add('coomment-field');
    commentField.name = 'comment';
    commentField.placeholder = 'Write a comment...';
    commentField.required = true;

    const submitButton = document.createElement('button');
    submitButton.type = 'submit';
    submitButton.id = 'comment-button-created';


    const paperPlaneIcon = document.createElement('i');
    paperPlaneIcon.classList.add('fa', 'fa-paper-plane');
    submitButton.appendChild(paperPlaneIcon);

    const commentError = document.createElement('span');
    commentError.id = 'CommentError';

    // Append form elements
    commentForm.appendChild(hiddenInput);
    commentForm.appendChild(commentField);
    commentForm.appendChild(submitButton);
    commentForm.appendChild(commentError);

    // Append the comments section and form to the article
    commentsSection.appendChild(commentForm);
    article.appendChild(commentsSection);


    return article;
}


