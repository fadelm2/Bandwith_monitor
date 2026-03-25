INSERT INTO `wan_capacities` (`wan_id`, `capacity_mbps`, `threshold_percent`, `description`, `created_at`) VALUES
('WAN-Primary', 1000, 80, 'Core Internet Link', NOW()),
('WAN-Secondary', 500, 75, 'Backup Link (LTE)', NOW()),
('WAN-Site-A', 100, 90, 'Branch Site A VPLS', NOW()),
('WAN-Site-B', 100, 90, 'Branch Site B VPLS', NOW());
