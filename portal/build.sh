#!/bin/bash
DIR="$( cd "$( dirname "$0"  )" && pwd  )"
cd $DIR

npm install
npm run build
