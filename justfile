spannerContainerName := "fixgen-spanner"

fmt:
    @goimports -w -local="github.com/earlgray283/fixgen" .
    @dprint fmt

test:
    @go test -v ./...

test-golden-all:
    @go test ./test/... -run="^Test_GoldenTest"

test-golden type:
    @go test ./test/{{ type }} -run="^Test_GoldenTest"

update-golden-all:
    @go test ./test/... -run="^Test_GoldenTest" -update

update-golden type:
    @go test ./test/{{ type }} -run="^Test_GoldenTest" -update

install:
    @go install .

update-example: install
    @cd .examples/yo && rm -rf fixture/* && fixgen yo
    @cd .examples/ent && rm -rf fixture/* && fixgen ent

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
