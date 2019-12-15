#!/bin/bash
kill $(pgrep auth_service -f)
kill $(pgrep mailpicker_service -f)
kill $(pgrep http_serviceg -f)

../auth_service &


../mailpicker_service -config=../config/mailpicker.config.json -dbuser go -dbpass postgres &

../http_service -config=../config/http_service.config.json -dbuser go -dbpass postgres &

#go run ../cmd/post/main.go -config=../config/post_service.config.json &
# chmod +x ./runall.sh
