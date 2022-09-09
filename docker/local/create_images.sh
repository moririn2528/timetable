#!/bin/bash
cd ../../
docker image build -t moririn2528/timetable-app:latest .
docker push moririn2528/timetable-app:latest
cd api/sql
docker image build -t moririn2528/timetable-db:latest .
docker push moririn2528/timetable-db:latest
# want: docker login