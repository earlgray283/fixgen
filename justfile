spannerContainerName := "fixgen-spanner"

test:
    @go test -v ./...

golden-test: golden-test-ent golden-test-yo

golden-test-update:
    @go test ./test/... -run="^Test_GoldenTest" -update

golden-test-ent:
    @go test -v ./test/ent/...

golden-test-yo:
    @go test -v ./test/yo/...

build-spanner-image:
    @docker build -t fixgen-spanner:latest -f ./dockerfiles/Dockerfile_spanner . 

run-spanner-image: build-spanner-image
    @docker stop {{ spannerContainerName }} 2>/dev/null || :
    @docker rm {{ spannerContainerName }} 2>/dev/null || :
    @docker run -d -p 9010:9010 -p 9020:9020 --name {{ spannerContainerName }} fixgen-spanner:latest
    @docker exec {{ spannerContainerName }} sh -c \
      'gcloud spanner instances create test-instance \
        --config=emulator-config --description="Test Instance" --nodes=1 \
      && gcloud spanner databases create test-db \
        --instance=test-instance --ddl-file=schema.sql'
