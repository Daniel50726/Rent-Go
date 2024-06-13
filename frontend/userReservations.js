document.addEventListener('DOMContentLoaded', (event) => {
    const reservationContainer = document.getElementById('reservationContainer');
    if (reservationContainer) {
        const reservations = JSON.parse(localStorage.getItem('userReservations'));
        if (reservations && reservations.length > 0) {
            reservations.forEach(car => {
                const tarjetContainer = document.createElement("div");
                tarjetContainer.setAttribute("class", "tarjetContainer");

                const tarjetName = document.createElement("h2");
                tarjetName.setAttribute("class", "tarjetName");
                tarjetName.textContent = car.name;

                const imageContainer = document.createElement("div");
                imageContainer.setAttribute("class", "imageContainer");

                const image = document.createElement("img");
                image.setAttribute("src", car.images.CarImage);
                image.setAttribute("class", "tarjetImage");

                const price = document.createElement("p");
                price.setAttribute("class", "tarjetPrice");
                price.textContent = `Precio: ${car.price}`;

                tarjetContainer.appendChild(tarjetName);
                imageContainer.appendChild(image); // Agregar imagen al contenedor de imagen
                tarjetContainer.appendChild(imageContainer); // Agregar contenedor de imagen al contenedor principal
                tarjetContainer.appendChild(price);

                tarjetContainer.addEventListener('click', () => {
                    // Redirect to reservationInfo.html and pass car details via localStorage
                    localStorage.setItem('selectedCar', JSON.stringify(car));
                    window.location.href = 'reservationInfo.html';
                });

                reservationContainer.appendChild(tarjetContainer);
            });
        } else {
            console.log('No reservations found.');
        }
    }
});
