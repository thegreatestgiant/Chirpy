# Chirpy

Chirpy is a monolith REST API written in Go that acts as a messaging platform.

## Boot.dev course context

I built this for Boot.dev's
[Learn HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang)
course (the "Chirpy" project).

| | |
|---|---|
| Provided by Boot.dev | The course lessons, project spec, and the CLI test harness used to grade each chapter. There was no starter application. |
| Written by me | All 22 commits and ~1,030 Go lines: routing, JSON handlers, JSON-file persistence, users and chirps, JWT access/refresh tokens with revocation, API-key webhook, metrics, CORS, and this README. |
| Not provided | No starter app and no committed test suite. The course tests run through the Boot.dev CLI. |

This is a learning project, not a production service.

# The frontend
The frontend can be accessed under /apps will be compiled from the root directory

# Admin
Accessable under /admin 

## Metrics
You can see the number of times a Chirp has been read only as a GET request

# The backend
Everything here will be accessed under /api

## GET Requests

### Health
Accessable under /healthz
This will tell you if the server is up

### /Reset
This will reset the metrics

### /Chirps
This will return all the chirps on the server

### /Chirps/{chirpID}
This will return the chirp with that specific id

### /Reset
This will reset your database

## POST Requests

### /Chirps
This will create a chirp

### /Users
This will create a user

### /Login
This will return a Refresh and Access Token

### /Refresh
Takes a refresh-token as input, it returns a new access token

### Revoke
Takes a Refresh Token as input and revokes it

## Webhook

### Polka
This is a POST Request and takes a body like this
```json
{
    "event": "user.upgraded",
    "data": {
        "user_id": 3
    }
}
```

## PUT Requests

### /Users
This updates the user

## DELETE Requests

### /Chirps/{chirpID}
deletes the Chirp under chirpID
