fmt:
    @goimports -w -local="github.com/earlgray283/fixgen" .
    @gofumpt -w .
    @dprint fmt

test:
    @go test -v ./...

test-golden:
    @go test ./test -run="^Test_GoldenTest"

update-golden:
    @go test ./test -run="^Test_GoldenTest" -update

install:
    @go install .

update-example: install
    @cd .examples/yo && rm -rf fixture/* && fixgen yo
    @cd .examples/ent && rm -rf fixture/* && fixgen ent

gen-yo: run-spanner-image
    @cd ./test/yo/test && go tool yo test-project test-instance test-db -o ./models

spanner_image_name := "fixgen-spanner"
spanner_container_name := "fixgen-spanner"
spanner_schema_path := "./test/yo/test/schema.sql"
spanner_project_id := "test-project"
spanner_instance_id := "test-instance"
spanner_database_id := "test-db"

build-spanner-image:
    @docker build \
      --build-arg PROJECT_ID={{ spanner_project_id }} \
      -t {{ spanner_image_name }}:latest \
      -f ./dockerfiles/Dockerfile_spanner .

run-spanner-image: build-spanner-image
    @docker stop {{ spanner_container_name }} 2>/dev/null || :
    @docker rm {{ spanner_container_name }} 2>/dev/null || :
    @docker run -d -p 9010:9010 -p 9020:9020 \
      -v {{ spanner_schema_path }}:/schema.sql \
      --name {{ spanner_container_name }} \
      {{ spanner_image_name }}:latest
    @docker exec {{ spanner_container_name }} \
      gcloud spanner instances create {{ spanner_instance_id }} \
        --config=emulator-config --description="Test Instance" --nodes=1
    @docker exec {{ spanner_container_name }} \
      gcloud spanner databases create {{ spanner_database_id }} \
        --instance={{ spanner_instance_id }} --ddl-file=schema.sql
