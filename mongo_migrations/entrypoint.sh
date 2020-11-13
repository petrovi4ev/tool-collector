#!/bin/sh
set -e
set -x
mongo_init_migrations_table=$(cat <<HEREDOC
const previous_tables_migration_event = { 
    "_id" : ObjectId("5fae885b61f6b79f18304e4a"), 
    "fileName" : "20201029100341-add_existing_tables.js", 
    "appliedAt" : ISODate("2020-11-13T13:21:31.998Z") 
}
db.getCollection('changelog_receive_tx_service').insertOne(previous_tables_migration_event)
HEREDOC
)
echo MIGRATE_MONGO
sed -i "s%MONGO_URI_PLACEHOLDER%$1%" ./migrate-mongo-config.js
cat ./migrate-mongo-config.js
migrate-mongo status
mongo "$1" --quiet --eval "$mongo_init_migrations_table"
migrate-mongo status
migrate-mongo up
migrate-mongo status
