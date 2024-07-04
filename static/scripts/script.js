const delay = 100;

async function logout() {
    const response = await fetch(`/logout`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }, credentials: 'include'
    });

    const data = await response.json();

    if (data.redirect) {
        setTimeout(() => {
            window.location.href = data.redirect;
        }, delay);
    }
}

// alert("Connected");

async function HandleGet(path) {
    const response = await fetch(path, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        
        }, credentials: 'include'
    });

    if (response.ok){
        setTimeout(() => {
            window.location.href = path;
        }, delay);
    } else {
        const data = await response.json();
        if (data.redirect) {
            setTimeout(() => {
                window.location.href = data.redirect;
            }, delay);
        }
    }  
}

async function HandlePost(path) {
    const response = await fetch(path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include'  // This ensures cookies are sent with the request
    });

    const data = await response.json();
    if (data.redirect) {
        // Add a small delay before redirect
        setTimeout(() => {
            window.location.href = data.redirect;
        }, delay);
    }
}

async function UpdateBook(id) {
    const response = await fetch('/admin/update/'+id, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }, credentials: "include"
    });

    if (response.ok){
        window.location.href = '/admin/update/'+id;
    } else {
        const data = await response.json();
        if (data.redirect) {
            setTimeout(() => {
                window.location.href = data.redirect;
            }, delay);
        }
    }  
}


async function DeleteBook(id) {
    // alert("Aaaaaa");
    const response = await fetch('/admin/delete/'+id, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }, credentials: "include"
    });

    const data = await response.json();
    if (data.redirect) {
        setTimeout(() => {
            window.location.href = data.redirect;
        }, delay);
    }
    
}

