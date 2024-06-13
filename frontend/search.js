function searchCar() {
    const searchInput = document.getElementById('searchInput').value;
    const searchCriteria = document.getElementById('searchCriteria').value;

    const requestBody = {};
    requestBody[searchCriteria] = searchInput;

    fetch('http://localhost:8080/api/cars/search', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
    })
        .then(response => response.json())
        .then(data => {
            console.log('Search Results:', data); // Imprimir resultados en consola
            localStorage.setItem('searchResults', JSON.stringify(data));
            window.location.href = 'searchResult.html';
        })
        .catch(error => console.error('Error:', error));
}

document.addEventListener('DOMContentLoaded', (event) => {
    const resultsContainer = document.getElementById('resultsContainer');
    if (resultsContainer) {
        const results = JSON.parse(localStorage.getItem('searchResults'));
        if (results && results.length > 0) {
            results.forEach(car => {
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
                price.textContent = `PRECIO: ${car.price}`;

                const rentButton = document.createElement("button");
                rentButton.setAttribute("class", "rentButton");

                if (car.reservation) {
                    rentButton.textContent = "Bloqueado";
                    rentButton.style.backgroundColor = "red";
                    rentButton.disabled = true;
                } else {
                    rentButton.textContent = "Rentar ya";
                    rentButton.addEventListener('click', () => {
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
                                rentButton.textContent = "Bloqueado";
                                rentButton.style.backgroundColor = "red";
                                rentButton.disabled = true;

                                // Save the car in localStorage and redirect
                                let reservations = JSON.parse(localStorage.getItem('userReservations')) || [];
                                reservations.push(car);
                                localStorage.setItem('userReservations', JSON.stringify(reservations));
                                window.location.href = 'userReservations.html';
                            })
                            .catch(error => console.error('Error:', error));
                    });
                }

                imageContainer.appendChild(image);
                tarjetContainer.appendChild(tarjetName);
                tarjetContainer.appendChild(imageContainer);
                tarjetContainer.appendChild(price);
                tarjetContainer.appendChild(rentButton);

                resultsContainer.appendChild(tarjetContainer);
            });
        } else {
            console.log('No search results found.');
        }
    }
});

