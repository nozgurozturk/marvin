<p align="center">
  <a href="https://github.com/nozgurozturk/noo-analytics">
    <img style="-webkit-user-select: none; display: block; margin: auto; padding: env(safe-area-inset-top) env(safe-area-inset-right) env(safe-area-inset-bottom) env(safe-area-inset-left); cursor: zoom-in;" src="https://i.ibb.co/fCfFyHp/marvin-logo.png" width=" 190" height="210">
  </a>

  <h3 align="center">Marvin</h3>

  <p align="center">
    marvin is a paranoid !android that is control repositories and reports deprecated dependencies
  </p>
</p>

## Install
First of all, [download](https://golang.org/dl/) and install Go. `1.14` or higher is required for api and notification service.

After that, [download](https://nodejs.org/en/download/) and instal node. `12.14.0` or higher is required for cli.

When installation is done clone the repo with command:

```bash
git clone https://github.com/nozgurozturk/marvin
 ```
After that install missing dependencies for api and notifier with command:

```bash
cd ./server
go mod tidy
```

```bash
cd ./notifier
go mod tidy
```

## Run the api

    go run ./server/cmd

## Run the notifier

    go run ./notifier/cmd

## Run the cli

    cd ./cli
    yarn marvin <command>

## Docker

    docker-compose up --build

## Database Support

marvin works with **MongoDB** and **Redis** . You need to install dbs for local development.

## Environment Variables

### For API and Notifier

You can find sample .env files in each folder

**MongoDB Variables:**
```.env
MONGO_DB_USERNAME = admin
MONGO_DB_PASSWORD = password
MONGO_DB_HOST = 127.0.0.1
MONGO_DB_PORT = 27017
MONGO_DB_NAME = marvin
MONGO_DB_QUERY = 
```

**Server Variables:**
```.env
HOST = localhost
PORT = 8081
ACCESS_SECRET = superaccesssecret
REFRESH_SECRET = superrefreshsecret
SUB_SECRET = supersubsecret
# Minute
ACCESS_EXPIRE = 60
# Hour
REFRESH_EXPIRE = 720
# Hour
SUB_EXPIRE = 24
```

**Email Variables**
```.env
EMAIL_PORT = :587
EMAIL_HOST = mail.hostservice.com
EMAIL_FROM = example@mail.com
EMAIL_PASSWORD = dummypass
```

### For API Only

**Redis Variables:**
```.env
REDIS_DB_ADDRESS = 127.0.0.1:6379
REDIS_DB_USERNAME =
REDIS_DB_PASSWORD =
REDIS_DB = 0
```

### For Notifier Only

MAIN_HOST variable must be same as HOST in **server**
MAIN_PORT variable must be same as PORT in **server**

**Server Variables**
```.env
MAIN_HOST = localhost
MAIN_PORT = 8081
```

### For CLI Only

**Server Variables**
```.env
LOCAL_DIR = /.marvin
API_HOST = http://localhost:8081
```

## API Documentation

You can find api documentation on [http://localhost:8081/docs/index.html](http://localhost:8081/dics/index.html).

## CLI Documentation

You can find CLI documentation whit run command in `./cli`

    yarn marvin -h

    Usage: index [options] [command]

    Options:
    -V, --version   output the version number
    -h, --help      display help for command

    Commands:
    login           Authenticates users
    signup          Creates user
    use             Sets default user
    create <url>    Creates new repository
    update          Updates repository packages
    list-repo       List all repositories of user
    delete          Deletes repository
    sub-add         Add subscriber to repository
    sub-list        List all subscribers belongs to repository, reds are unverfied
    
    send            Sends mail to subscriber
    help [command]  display help for command
