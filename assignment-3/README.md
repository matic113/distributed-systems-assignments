Solution for the section-3 assignment

Docker Hub image
- https://hub.docker.com/r/f4knighty/go-app-section-3

Quick start (build & run locally)

1. Build the image:

```powershell
docker build -t go-app-section-3 .
```

2. Run the container and mount a local directory for persisted users:

```powershell
docker run -d -p 8080:8080 -v ${PWD}/saved_users:/app/saved_users go-app-section-3
```

If you prefer not to use Docker you can run the app directly from the `src` folder:

```powershell
cd src
go run .
```

API

- POST /users
   - Create a new user. Request body (JSON): {"Name":"Alice","Age":30,"Address":{...}}
   - Returns created user with assigned ID.

- GET /users?id=<id>
   - Retrieve user by ID.

Notes

- The application stores users as JSON files in `saved_users/`.
- The directory `saved_users` is created automatically when the app saves a user.