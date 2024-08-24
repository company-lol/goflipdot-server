# 📟 goflipdot-server

Welcome to the `goflipdot-server` repository! This project provides a Go-based server to control a Flipdot sign. 

## 📋 Summary of Project

`goflipdot-server` is a Golang application designed to interact with Flipdot signs. It provides a RESTful interface to update the Flipdot display via an HTTP API.

## 🚀 How to Use

1. **Clone the Repository**
    ```bash
    git clone https://github.com/harperreed/goflipdot-server.git
    cd goflipdot-server
    ```

2. **Configure Server**

   - Copy the example configuration file and update it according to your setup:
     ```bash
     cp config.ini.example config.ini
     ```
   - Update `config.ini` with your specific Flipdot sign and server configurations.

3. **Build & Run**
   - Build the application:
     ```bash
     make build
     ```
   - Run the application:
     ```bash
     make run
     ```

4. **Using Docker**
   - Build the Docker image:
     ```bash
     docker-compose build
     ```
   - Run the Docker container:
     ```bash
     docker-compose up
     ```

5. **Access the API**
   - You can now interact with the API:
     - To post a new array of dots to the Flipdot sign:
       ```bash
       curl -X POST http://localhost:8080/api/dots -d '[[true,false,...],...]' -H "Content-Type: application/json"
       ```
     - To view documentation:
       Visit `http://localhost:8080/documentation`

## ⚙️ Tech Info

- **Dockerfile**: Defines the multi-stage Docker build process.
- **Makefile**: Provides various commands to build, run, test, and clean the project.
- **Configuration**: 
  - `config.ini`: Main configuration file for the server.
  - `config.ini.example`: Example configuration file.
- **Dependencies**:
  - Defined in `go.mod` and `go.sum`.
- **Main Application**:
  - `main.go`: The entry point of the application containing server setup and handlers.

The directory structure is as follows:

```plaintext
goflipdot-server/
├── Dockerfile
├── Makefile
├── config.ini
├── config.ini.example
├── docker-compose.yaml
├── go.mod
├── go.sum
├── main.go
```

Dive in and start customizing the `goflipdot-server` for your own unique Flipdot display setups! If you have questions, feel free to ask. Happy coding! 🚀
