echo -en "yes\n" | git clone $1 repo --depth 1 --branch main --single-branch --quiet
cd repo

ego-go mod tidy
ego-go build
ego sign m 

