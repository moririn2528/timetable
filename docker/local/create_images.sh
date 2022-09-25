#!/bin/bash
echo a | docker login
code=$?
if [ $code -ne 0 ]; then
    echo "Error: please login docker"
    exit
fi
cd ../../
docker image build -t moririn2528/timetable-app:latest .
docker push moririn2528/timetable-app:latest
cd api/sql
docker image build -t moririn2528/timetable-db:latest .
docker push moririn2528/timetable-db:latest
# want: docker login