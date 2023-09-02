-- +migrate Up
CREATE TABLE "trophyGroup" (
    "id" INTEGER CONSTRAINT "trophyGroup_pk" PRIMARY KEY,
    "psnID" TEXT NOT NULL,
    "gameID" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "iconURL" TEXT NOT NULL
);

CREATE UNIQUE INDEX "trophyGroup_psnID_gameID" ON "trophyGroup" ("psnID", "gameID");

CREATE TABLE "trophy" (
    "id" INTEGER CONSTRAINT "trophy_pk" PRIMARY KEY,
    "psnID" INTEGER NOT NULL,
    "gameID" INTEGER NOT NULL,
    "trophyGroupID" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "rarity" TEXT NOT NULL,
    "iconURL" TEXT NOT NULL,
    "hidden" BOOLEAN NOT NULL
);

CREATE UNIQUE INDEX "trophy_psnID_gameID" ON "trophy" ("psnID", "gameID");

-- +migrate Down
DROP TABLE "trophy";
DROP TABLE "trophyGroup";
