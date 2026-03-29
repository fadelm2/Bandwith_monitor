-- Adding actual WAN IDs from Telegraf data
INSERT INTO `wan_capacities` (`wan_id`, `capacity_mbps`, `threshold_percent`, `description`, `created_at`) VALUES
('---WAN-123-TIS FO---', 1000, 80, 'Harapan Mulya Site', NOW()),
('WAN-1234-PENGGILINGAN_HSP FO', 500, 80, 'Penggilingan Site', NOW())
ON DUPLICATE KEY UPDATE capacity_mbps=VALUES(capacity_mbps), description=VALUES(description);
