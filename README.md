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

GET /user/{id} - Get users (no auth required): takes a userID and return the user info 

POST /users - Register (no auth required): takes a username, password and optional string for
mobile_token. Registers the user and authenticates the client as the newly created user

POST /login - Sign in/authenticate: takes a username and password, and authenticates the
user

PUT /users - Update User (must be signed in as the user): updates password and/or
mobile_token of the user

DELETE /users - Delete User (must be signed in as the user): deletes the user

```

## Room Management: 
```
The server is able to handle creating conference rooms.

Room has: name (non-unique), guid, host user, participants (users) in the room, and a capacity limit.

```

##Routes:
```

POST /rooms - Create a room (signed in as a user): creates a room hosted by the current user, with an optional capacity limit. Default is 5.

GET /rooms/{guid} - given a room guid, gets information about a room

POST /rooms/{guid}/users - Join room(signed in as a user): joins the room as the current user

DELETE /rooms/{guid}/users - leave room(signed in as a user): leaves the room as the current user

PUT /rooms/{guid} - Change host (must be signin as the host): changes the host of the user from the
current user to another user, which should be a participant

GET /users/{id}/rooms - Search for the rooms that a user is in: given a username, returns a list of rooms
that the user is in.
 
```
