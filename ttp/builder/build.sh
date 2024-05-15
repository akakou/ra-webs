echo -en "yes\n" | git clone $1 repo --depth 1 --branch main --single-branch --quiet
cd repo

ego-go build -buildvcs=false -trimpath=true
ego sign m 

