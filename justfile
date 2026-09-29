default: run
run:
    echo "Running Cyber Navigate"

build service_name:
    docker build -f ./docker/{{ service_name }}.dockerfile -t cyber-navigate/{{ service_name }} .
