#!/bin/bash

cd /app
go build -o /app/docker-gs-ping-roach
/app/docker-gs-ping-roach migrate
/app/docker-gs-ping-roach
