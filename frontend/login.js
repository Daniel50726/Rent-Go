document.getElementById('loginForm').addEventListener('submit', async function (event) {
    event.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const response = await axios.post('http://localhost:8080/api/login', {
            username: username,
            password: password
        }, {
            headers: {
                'Content-Type': 'application/json'
            }
        });

        document.getElementById('message').textContent = response.data;
        window.location.href = 'search.html';
    } catch (error) {
        if (error.response) {
            document.getElementById('message').textContent = error.response.data;
        } else {
            document.getElementById('message').textContent = 'An error occurred';
        }
    }
});

// Add event listener to the register button
document.querySelector('.registerUser').addEventListener('click', function() {
    window.location.href = 'register.html';
});
