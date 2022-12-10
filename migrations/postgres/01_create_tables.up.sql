CREATE TABLE IF NOT EXISTS "order" (
	"id" CHAR(36) PRIMARY KEY,
	"customer_name" VARCHAR(255) NOT NULL,
	"customer_address" VARCHAR(255) NOT NULL ,
	"customer_phone"   VARCHAR NOT NULL,
	"total_price" INT,
	"created_at" TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS "order_items" (
   	"id" CHAR(36) PRIMARY KEY,
   	"order_id" CHAR(36)  REFERENCES "order"("id")NOT NULL,
   	"product_id" CHAR(36) NOT NULL,
	"quantity" INT NOT NULL
	
);
