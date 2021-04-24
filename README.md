#Go Server
---
A web server built in golang

##Get Started
---

###Docker
---
Install Docker
https://docs.docker.com/get-docker/

###Environment
---
Environment needs to be specified in an `.env` file

###Initialize Database
---
```
make db-run
make db-create
make db-migrate
```

###Run the web server
---
```
make run
```

###Test
---
```
make test
```

###Build and deploy
---
```
make build
make deploy
```
