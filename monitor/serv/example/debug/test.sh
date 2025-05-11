mkdir -p ./tmp
cp -r ./api ./tmp
cp -r ../static ./tmp
cp -r ../views/verification-status ./tmp
cd ./tmp
python3 -m http.server 8080
