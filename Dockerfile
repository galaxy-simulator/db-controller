FROM golang:latest

# Copy the source files into the container
COPY . . 
COPY frontend/ frontend/
COPY backend/ frontend/

# Get dependencies
RUN ["go", "get", "git.darknebu.la/GalaxySimulator/structs"]
#RUN ["go", "get", "git.darknebu.la/GalaxySimulator/db-controller/frontend"]
RUN ["go", "get", "git.darknebu.la/GalaxySimulator/db-controller/backend"]
RUN ["go", "get", "github.com/gorilla/mux"]
RUN ["go", "get", "github.com/lib/pq"]

# build an executable
RUN ["go", "build", "-o", "db-controller", "."]

# Start the webserver
ENTRYPOINT ["./db-controller"]
