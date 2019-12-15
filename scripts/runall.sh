#!/bin/bash
#kill $(pgrep ../build -f)

../build/auth_service -config ../config/auth.config.json -dbuser go -dbpass postgres &

#go run ../cmd/MailPicker/main.go -config=../config/mailpicker.config.json -dbuser go -dbpass postgres &

../build/mailpicker_service -config=../config/mailpicker.config.json -dbuser go -dbpass postgres &

../build/http_service -config=../config/linux_http_service.config.json -dbuser go -dbpass postgres &

#go run ../cmd/post/main.go -config=../config/post_service.config.json &
# chmod +x ./runall.sh
