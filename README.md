# video_server
It's a basic server having following functionalities:

## User Management: 
```
The server is able to register and authenticate users.
User has: username, password, and an optional mobile_token (string)
```

##Routes:
```
GET /users - Get users (no auth required): returns a list of all users

GET /user/{id} - Get users (no auth required): takes a username and return the user with
matching username

POST /users - Register (no auth required): takes a username, password and optional string for
mobile_token. Registers the user and authenticates the client as the newly created user

POST /login - Sign in/authenticate: takes a username and password, and authenticates the
user

PUT /user/{id} - Update User (must be signed in as the user): updates password and/or
mobile_token of the user

DELETE /user/{id} - Delete User (must be signed in as the user): deletes the user

```
