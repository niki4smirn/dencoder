# Dencoder - Simple Video Storing and Streaming Service

Dencoder is a learning project, a video storing, and streaming service built using Golang, Amazon S3, PostgreSQL, and HTMX.

## Features

- **Video Upload:** Dencoder allows users to easily upload videos to the platform. Uploaded videos are securely stored on Amazon S3, ensuring scalability and reliability.

- **Video Streaming:** Users can watch videos through a seamless streaming experience. Dencoder uses HTMX to provide a dynamic and fast-loading video player interface.

- **Database Integration:** Video metadata and user information are stored in a PostgreSQL database, making it easy to manage and organize your video content.

## Getting Started

Follow these steps to set up and run Dencoder on your local machine:

1. **Clone the Repository:** 
   ```bash
   git clone https://github.com/your-username/dencoder.git
   cd dencoder
   ```

2. **Install Dependencies:** 
   Make sure you have Go, PostgreSQL, and Amazon S3 credentials set up.
   
3. **Configuration:**
   Configure your application by updating the `config.yml` file.

4. **Database Setup:**
   Create the necessary PostgreSQL database table by running the provided SQL script. You can find it in the `database` directory.

5. **Run:**
   Run the Go application (don't forget to set your config path and postgresql credentials):
   ```bash
   go run CONFIG_PATH=... PGX_PASS=... PGX_USER=...
   ```

6. **Access Dencoder:**
   Open your web browser and navigate to `http://localhost:8080` to access the Dencoder web interface.
