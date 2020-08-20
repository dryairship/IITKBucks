# IITKBucks

[![Project Status: Inactive – The project has reached a stable, usable state but is no longer being actively developed; support/maintenance will be provided as time allows.](https://www.repostatus.org/badges/latest/inactive.svg)](https://www.repostatus.org/#inactive) 
[![Docker Image Size](https://img.shields.io/docker/image-size/dryairship/iitkbucks/latest?color=022044&label=Docker%20Image)](https://hub.docker.com/r/dryairship/iitkbucks) 
[![Travis CI Build Status](https://img.shields.io/travis/dryairship/IITKBucks)](https://travis-ci.org/github/dryairship/IITKBucks) 
![Go version](https://img.shields.io/github/go-mod/go-version/dryairship/IITKBucks) 
[![Codecov](https://img.shields.io/codecov/c/github/dryairship/IITKBucks)](https://codecov.io/gh/dryairship/IITKBucks) 
[![LICENSE](https://img.shields.io/github/license/dryairship/IITKBucks?color=purple)](https://github.com/dryairship/IITKBucks/blob/master/LICENSE)


IITKBucks is a cryprocurrency project that I mentored as a part of Programming Club's Summer Camp 2020. This repository contains the code for my client, implemented in Go.

## The Project

A short documentation of the project is available at [https://iitkbucks.pclub.in/](https://iitkbucks.pclub.in/). The messages and notes shared with the students who did the project are available at [dryairship/IITKBucks-meta](https://github.com/dryairship/IITKBucks-meta/).

## Running

### Pull docker image and run
Pull the docker image:
```
docker pull dryairship/iitkbucks
```

Get a public key, update the [config file](https://github.com/dryairship/IITKBucks/blob/master/iitkbucks-config.yml), and then use [docker-compose](https://github.com/dryairship/IITKBucks/blob/master/docker-compose.yml) to start the container:
```
docker-compose up -d iitkbucks
```

### Build from source and run
Clone the repo:
```
git clone git@github.com:dryairship/IITKBucks.git
cd IITKBucks
git submodule update --init
```

Build the frontend: (This will get you a `build` folder inside the `frontend` directory).
```
cd frontend
npm i
npm run build
cd ..
```

Build the backend: (This will get you a binary file named `iitkbucks` inside the project directory).
```
go build ./cmd/iitkbucks
```

Get a public key, update the [config file](https://github.com/dryairship/IITKBucks/blob/master/iitkbucks-config.yml), and then start the server:
```
./iitkbucks
```
