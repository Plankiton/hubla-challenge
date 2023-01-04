if ! command -v docker &> /dev/null
then
    echo "MISSING DEPENDENCY: docker is not installed"
    exit 1
fi

DB_USER="${DB_USER-hubla}"
DB_NAME="${DB_NAME-hubla}"
DB_PASS="${DB_PASS-supersecret}"

if ! docker container inspect hubla_db > /dev/null 2>&1
then
    echo "- Starting a postgres on port 5432..."
    docker run --name hubla_db -d -p 5432:5432 \
        -e POSTGRES_USER="$DB_USER" \
        -e POSTGRES_PASSWORD="$DB_PASS" \
        -e POSTGRES_DB="$DB_NAME" \
        postgres:13.3
    echo ""
fi

docker start hubla_db > /dev/null
docker exec -i hubla_db psql -U "$DB_USER" -d "$DB_NAME" < ./pkg/db/sql/init.sql

echo "- Waiting for postgres to be ready..."
RETRIES=10
until docker run -it --rm --link hubla_db:pg postgres:13.3 psql "postgres://$DB_USER:$DB_PASS@pg:5432/$DB_NAME" -c "select 1" > /dev/null || [ $RETRIES -eq 0 ]; do
    echo "    $((RETRIES--)) remaining attempts..."
    sleep 1;
done

if ! air -h > /dev/null 2>&1
then
  go build -o ./bin/server/ ./cmd/server && ./bin/server
else
  air -build.cmd "go build -o ./bin/server ./cmd/server" -build.bin "./bin/server"
fi
