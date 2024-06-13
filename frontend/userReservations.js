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

                const deleteButton = document.createElement("button");
                deleteButton.setAttribute("class", "deleteButton");
                deleteButton.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill="red" d="M3 6l3 18h12l3-18H3zm18-2h-4l-1-1H8L7 4H3V2h18v2z"/></svg>';

                deleteButton.addEventListener('click', () => {
                    fetch('http://localhost:8080/api/cars/toggleReservation', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({ carName: car.name }),
                    })
                        .then(response => response.json())
                        .then(data => {
                            console.log('Car Reservation Status Changed:', data);

                            // Remove the car from the local storage
                            let reservations = JSON.parse(localStorage.getItem('userReservations')) || [];
                            reservations = reservations.filter(reservedCar => reservedCar.name !== car.name);
                            localStorage.setItem('userReservations', JSON.stringify(reservations));

                            // Remove the card from the frontend
                            reservationContainer.removeChild(tarjetContainer);
                        })
                        .catch(error => console.error('Error:', error));
                });

                imageContainer.appendChild(image);
                tarjetContainer.appendChild(tarjetName);
                tarjetContainer.appendChild(imageContainer);
                tarjetContainer.appendChild(price);
                tarjetContainer.appendChild(deleteButton);

                reservationContainer.appendChild(tarjetContainer);
            });
        } else {
            console.log('No reservations found.');
        }
    }
});
