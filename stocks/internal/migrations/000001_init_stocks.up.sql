
CREATE TABLE IF NOT EXISTS sku (
    "sku_id" BIGINT PRIMARY KEY,
    "name" TEXT NOT NULL UNIQUE,
    "type" TEXT
);

ALTER TABLE "sku" OWNER TO "user_stocks";

INSERT INTO sku (sku_id, name, type) VALUES
(1001, 't-shirt', 'apparel'),
(2020, 'cup', 'accessory'),
(3033, 'book', 'stationery'),
(4044, 'pen', 'stationery'),
(5055, 'powerbank', 'electronics'),
(6066, 'hoody', 'apparel'),
(7077, 'umbrella', 'accessory'),
(8088, 'socks', 'apparel'),
(9099, 'wallet', 'accessory'),
(10101, 'pink-hoody', 'apparel')
ON CONFLICT (sku_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS items (
	"id" SERIAL PRIMARY KEY,
	"user_id" INT NOT NULL,
	"sku" BIGINT UNIQUE,
	"count" INT NOT NULL DEFAULT 0,
	"price" INT NOT NULL DEFAULT 0,
	"location" TEXT,
	"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT items_sku_fk
		FOREIGN KEY ("sku")
				REFERENCES sku("sku_id")
						ON DELETE SET NULL
);

ALTER TABLE "items" OWNER TO "user_stocks";
