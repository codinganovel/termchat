#!/bin/bash

echo "Building termchat..."
make build

echo "Starting test session..."
echo "Run the following command in another terminal:"
echo "./build/termchat join localhost:$(./build/termchat start | grep 'Session started:' | cut -d' ' -f3)"

./build/termchat start