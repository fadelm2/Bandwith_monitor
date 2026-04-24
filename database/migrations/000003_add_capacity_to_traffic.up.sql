ALTER TABLE wan_traffics ADD COLUMN capacity_mbps DOUBLE NOT NULL AFTER tx_mbps;
ALTER TABLE wan_traffics ADD COLUMN utilization_percent DOUBLE NOT NULL AFTER capacity_mbps;
