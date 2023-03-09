CREATE TABLE geoloc.entry (
    ip inet PRIMARY KEY,
    code text,
    country text,
    city text,
    lat double,
    lon double
);
