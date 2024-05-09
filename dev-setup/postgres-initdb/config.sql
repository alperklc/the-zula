ALTER SYSTEM SET ssl_cert_file TO '/var/lib/postgresql/certs/cert.pem';
ALTER SYSTEM SET ssl_key_file TO '/var/lib/postgresql/certs/privkey.pem';
ALTER SYSTEM SET ssl TO 'ON';

CREATE DATABASE root;
CREATE ROLE root WITH LOGIN SUPERUSER PASSWORD 'zitadel';

create database zitadel;
