-- +migrate Up
CREATE TABLE "game" (
    "id" INTEGER CONSTRAINT "game_pk" PRIMARY KEY,
    "psnID" TEXT NOT NULL UNIQUE,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "iconURL" TEXT NOT NULL,
    "platform" TEXT NOT NULL,
    "psnServiceName" TEXT NOT NULL
);

-- +migrate Down
DROP TABLE "game";