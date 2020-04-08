#!/usr/bin/env bash

# server should be up on 127.0.0.1:8080 before running this script
# non zero return value indicates failure

curl --silent --show-error --fail \
     127.0.0.1:8080/api/fetcher -X POST \
     -d '{"url":"https://httpbin.org/range/15","interval":60}' \
     > /dev/null

if [ $? -ne 0 ]; then
    exit 1
fi


curl --silent --show-error --fail \
     127.0.0.1:8080/api/fetcher -X POST \
     -d '{"url":"https://httpbin.org/delay/10","interval":120}' \
     > /dev/null

if [ $? -ne 0 ]; then
    exit 2
fi


curl --silent --fail \
    127.0.0.1:8080/api/fetcher -X POST \
    -d 'niepoprawny json'

if [ $? -eq 0 ]; then
    exit 3
fi

# make large request
cat /dev/urandom | \
head -n 1000000 | \
curl -si \
    127.0.0.1:8080/api/fetcher -X POST \
    -H "Content-Type: application/json" \
    -d @- | \
grep "413 Request Entity Too Large" -q

if [ $? -ne 0 ]; then
    exit 4
fi


curl -s 127.0.0.1:8080/api/fetcher | \
grep '"url":"https://httpbin.org/delay/10","interval":120' -q

if [ $? -ne 0 ]; then
    exit 5
fi


curl -s 127.0.0.1:8080/api/fetcher | \
grep '"url":"https://httpbin.org/range/15","interval":60' -q

if [ $? -ne 0 ]; then
    exit 6
fi


# change interval
curl --silent --show-error --fail \
     127.0.0.1:8080/api/fetcher -X POST \
     -d '{"url":"https://httpbin.org/delay/10","interval":240}' \
     > /dev/null

if [ $? -ne 0 ]; then
    exit 7
fi


curl -s 127.0.0.1:8080/api/fetcher | \
grep '"url":"https://httpbin.org/delay/10","interval":240' -q

if [ $? -ne 0 ]; then
    exit 8
fi


# deletion
curl --silent --fail \
    127.0.0.1:8080/api/fetcher/12345678 -X DELETE

if [ $? -eq 0 ]; then
    exit 9
fi

curl --silent --fail \
    127.0.0.1:8080/api/fetcher/1 -X DELETE

curl --silent --fail 127.0.0.1:8080/api/fetcher/1/history

if [ $? -eq 0 ]; then
    exit 11
fi
