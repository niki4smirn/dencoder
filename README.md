# Dencoder - Simple Video Storing and Streaming Service

Dencoder is a learning project, a video storing, and streaming service built using Golang, Amazon S3, PostgreSQL, and HTMX.

## Features

- **Video Upload:** Dencoder allows users to easily upload videos to the platform. Uploaded videos are securely stored on Amazon S3, ensuring scalability and reliability.

- **Video Streaming:** Users can watch videos through a seamless streaming experience. Dencoder uses HTMX to provide a dynamic and fast-loading video player interface.

## Getting Started

Follow these steps to set up and run Dencoder:

1. **Clone the Repository:** 
   ```bash
   git clone https://github.com/niki4smirn/dencoder.git
   cd dencoder
   ```

2. **Install Dependencies:** 
   Make sure you have Go, PostgreSQL, and Amazon S3 credentials set up.
   
3. **Configuration:**
   Configure your application by updating the `config.yml` file.

4. **Database Setup:**
   Create the necessary PostgreSQL database table by running the provided SQL script. You can find it in the `database` directory.

5. **Run:**
   Run the Go application (don't forget to set all env variables):
   ```bash
   go run ./...
   ```
## Screenshots

### Minimalistic main page :)
![image](https://github.com/niki4smirn/dencoder/assets/66160046/17f0ac0b-529e-489a-85a1-67e566eb7376)

### The screenshot showcases charming features, including the ability to display adorable cats on your screen 🐱
![image](https://github.com/niki4smirn/dencoder/assets/66160046/fd3a360e-285e-46ce-b990-f78e56385853)

