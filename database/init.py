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

cur.execute("""
CREATE TABLE IF NOT EXISTS "items" (
  "id" SERIAL,
  "name" varchar(255) NOT NULL,
  "brand" varchar(255) NOT NULL,
  "description" varchar(255) NOT NULL,
  "price" decimal(6,2) NOT NULL,
  "image" varchar(255) NOT NULL,
  PRIMARY KEY ("id")
);

-- added IDs to prevent re-inserting
INSERT INTO "items" ("id", "name", "brand", "description", "price", "image") VALUES
(1,'Teddy Bear', 'LotFancy', 'A cute teddy bear', 14.99, 'teddy-bear.png'),
(2,'Doll', 'Barbie', 'A sweet doll', 9.48, 'doll.png'),
(3,'Toy Car 1', 'Hot Wheels', 'A cute datsun toy car', 57.14, 'toy-car-1.png'),
(4,'Toy Car 2', 'Hot Wheels', 'A cute ford toy truck', 39.99, 'toy-car-2.png'),
(5,'Toy Car 3', 'Hot Wheels', 'A cute mercedes toy racecar', 54.74, 'toy-car-3.png'),
(6,'Lego Classics', 'Lego', 'A cute lego monsters set', 9.97, 'lego-monsters.png'),
(7,'Lego Dots', 'Lego', 'A cute lego dots set', 17.94, 'lego-dots.png'),
(8,'Sweet Jumperoo', 'Fisher-Price', 'A sweet ride jumperoo', 124.99, 'ride-jumperoo.png'),
(9,'Musical Keyboard', 'CoComelon', 'A sweet musical keyboard', 26.99, 'musical-keyboard.png'),
(10,'T-Shirt & Shorts Set 1', 'CoComelon', 'A sweet t-shirt & shorts set', 18.99, 'tshirt-shorts-1.png'),
(11,'T-Shirt & Shorts Set 2', 'CoComelon', 'A sweet t-shirt & shorts set', 18.99, 'tshirt-shorts-2.png'),
(12,'T-Shirt & Shorts Set 3', 'CoComelon', 'A sweet t-shirt & shorts set', 18.99, 'tshirt-shorts-3.png'),
(13,'N-Strike Blaster', 'Nerf', 'A powerful blaster gun', 34.99, 'strike-blaster.png'),
(14,'Baby Mickey Mouse', 'Disney', 'A sweet baby Mickey plush', 18.88, 'baby-mickey.png'),
(15,'Baby Minnie Mouse', 'Disney', 'A sweet baby Minnie plush', 51.23, 'baby-minnie.png'),
(16,'3D Toddler Scooter', 'Bluey', 'A fantastic 3-wheel scooter', 29.00, 'toddler-scooter.png'),
(17,'Cottage Playhouse', 'Litte Tikes', 'A fancy blue playhouse', 139.99, 'cottage-playhouse.png'),
(18,'2-in-1 Motor/Wood Shop', 'Little Tikes', 'A realistic motor/wood shop', 99.00, '2x1-motor-shop.png'),
(19,'UNO Collector Tin', 'UNO', 'A premium quality uno set', 49.39, 'uno-phase10-snappy.png'),
(20,'Razor MX350 Bike', 'Razor', 'A powerful electric bike', 328.00, 'mx350-bike.png')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "email" TEXT UNIQUE,
    "salt" TEXT,
    "hashed_password" BYTEA
);

CREATE TABLE IF NOT EXISTS "purchases" (
    "user_id" INT NOT NULL, 
    "item_id" INT NOT NULL,
    "quantity" INT NOT NULL DEFAULT 1,
    "purchase_date" timestamp DEFAULT NOW(),
    PRIMARY KEY ("user_id", "item_id", "purchase_date"),
    FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
    FOREIGN KEY ("item_id") REFERENCES "items"("id") ON DELETE CASCADE
);
""")

conn.commit()
