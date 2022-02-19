# rcodingtakehome

This code base is for a take home coding challenge for a job application

## Goals

A simple server with one path ("foo") that supports three methods: GET, POST, and DELETE.

### GET

`http://127.0.0.1:8080/foo/{id}`

Gets a foo record if found

Response Body:

```
{"name": "{name}"
 "id": "{id}"
```

### POST

Adds a new foo record
`http://127.0.0.1:8080/foo`

Request Body:

```
{"name": "{name}"}
```

Response Body:

```
{"name": "{name}"
 "id": "{id}"
```

### DELETE

`http://127.0.0.1:8080/foo/{id}`

Removes a record, if found, from the in-memory data
