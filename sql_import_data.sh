#!/bin/bash -x
mysql --user="$DB_USER" --password="$DB_PWD" -h $DB_HOST $DB_NAME -e "source $1"
