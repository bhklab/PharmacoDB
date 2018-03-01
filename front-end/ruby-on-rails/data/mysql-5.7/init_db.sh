#!/bin/bash

#/usr/bin/mysqld_safe &
#sleep 5

mysql -u root -e "CREATE DATABASE pharmacodb_development;"
mysql -u root pharmacodb_development < /docker-entrypoint-initdb.d/dump.sql
