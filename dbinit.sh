#! /bin/bash

echo "DB_DROP: ${DB_DROP}"

if [[ ${DB_DROP} = 'yes' ]]; then
  echo "dropping database ${DB_NAME}"
  dropdb -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} ${DB_NAME}

  echo "creating database ${DB_NAME}"
  createdb -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} ${DB_NAME}
fi
