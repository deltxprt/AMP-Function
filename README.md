# AMP-Function

# What is this project?

This project is partly to practice writing GO applications and web application in general.
It's also an app to manage AMP from Discord via API calls.

# How to use

You will need to create a config file in the same directory as the executable. The config file is a YAML file.

# Dependencies

- [Redis](https://redis.io/)

## Container Example

### Custom DockerFile/Container File

see [ContainerFile](/Containerfile)

### adhoc

`docker run -d -v /path/to/config.yaml:/etc/api/config/config.yaml --name amp-function --restart unless-stopped amp-function:latest`

### docker-compose

see [docker-compose.yml](/examples/docker-compose.yaml)

### Podman (WIP)

see [podman-pods.yml](/podman-pods.yaml)

### Kubernetes

see [kubernetes example folder](/examples/kuberntes)

## Config Examples

*config.yaml*
```yaml
port: 8081
env: "Production"
RefreshInterval: "15s"
postgres:
  dsn: "postgres://db_user:4n0Th3rP4$$W0rd@172.16.1.1/mydb?sslmode=disable"
  maxOpenConns: 5
  maxIdleConns: 5
  maxIdleTime: "15m"
amp:
  url: "https://my.amp.enpoint.com"
  username: "ReadUser"
  password: "5uper$eCretP4$$W0rd"
```
environment variables

```editorconfig
PORT=8080

ENVIRONNEMENT=Production

REFRESH_INTERVAL=15s

AMP_URL=https://my.amp.enpoint.com

AMP_USERNAME=ReadUser

AMP_PASSWORD=5uper$eCretP4$$W0rd

#your otp token if the account have one
AMP_TOKEN=123456

AMP_REMEMBER_ME=false

#if you use ssl mode switch to enable
POSTGRES_DSN=postgres://db_user:4n0Th3rP4$$W0rd@172.16.1.1/mydb?sslmode=disable 

POSTGRES_MAX_OPEN_CONNS=5

POSTGRES_MAX_IDLE_CONNS=5

POSTGRES_MAX_IDLE_TIME=15m
```

## AMP user Required Permissions

- AMPCore
  - App Management
    - Read Console
- All Instances
  - The Instance Name
    - Manage

* If you want to make the use not see an instance just don't give them the Manage permission for that instance. (see example below)

here's an example of the permissions i gave to the User:
![AMP_Permissions_Menu](/docs/AMP_Permissions_Menu.png)
