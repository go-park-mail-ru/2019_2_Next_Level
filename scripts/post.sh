#!/bin/bash

kill $(pgrep build/post -f)

./build/post_service -config ./config/post_service.config.json &