db.auth("root", "root-password");

db = db.getSiblingDB("gcsim-database");

db.createCollection("data");
