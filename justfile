set shell := ["bash", "-uc"]
set dotenv-load

up:
    docker compose up --build

pgcli:
    pgcli postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB

mycli:
    mycli mysql://$MYSQL_USER:$MYSQL_PASSWORD@$MYSQL_ROOT_HOST:$MYSQL_PORT/$MYSQL_DATABASE
