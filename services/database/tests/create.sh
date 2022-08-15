echo "CREATING DB:\n"
psql -U user -d db -f /scripts/create_db.sql
echo "ADDING FAKE DATA:\n"
psql -U user -d db -f /scripts/test.sql