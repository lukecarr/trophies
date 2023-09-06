-- +migrate Up
CREATE TABLE "gameMetadata" (
    "gameID" INTEGER NOT NULL UNIQUE,
    "backgroundImageURL" TEXT,
    "metacriticScore" INTEGER,
    "releaseDate" TEXT
);
INSERT INTO "gameMetadata" ("gameID", "backgroundImageURL", "metacriticScore", "releaseDate")
SELECT "id", "backgroundImageURL", "metacriticScore", "releaseDate" FROM "game" WHERE "backgroundImageURL" IS NOT NULL OR "metacriticScore" IS NOT NULL OR "releaseDate" IS NOT NULL;
ALTER TABLE "game" DROP COLUMN "backgroundImageURL";
ALTER TABLE "game" DROP COLUMN "metacriticScore";
ALTER TABLE "game" DROP COLUMN "releaseDate";

-- +migrate Down
ALTER TABLE "game" ADD COLUMN "releaseDate" TEXT;
ALTER TABLE "game" ADD COLUMN "metacriticScore" INTEGER;
ALTER TABLE "game" ADD COLUMN "backgroundImageURL" TEXT;
UPDATE "game" SET "backgroundImageURL" = "gameMetadata"."backgroundImageURL", "metacriticScore" = "gameMetadata"."metacriticScore", "releaseDate" = "gameMetadata"."releaseDate"
FROM (SELECT "gameID", "backgroundImageURL", "metacriticScore", "releaseDate" FROM "gameMetadata") AS "gameMetadata"
WHERE "game"."id" = "gameMetadata"."gameID";
DROP TABLE "gameMetadata";
