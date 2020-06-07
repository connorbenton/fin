BEGIN TRANSACTION;
DROP TABLE IF EXISTS `accounts`;
DROP TABLE IF EXISTS `categories`;
DROP TABLE IF EXISTS `salt_edge__categories`;
DROP TABLE IF EXISTS `plaid__categories`;
DROP TABLE IF EXISTS `currency_rates`;
DROP TABLE IF EXISTS `transactions`;
DROP TABLE IF EXISTS `analysis_trees`;
COMMIT;