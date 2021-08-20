# Go Server

A web server built in golang

## Set Up

### Docker

Install Docker
https://docs.docker.com/get-docker/

### Environment

Environment needs to be specified in an `.env` file

### Initialize Database

```
make db-start
make db-create
make db-migrate
```

### Run the web server

```
make run-docker
```

### Test

```
make test-docker
```

### Build

```
make build-docker
```

## Operations

### Create user

```
POST
/api/users/
{ "FirstName": "MyFirstName" }
```

### Get user

```
GET
/api/users/:userId
```

### List users

```
GET
/api/users
```

### Modify user

```
PUT
/api/users/:userId
{ "FirstName": "DifferentFirstName" }
```

### Delete user

```
DELETE
/api/users/:userId
```
