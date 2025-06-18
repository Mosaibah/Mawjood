# steps
- brew install cockroachdb/tap/cockroach
  992  brew install cockroachdb/tap/cockroach\n
  993  docker pull cockroachdb/cockroach:latest
  994  docker network create -d bridge roachnet\n
  995  docker volume create roach1\n
  996  docker run -d \\n--name=roach1 \\n--hostname=roach1 \\n--net=roachnet \\n-p 26257:26257 \\n-p 8080:8080 \\n-v "roach1:/cockroach/cockroach-data" \\ncockroachdb/cockroach:v25.2.1 start \\n  --advertise-addr=roach1:26357 \\n  --http-addr=roach1:8080 \\n  --listen-addr=roach1:26357 \\n  --sql-addr=roach1:26257 \\n  --insecure \\n  --join=roach1:26357,roach2:26357,roach3:26357
  

  docker exec -i roach1 cockroach sql --insecure < db_setup.sql