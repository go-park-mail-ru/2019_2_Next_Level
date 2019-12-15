#!/bin/bash

kill $(pgrep ../build -f)

./build/post_service -config ./config/post_service.config.json &