use mysql;
create user 'replicator'@'%' identified by 'replicator';
grant replication slave on *.* to 'replicator'@'%';
FLUSH PRIVILEGES;
SHOW MASTER STATUS;
SHOW VARIABLES LIKE 'server_id';
