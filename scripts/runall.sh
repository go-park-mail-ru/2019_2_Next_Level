#!/bin/bash

../build/auth_service &

#go run ../cmd/MailPicker/main.go -config=../config/mailpicker.config.json -dbuser go -dbpass postgres &

../build/http_service -config=../config/http_service.config.json -dbuser go -dbpass postgres &

#go run ../cmd/post/main.go -config=../config/post_service.config.json &
# chmod +x ./runall.sh
