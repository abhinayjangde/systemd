/*
  Warnings:

  - You are about to drop the column `webhookId` on the `Event` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "Event" DROP CONSTRAINT "Event_webhookId_fkey";

-- AlterTable
ALTER TABLE "Event" DROP COLUMN "webhookId";

-- CreateTable
CREATE TABLE "WebhookDelivery" (
    "id" TEXT NOT NULL,
    "eventId" TEXT NOT NULL,
    "webhookId" TEXT NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'PENDING',
    "attempts" INTEGER NOT NULL DEFAULT 0,
    "lastError" TEXT,
    "deliveredAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "WebhookDelivery_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "WebhookDelivery_eventId_webhookId_key" ON "WebhookDelivery"("eventId", "webhookId");

-- AddForeignKey
ALTER TABLE "WebhookDelivery" ADD CONSTRAINT "WebhookDelivery_eventId_fkey" FOREIGN KEY ("eventId") REFERENCES "Event"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "WebhookDelivery" ADD CONSTRAINT "WebhookDelivery_webhookId_fkey" FOREIGN KEY ("webhookId") REFERENCES "WebhookEndpoint"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
