document.addEventListener('DOMContentLoaded', (event) => {
    const detailsContainer = document.getElementById('detailsContainer');
    if (detailsContainer) {
        const selectedCar = JSON.parse(localStorage.getItem('selectedCar'));

        if (selectedCar) {
            const detailsContent = document.createElement('div');
            detailsContent.setAttribute('class', 'details-content');

            const image = document.createElement('img');
            image.setAttribute('src', selectedCar.images.CarImage);
            image.setAttribute('alt', selectedCar.name);
            image.setAttribute('class', 'carImage');

            const description = document.createElement('div');
            description.setAttribute('class', 'details-description');

            const name = document.createElement('h2');
            name.textContent = selectedCar.name;

            const detailsList = document.createElement('ul');

            const detailsItems = [
                { label: 'Descripción', value: selectedCar.biography.DescriptionCar },
                { label: 'Marca', value: selectedCar.brand },
                { label: 'Modelo', value: selectedCar.model },
                { label: 'Color', value: selectedCar.color },
                { label: 'Tipo', value: selectedCar.type },
                { label: 'Asistencia', value: selectedCar.assistance ? 'Sí' : 'No' },
                { label: 'Seguro', value: selectedCar.insurance ? 'Sí' : 'No' },
                { label: 'Asiento para bebé', value: selectedCar.babySeat ? 'Sí' : 'No' },
                { label: 'Transmisión', value: selectedCar.transmission },
                { label: 'Asiento de lujo', value: selectedCar.luxurySeat ? 'Sí' : 'No' },
                { label: 'Combustible', value: selectedCar.fuel },
                { label: 'Precio', value: selectedCar.price },
            ];

            detailsItems.forEach(item => {
                const li = document.createElement('li');
                li.innerHTML = `<span>${item.label}:</span> ${item.value}`;
                detailsList.appendChild(li);
            });

            const deleteButton = document.createElement('button');
            deleteButton.setAttribute('class', 'deleteButton');
            deleteButton.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill="red" d="M3 6l3 18h12l3-18H3zm18-2h-4l-1-1H8L7 4H3V2h18v2z"/></svg>';

            deleteButton.addEventListener('click', () => {
                fetch('http://localhost:8080/api/cars/toggleReservation', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ carName: selectedCar.name }),
                })
                    .then(response => response.json())
                    .then(data => {
                        console.log('Car Reservation Status Changed:', data);

                        // Remove the car from the local storage
                        let reservations = JSON.parse(localStorage.getItem('userReservations')) || [];
                        reservations = reservations.filter(reservedCar => reservedCar.name !== selectedCar.name);
                        localStorage.setItem('userReservations', JSON.stringify(reservations));

                        // Redirect back to search.html after deletion
                        window.location.href = 'search.html';
                    })
                    .catch(error => console.error('Error:', error));
            });

            description.appendChild(name);
            description.appendChild(detailsList);

            detailsContent.appendChild(image);
            detailsContent.appendChild(description);
            detailsContent.appendChild(deleteButton);

            detailsContainer.appendChild(detailsContent);
        } else {
            console.log('No car details found.');
        }
    }
});
