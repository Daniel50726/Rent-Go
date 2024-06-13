COMANDO PARA EJECUTAR: docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
IP: docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' some-postgres
CADENA DE CONEXION: dsn := "user=postgres password=mysecretpassword dbname=postgres sslmode=disable host=172.17.0.2 port=5432"