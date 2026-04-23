INSERT INTO `wan_capacities` (`wan_id`, `capacity_mbps`, `threshold_percent`, `description`, `created_at`) VALUES
('WAN-001', 100, 80, 'Simulated Site 1', NOW()),
('WAN-002', 200, 80, 'Simulated Site 2', NOW()),
('WAN-003', 300, 80, 'Simulated Site 3', NOW()),
('WAN-004', 400, 80, 'Simulated Site 4', NOW()),
('WAN-005', 500, 80, 'Simulated Site 5', NOW())
ON DUPLICATE KEY UPDATE capacity_mbps=VALUES(capacity_mbps), description=VALUES(description);
