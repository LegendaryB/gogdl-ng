<p align="center"><img alt="gogdl-ng" height="60" src="https://raw.githubusercontent.com/LegendaryB/gogdl-ng/develop/.github/assets/banner.png"></p>

<div align="center">

[![forthebadge](https://forthebadge.com/images/badges/fuck-it-ship-it.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

[![GitHub license](https://img.shields.io/github/license/gogdl-ng/gogdl-ng.svg?longCache=true&style=flat-square)](https://github.com/gogdl-ng/gogdl-ng/blob/main/LICENSE)
[![GitHub license](https://img.shields.io/docker/v/legendaryb/gogdl-ng/latest?longCache=true&style=flat-square)](https://hub.docker.com/r/legendaryb/gogdl-ng)

A self-hostable application to download files in a folder from Google Drive powered by Go.
</div><br>

## 🎯 Features
* Google Drive v3 API based
* With additional Firefox and Chrome extension
* OAuth 2.0 authorization
* Integrity checks (MD5 checksum)
* Transfer retries
* Hassle-free setup thanks to Docker ❤︎

## 📝 Requirements
- A Google Cloud Platform project with the Drive API enabled. [Guide](https://developers.google.com/drive/api/v3/quickstart/go#step_1_turn_on_the)
- Basic docker and docker-compose knowledge.
- The gogdl-ng browser [extension](https://github.com/gogdl-ng/gogdl-ng-webext).

## Installation
1. Create a `config` folder and create a `config.toml` with the following content:
```toml
title = "gogdl-ng"

[application]
# Defines the port on which the application is listening for requests.
listenPort = 3200

[logging]
# Defines the log level. Must be one of the following: trace, debug, info, warn, error, fatal, panic
logLevel = "info"

# Defines the location where to write the application log file.
logFilePath = "./config/gogdl-ng.log"

# Defines if the messages should also be logged to console or not
logToConsole = true

[queue]
# Defines the maximum capacity of the job queue.
size = 1000

# Defines how many workers can run concurrently.
maxWorkers = 2

[gdrive]
# query string which will be appended to the base query which is: "'%drive_folder_id%' in parents and"
query = "trashed=false"

[download]
# Defines how many times a failed download should be retried.
retryThreeshold = 5
```
2. Copy the `*.json ` file which you got at the end of the Google Cloud Platform project guide into the `config` folder. Rename it to `credentials.json`.
3. Start the container once with `docker run`. Like this:  
`docker run -i -p 3200:3200 -v /path/to/config:/config -v /path/to/downloads:/downloads legendaryb/gogdl-ng:latest`  
Follow the instructions as shown in the terminal. You need to enter the authorization code. After that you should exit via pressing CTRL+C
4. Create the docker-compose.yml file (adjust it as you need): 
```
version: '3'

services:
  gogdl-ng:
    image: legendaryb/gogdl-ng:latest
    container_name: gogdl-ng
    volumes:
      - ./config:/config
      - ./downloads:/downloads
    ports:
      - 3200:3200
    restart: always
```
5. Now you can bring the service up: `docker-compose up`
