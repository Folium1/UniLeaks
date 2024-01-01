# Anonymous Task Sharing Platform

This repository contains a project developed in Go, utilizing MySQL for data storage and JWT for authentication. The project aims to provide a platform where students can anonymously share files related to university tasks. The files are stored on Google Drive.

## Table of Contents
- [Anonymous Task Sharing Platform](#anonymous-task-sharing-platform)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Usage](#usage)

## Features
- Anonymous file sharing: Students can upload files without revealing their identity.
- Task organization: Files can be organized based on university subjects or specific tasks.
- Authentication: JWT is used for user authentication, ensuring secure access to the platform.
- Role-based access control: Different user roles (e.g., admin, student) can be defined to manage access permissions.
- Google Drive integration: Files are stored on Google Drive for efficient and scalable storage.

## Prerequisites
To run this project, you need to have the following prerequisites installed:

- Go (1.20 or higher)
- MySQL (5.7 or higher)
- Git
- Docker (24.0 or higher)

## Installation
1. Clone the repository:
```bash
   git clone https://github.com/Folium1/UniLeaks.git
```

2. Run the following commands:
```bash
    make generate-certs
    make build
```

3. Set up the database:
   - Create a new MySQL database for the project.

4. Configure Google Drive API:
   - Go to the [Google Cloud Console](https://console.cloud.google.com/).
   - Create a new project or select an existing one.
   - Enable the Google Drive API for the project.
   - Create credentials (OAuth 2.0 client ID) and download the JSON file.
   - Rename the downloaded JSON file to `client_credentials.json` and place it in the project root directory.

## Configuration
Before running the project, you need to configure the necessary environment variables. Create a `.env` file in the project root directory and populate it with the following variables:

```plaintext
   MYSQL=""
   SALT=""
   # api for checking files for viruses
   CLOUD_MERSIVE_API=""
   GOOGLE_CLIENT_ID=""
   GOOGLE_SECRET=""
   # used to check if the user is from a particular university
   MAIL_DOMAIN=""
```

## Usage
To start the application in docker, run the following command:
```bash
    docker build -t leaks .
    # sometimes can fail, because app can run quicker then db container, still don't know how to fix it
    # but it will try to reconnect to db later, you don't need to to extra moves
    docker-compose up --build
```

The application will be accessible at `https://localhost:8080/leaks`.