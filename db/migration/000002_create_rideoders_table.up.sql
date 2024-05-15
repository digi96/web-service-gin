CREATE TABLE "rideorder" (
  "rideorder_id" SERIAL PRIMARY KEY,
  "contact_id" uuid NOT NULL,
  "rider_name" varchar NOT NULL,
  "rider_phone" varchar NOT NULL,
  "destination" varchar NOT NULL,
  "pickup_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp
);
ALTER TABLE "rideorder"
ADD FOREIGN KEY ("contact_id") REFERENCES "contacts" ("contact_id");