drop KEYSPACE if exists demo

CREATE KEYSPACE  IF NOT EXISTS demo
  WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
  
  
CREATE TABLE IF NOT EXISTS demo.messages (
  address int,
  message text,
  value int,
  created text;
  PRIMARY KEY(created, address)  
);