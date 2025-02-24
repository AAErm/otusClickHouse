#!/bin/bash

if ! command -v jq &> /dev/null; then
    apk add --no-cache jq
fi

DATA_DIR="./home/data"
for file in "$DATA_DIR"/services_*.json; do
    if [[ "$file" == *_lines.json ]]; then
        continue
    fi

    jq -c '.[]' "$file" > "${file%.json}_lines.json"

    clickhouse-client --query="
        INSERT INTO services (id, Name, Category)
        SELECT
            rowNumberInAllBlocks() + (SELECT max(id) FROM services) AS id,
            Name,
            Category
        FROM input('Name String, Category String') FORMAT JSONEachRow" < "${file%.json}_lines.json"
done