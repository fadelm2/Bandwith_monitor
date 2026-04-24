CREATE TABLE IF NOT EXISTS wan_capacities (
    wan_id VARCHAR(100) PRIMARY KEY,
    capacity_mbps DOUBLE NOT NULL,
    threshold_percent DOUBLE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wan_traffics (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    wan_id VARCHAR(100) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    interface_name VARCHAR(255) NOT NULL,
    rx_mbps DOUBLE NOT NULL,
    tx_mbps DOUBLE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_wan_traffic_wan_id (wan_id),
    INDEX idx_wan_traffic_created_at (created_at)
);
