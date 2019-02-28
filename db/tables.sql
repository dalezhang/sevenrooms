 -- ----------------------------------------------------------------------------
-- Schema sevenrooms
-- ----------------------------------------------------------------------------
-- DROP DATABASE IF EXISTS `sevenrooms`;

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tables` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `restaurant_party_id` int(11) DEFAULT NULL,
  `restaurant_table_id` int(11) DEFAULT NULL,
  `restaurant_updated_at` timestamp NULL DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `store_id` int(11) DEFAULT NULL,
  `check_ref` varchar(255) DEFAULT NULL,
  `check_open_time` int(64) DEFAULT NULL,
  `check_close_time` int(64) DEFAULT NULL,
  `table` varchar(255) DEFAULT NULL,
  `guest_count` int(11) DEFAULT NULL,
  `chind_count` int(11) DEFAULT NULL,
  `has_sync` tinyint(1) DEFAULT '0',
  `total_amount` decimal(12,2) DEFAULT NULL,
  `org_table` varchar(255) DEFAULT NULL,
  `order_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tables_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `transactions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `amount` decimal(16,2) DEFAULT '0.00',
  `gateway_transaction_id` int(11) DEFAULT NULL,
  `gateway_created_at` datetime DEFAULT NULL,
  `gateway_action` int(11) DEFAULT NULL,
  `gateway_voided_at` datetime DEFAULT NULL,
  `table_id` int(11) DEFAULT NULL,
  `has_sync` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_transactions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `line_items` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `gateway_line_item_id` int(11) DEFAULT NULL,
  `table_id` int(11) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `price` decimal(16,2) DEFAULT '0.00',
  `net_total` decimal(16,2) DEFAULT '0.00',
  `qty` decimal(12,4) DEFAULT NULL,
  `listing_barcode` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_line_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;