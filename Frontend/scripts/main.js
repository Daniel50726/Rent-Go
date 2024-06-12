document.addEventListener("DOMContentLoaded", function() {
    const searchButton = document.getElementById("searchButton");

    searchButton.addEventListener("click", function() {
        window.location.href = "./resultadoBusqueda.html";
    });
});