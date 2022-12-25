import psycopg2 as pg
import os


PORT = 5432
user = os.getenv("username")
password = os.getenv("password")
host = os.getenv("host")
dbname = os.getenv("database")

conn = pg.connect(
    host=host,
    port=PORT,
    user=user,
    password=password,
    database=dbname,
    connect_timeout=10,
)
cur = conn.cursor()

cur.execute(open("schema.sql", "r").read())

conn.commit()
