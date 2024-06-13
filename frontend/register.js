document.getElementById('registerForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;
    const messageDiv = document.getElementById('message');

    if (password !== confirmPassword) {
        messageDiv.textContent = 'Las contraseñas no coinciden.';
        messageDiv.style.color = 'red';
        return;
    }

    try {
        const response = await axios.post('http://localhost:8080/api/register', {
            username: username,
            password: password
        }, {
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (response.status === 200) {
            messageDiv.textContent = `Usuario ${username} registrado correctamente.`;
            messageDiv.style.color = 'green';

            // Redireccionar al usuario a la página de login después de un registro exitoso
            setTimeout(() => {
                window.location.href = 'login.html';
            }, 2000); // Retraso de 2 segundos antes de redirigir (opcional)
        }
    } catch (error) {
        messageDiv.textContent = 'Error al registrar el usuario.';
        messageDiv.style.color = 'red';
        console.error('Error:', error);
    }
});
