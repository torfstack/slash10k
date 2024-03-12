-- Create "bot_setup" table
CREATE TABLE "bot_setup" (
  "channel_id" text NOT NULL,
  "message_id" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now()
);
-- Create "player" table
CREATE TABLE "player" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "player_name_key" to table: "player"
CREATE UNIQUE INDEX "player_name_key" ON "player" ("name");
-- Create "char" table
CREATE TABLE "char" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  "class" character varying(255) NOT NULL,
  "user_id" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "char_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "player" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "debt" table
CREATE TABLE "debt" (
  "id" serial NOT NULL,
  "amount" bigint NOT NULL,
  "last_updated" timestamp NOT NULL DEFAULT now(),
  "user_id" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "debt_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "player" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "debt_user_id_key" to table: "debt"
CREATE UNIQUE INDEX "debt_user_id_key" ON "debt" ("user_id");
-- Create "debt_journal" table
CREATE TABLE "debt_journal" (
  "id" serial NOT NULL,
  "amount" bigint NOT NULL,
  "description" text NOT NULL DEFAULT '',
  "date" date NOT NULL DEFAULT now(),
  "user_id" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "debt_journal_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "player" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
