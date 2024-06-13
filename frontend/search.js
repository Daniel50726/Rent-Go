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
                tarjetName.textContent = car.name; // Access the name property directly

                const imageContainer = document.createElement("div");
                imageContainer.setAttribute("class", "imageContainer");

                const image = document.createElement("img");
                image.setAttribute("src", car.images.CarImage); // Access the CarImage property directly
                image.setAttribute("class", "tarjetImage");

                const price = document.createElement("p");
                price.setAttribute("class", "tarjetPrice");
                price.textContent = `Price: ${car.price}`;

                const rentButton = document.createElement("button");
                rentButton.setAttribute("class", "rentButton");
                rentButton.textContent = "Rentar ya";

                imageContainer.appendChild(image);
                tarjetContainer.appendChild(tarjetName);
                tarjetContainer.appendChild(imageContainer);
                tarjetContainer.appendChild(price);
                tarjetContainer.appendChild(rentButton);

                resultsContainer.appendChild(tarjetContainer);
            });
        } else {
            console.log('No search results found.'); // Log if no results found
        }
    }
});

