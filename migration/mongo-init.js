db.auth("root", "root-password");

db = db.getSiblingDB("store");

db.createCollection("data");

db.data.setIndex("create_time");