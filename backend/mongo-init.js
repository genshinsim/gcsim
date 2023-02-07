db.auth("root", "root-password");

db = db.getSiblingDB("gcsim-database");

db.createCollection("data");

db.createView(
    "gcsimvaliddb",
    "data",
    [ { $match: { is_db_valid: true } }]
)