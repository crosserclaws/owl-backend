TRUNCATE TABLE `hosts`;
ALTER TABLE `hosts` ADD COLUMN `bonding` int(3) UNSIGNED DEFAULT NULL;
ALTER TABLE `hosts` ADD COLUMN `speed` int(8) UNSIGNED DEFAULT NULL;
