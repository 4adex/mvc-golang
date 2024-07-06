// const delay = 100;
console.log("Script loaded");
async function logout() {
    HandlePost('/logout');
}

// alert("Connected");

async function HandleGet(path) {
    // alert('Foo');
    const response = await fetch(path, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        
        }, credentials: 'include'
    });

    if (response.ok){
        setTimeout(() => {
            window.location.href = path;
        }, 100);
    } else {
        const data = await response.json();
        if (data.flashMessage && data.flashType) {
            // Store flash message in localStorage
            localStorage.setItem('flashMessage', data.flashMessage);
            localStorage.setItem('flashType', data.flashType);
        }
        if (data.redirect) {
            console.log("Redirecting to:", data.redirect);
            window.location.href = data.redirect;
        }
    }  
}

async function HandlePost(path) {
    console.log("HandlePost called with path:", path);
    try {
        const response = await fetch(path, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include'
        });

        const data = await response.json();
        console.log("Response received:", data);
        
        if (data.flashMessage && data.flashType) {
            // Store flash message in localStorage
            localStorage.setItem('flashMessage', data.flashMessage);
            localStorage.setItem('flashType', data.flashType);
        }
        
        if (data.redirect) {
            console.log("Redirecting to:", data.redirect);
            window.location.href = data.redirect;
        }
    } catch (error) {
        console.error("Error in HandlePost:", error);
    }
}

async function UpdateBook(id) {
    HandleGet('/admin/update/'+id);
}


async function DeleteBook(id) {
    HandlePost('/admin/delete/'+id);
}

function displayFlashMessage() {
    // alert("Youtube");
    const flashMessage = localStorage.getItem('flashMessage');
    const flashType = localStorage.getItem('flashType');
    
    if (flashMessage && flashType) {
        // Display the flash message (adjust this based on your HTML structure)
        flashElement = document.getElementById('message-div');
        flashElement.textContent = flashMessage;
        flashElement.className = `alert alert-${flashType}`;
        // document.body.insertBefore(flashElement, document.body.firstChild);
        
        // Clear the flash message from localStorage
        localStorage.removeItem('flashMessage');
        localStorage.removeItem('flashType');
        
        // Optionally, remove the flash message after a few seconds
        setTimeout(() => {
            flashElement.remove();
        }, 5000);
    }
}

