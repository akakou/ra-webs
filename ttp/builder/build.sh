echo -en "yes\n" | git clone $1 repo --depth 1 --branch debug-pub-5-15 --single-branch --quiet
#echo -en "yes\n" | git clone $1 repo --depth 1 --branch main --single-branch --quiet

cd repo/ta/example
#cd repo

ego-go build -buildvcs=false
ego sign example
#ego sign m 

