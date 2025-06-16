CREATE TABLE IF NOT EXISTS "User" (
	"id" serial NOT NULL UNIQUE,
	"name" varchar(255) NOT NULL,
	"email" varchar(255) NOT NULL,
	"password" varchar(255) NOT NULL,
	"age" bigint NOT NULL,
	"gender" varchar(255) NOT NULL,
	"weight" double precision NOT NULL,
	"height" double precision NOT NULL,
	"goal" varchar(255) NOT NULL,
	"activity_level" varchar(255) NOT NULL,
	"created_at" timestamp with time zone NOT NULL,
	"role" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "food" (
	"id" serial NOT NULL UNIQUE,
	"name" varchar(255) NOT NULL,
	"serving_size_gram" double precision NOT NULL,
	"source" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "nutrient" (
	"id" serial NOT NULL UNIQUE,
	"name" varchar(255) NOT NULL,
	"category" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "food_nutrients" (
	"id" serial NOT NULL UNIQUE,
	"food_id" bigint NOT NULL,
	"nutrient_id" bigint NOT NULL,
	"amount_per_100g" double precision NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "meal_log" (
	"id" serial NOT NULL UNIQUE,
	"user_id" bigint NOT NULL,
	"created_at" date NOT NULL,
	"meal_type" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "meal_log_items" (
	"id" serial NOT NULL UNIQUE,
	"meal_log_id" bigint NOT NULL,
	"food_id" bigint NOT NULL,
	"quantity" bigint NOT NULL,
	"quantity_grams" double precision NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "user_biometrics" (
	"id" serial NOT NULL UNIQUE,
	"user_id" bigint NOT NULL,
	"created_at" date NOT NULL,
	"type" varchar(255) NOT NULL,
	"value" double precision NOT NULL,
	"unit" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);




ALTER TABLE "food_nutrients" ADD CONSTRAINT "food_nutrients_fk1" FOREIGN KEY ("food_id") REFERENCES "food"("id");

ALTER TABLE "food_nutrients" ADD CONSTRAINT "food_nutrients_fk2" FOREIGN KEY ("nutrient_id") REFERENCES "nutrient"("id");
ALTER TABLE "meal_log" ADD CONSTRAINT "meal_log_fk1" FOREIGN KEY ("user_id") REFERENCES "User"("id");
ALTER TABLE "meal_log_items" ADD CONSTRAINT "meal_log_items_fk1" FOREIGN KEY ("meal_log_id") REFERENCES "meal_log"("id");

ALTER TABLE "meal_log_items" ADD CONSTRAINT "meal_log_items_fk2" FOREIGN KEY ("food_id") REFERENCES "food"("id");
ALTER TABLE "user_biometrics" ADD CONSTRAINT "user_biometrics_fk1" FOREIGN KEY ("user_id") REFERENCES "User"("id");