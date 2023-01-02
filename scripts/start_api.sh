echo "Source env";
source .env;

echo "Starting redis ..."
systemctl start redis;


echo "running api .."
cd api && go run .
