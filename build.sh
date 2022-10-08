#!/bin/bash

DIR="$( cd "$( dirname "$0"  )" && pwd  )"

bash $DIR/portal/build.sh
go build .
mkdir -p $DIR/dist
rm -rf $DIR/dist/*
cp -r $DIR/portal/build $DIR/dist/static
cp $DIR/etcd-dashboard $DIR/dist/etcd-dashboard
cd $DIR
