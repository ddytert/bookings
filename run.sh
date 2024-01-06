go build -o bookings cmd/web/*.go && \
./bookings -dbname=bookings -dbuser=ddytert -dbpass=password -cache=false -production=false