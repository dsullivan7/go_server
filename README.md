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
make db-run
make db-create
make db-migrate
```

### Run the web server

```
make run
```

### Test

```
make test
```

### Build and deploy

```
make build
make deploy
```

## Operations

###Create user

```
POST
/api/users/
{ "FirstName": "MyFirstName" }
```

###Get user
```
GET
/api/users/:userId
```

###List users
```
GET
/api/users
```

###Modify user
```
PUT
/api/users/:userId
{ "FirstName": "DifferentFirstName" }
```

###Delete user
```
DELETE
/api/users/:userId
```
