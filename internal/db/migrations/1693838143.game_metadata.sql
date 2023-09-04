-- +migrate Up
ALTER TABLE "game" ADD COLUMN "backgroundImageURL" TEXT;
ALTER TABLE "game" ADD COLUMN "metacriticScore" INTEGER;
ALTER TABLE "game" ADD COLUMN "releaseDate" TEXT;

-- +migrate Down
ALTER TABLE "game" DROP COLUMN "releaseDate";
ALTER TABLE "game" DROP COLUMN "metacriticScore";
ALTER TABLE "game" DROP COLUMN "backgroundImageURL";
