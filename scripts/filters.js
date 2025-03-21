let  radiosbtn  =  document.querySelectorAll("input[name='filter']");
radiosbtn.forEach(radio => {    
    radio.addEventListener("change", () => {
        let  filter  =  radio.value;
        filterPosts(filter);
    })
})



const posts = document.querySelectorAll('article[id^="article_post_"]');
function  filterPosts(filter, psts) {
    if (psts == null) psts = posts ;
    psts.forEach(post => {
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
function search_field(pst) {
    if (pst == null) pst = posts ;
    searchfield.addEventListener("input", () => {
        let  searchTerm = searchfield.value.toLowerCase();
        pst.forEach(post => {
            if (post.querySelector('.posts_categories').textContent.toLowerCase().includes(searchTerm)) {
                post.style.display = 'block';
            } else {
                post.style.display = 'none';
            }
        });
        newposts = document.querySelectorAll('article[id^="article_post_"][style="display: block;"]');
        let Eroor = document.querySelector("#ErrorPosts") 
        if(newposts.length === 0) {
        Eroor.innerHTML = `<p>No posts to display. Be the first to add post </p>`
        Eroor.style.display = "block";
        }else {
            Eroor.style.display = "none";
        }
    })
}


search_field();