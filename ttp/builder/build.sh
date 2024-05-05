echo -en "yes\n" | git clone $1 repo --depth 1 --branch develop --single-branch --quiet
cd repo

# for debug
cd ta/example 

ego-go mod tidy
ego-go build

# for debug
#ego sign m 
ego sign example
