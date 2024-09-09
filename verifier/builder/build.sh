### Clone
echo -en "yes\n" | git clone $1 $2 --depth 1 --branch $3 --single-branch --quiet >&2
cd $2/$4 >&2

### Show Git Log
git log -1 --pretty=format:%H
echo ""

### Build
ego-go build -buildvcs=false -trimpath=true >&2
ego sign m >&2

### Show UniqueID
ego uniqueid $5
