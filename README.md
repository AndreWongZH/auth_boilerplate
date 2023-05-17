# Auth Boilerplate

This is a boilerplate project for implementing user authentication using a Next.js frontend, Go backend, and MongoDB database.

Here is a link to my [blog post](https://blog.wongandre.com/authentication-with-email-and-code-verification) on the explanation.

<img src="https://github.com/AndreWongZH/auth_boilerplate/blob/main/images/verify_page.png?raw=true" alt="verify_page" width="500px"/>
<br />
<img src="https://github.com/AndreWongZH/auth_boilerplate/blob/main/images/register_filled.png?raw=true" alt="register_page" width="350px"/>
<br />
<img src="https://github.com/AndreWongZH/auth_boilerplate/blob/main/images/verify_email.png?raw=true" alt="verify_email" />

<br />

## Features

- User registration with email verification
- User login with password hashing
- MongoDB for data storage
- Next.js for frontend development
- Go for backend development

## Prerequisites

Before running the application, ensure you have the following prerequisites installed on your machine:

- Node.js and NPM
- Go
- MongoDB

## Getting Started

To get started with the Auth Boilerplate, follow these steps:

1. Configuration
The configuration for the application can be found in the .env file in the project root directory. Modify the environment variables in this file to add in your private keys.

1. Clone the repository:
    ```
    git clone https://github.com/AndreWongZH/auth_boilerplate.git
    ```

1. Navigate to the project directory:
    ```
    cd auth_boilerplate
    ```

1. Set up the frontend:
    ```
    // Navigate to the frontend directory:
    cd frontend

    // Install dependencies:
    yarn

    // Start the development server:
    yarn dev
    ```

1. Set up the backend:
    ```
    // Navigate to the backend directory:
    cd backend

    Install dependencies:
    // go mod download

    // Start the backend server
    go run .
    ```

1. Access the application:
Open your browser and visit http://localhost:3000 to access the frontend.
