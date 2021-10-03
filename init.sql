begin;
CREATE TABLE if not exists records (
    sku varchar(50) NULL CONSTRAINT recordspk PRIMARY KEY
);
commit;
