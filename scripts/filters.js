let  radiosbtn  =  document.querySelectorAll("input[name='filter']");
radiosbtn.forEach(radio => {    
    radio.addEventListener("change", () => {
        let  filter  =  radio.value;
        filterPosts(filter);
    })
})



const posts = document.querySelectorAll('article[id^="article_post_"]');
function  filterPosts(filter) {
    
    console.log(posts)
    posts.forEach(post => {
        if (filter === 'all') {
            post.style.display = 'block';
        } else if (filter === 'created') {
            if (post.classList.contains('created')) {
                post.style.display = 'block';
            } else {
                        post.style.display = 'none';
            }
        } else if (filter === 'liked') {
            if (post.classList.contains('liked')) {
                post.style.display = 'block';
            } else {
                post.style.display = 'none';
            }
        }
    });
}




let  searchfield = document.querySelector(".search-input");
searchfield.addEventListener("input", () => {
    console.log(searchfield.value)
    let  searchTerm = searchfield.value.toLowerCase();
    posts.forEach(post => {

        if (post.querySelector('.posts_categories').textContent.toLowerCase().includes(searchTerm)) {
            post.style.display = 'block';
        } else {
            post.style.display = 'none';
        }
    });
    newposts = document.querySelectorAll('article[id^="article_post_"][style="display: block;"]');
    let Eroor = document.querySelector("#ErrorPosts") 
    if(newposts.length === 0) {
     Eroor.innerHTML = `<p>No posts to display. Be the first to <a href="/make-post">add a post</a>!</p>`
     Eroor.style.display = "block";
    }else {
        Eroor.style.display = "none";
    }
})