#!/bin/bash

# Run the create_table.sql script
sqlite3 permutations.db < create_table.sql

# Remove the create_table.sql file
rm create_table.sql

# Run each numbered SQL file, vacuum the database, and remove the file
for sql_file in sql_statements_*.sql; do
  sqlite3 permutations.db < "$sql_file"
  sqlite3 permutations.db "VACUUM;"
  rm "$sql_file"
done

# Vacuum the database one more time
sqlite3 permutations.db "VACUUM;"