@echo off
set cloud=%1
if "%cloud%"=="" (
    set cloud=gcp
)
echo %cloud%

ssh %cloud% "sudo apt install -y git && git clone https://github.com/moririn2528/timetable.git && cd timetable && git checkout moririn"
scp api/sql/0_create_public.sql %cloud%:timetable/api/sql/0_create_public.sql
scp api/sql/0_init.sql %cloud%:timetable/api/sql/0_init.sql
scp api/sql/9_conceal.sql %cloud%:timetable/api/sql/9_conceal.sql
scp api/sql/my.cnf %cloud%:timetable/api/sql/my.cnf

scp docker/cloud/.env %cloud%:timetable/docker/cloud/.env
scp docker/cloud/public.env %cloud%:timetable/docker/cloud/public.env

scp cloud/startup.sh %cloud%:startup.sh

ssh %cloud% "sudo chmod 755 startup.sh && ./startup.sh"
ssh %cloud% "cd timetable/docker/cloud && docker compose up -d"

echo finished