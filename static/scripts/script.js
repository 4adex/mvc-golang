async function logout() {
    const response = await fetch(`/logout`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    });

    const data = await response.json();

    if (data.redirect) {
        window.location.href = data.redirect;
    }
}

// alert("Connected");

async function HandleGet(path) {
    const response = await fetch(path, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });

    if (response.ok){
        window.location.href = path;
    } else {
        const data = await response.json();
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    }  
}

async function HandlePost(path) {
    const response = await fetch(path, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    });

    const data = await response.json();
    if (data.redirect) {
        window.location.href = data.redirect;
    }
}

async function UpdateBook(id) {
    const response = await fetch('/admin/update/'+id, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });

    if (response.ok){
        window.location.href = '/admin/update/'+id;
    } else {
        const data = await response.json();
        if (data.redirect) {
            window.location.href = data.redirect;
        }
    }  
}


async function DeleteBook(id) {
    // alert("Aaaaaa");
    const response = await fetch('/admin/delete/'+id, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    });

    const data = await response.json();
    if (data.redirect) {
        window.location.href = data.redirect;
    }
    
}

