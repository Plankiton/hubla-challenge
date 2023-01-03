if ! command -v docker &> /dev/null
then
    echo "MISSING DEPENDENCY: docker is not installed"
    exit 1
fi

if ! docker container inspect hubla_db > /dev/null 2>&1
then
    echo "- Starting a postgres on port 5432..."
    docker run --name hubla_db -d -p 5432:5432 \
        -e POSTGRES_USER=hubla \
        -e POSTGRES_PASSWORD=supersecret \
        -e POSTGRES_DB=hubla \
        postgres:13.3
    echo ""
fi

docker start hubla_db > /dev/null
docker exec -i hubla_db psql -U postgres -d hubla < ./pkg/db/sql/init.sql
