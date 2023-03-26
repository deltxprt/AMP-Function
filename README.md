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

## Config File Example

*config.yaml*
```yaml
port: 8081
env: "Production"
database:
  address: "localhost:6379"
  database: 0
amp:
  url: "https://my.amp.enpoint.com"
  username: "ReadUser"
  password: "5uper$eCretP4$$W0rd"
```

## AMP user Required Permissions

- AMPCore
  - App Management
    - Read Console
- All Instances
  - The Instance Name
    - Manage

* If you want to make the use not see an instance just don't give them the Manage permission for that instance. (see example below)

here an example of the permissions i gave to the User:
![AMP_Permissions_Menu](/docs/AMP_Permissions_Menu.png)
