import React, { useState, useEffect } from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { ArrowDown, ArrowUp, Activity, AlertTriangle, ChevronRight } from 'lucide-react';
import { api } from '../api/client';

const DashboardPage = ({ onAlertClick }) => {
    const [trafficData, setTrafficData] = useState([]);
    const [stats, setStats] = useState({ rx: 0, tx: 0 });
    const [alerts, setAlerts] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const res = await api.searchTraffic({ page: 1, size: 100 });
                const data = res.data.data || [];

                // Chart: last 20, reversed oldest→newest
                const chartData = data.slice(0, 20).slice().reverse().map(item => ({
                    time: new Date(item.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }),
                    rx: parseFloat(item.rx_mbps.toFixed(2)),
                    tx: parseFloat(item.tx_mbps.toFixed(2)),
                }));
                setTrafficData(chartData);

                // Stats from latest
                if (data.length > 0) {
                    setStats({ rx: data[0].rx_mbps, tx: data[0].tx_mbps });
                }

                // Alerts: Get latest per WAN, check if > threshold (80%)
                const latestPerWan = {};
                data.forEach(item => {
                    if (!latestPerWan[item.wan_id]) {
                        latestPerWan[item.wan_id] = item;
                    }
                });
                const highUtil = Object.values(latestPerWan)
                    .filter(item => item.utilization_percent > 80)
                    .sort((a, b) => b.utilization_percent - a.utilization_percent)
                    .slice(0, 10); // Top 10 most critical
                setAlerts(highUtil);
            } catch (err) {
                console.error("Failed to fetch traffic", err);
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 5000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div>
            <h1 style={{ marginBottom: '32px' }}>Network <span style={{ color: 'var(--accent-glow)' }}>Dashboard</span></h1>

            <div className="dashboard-grid" style={{ padding: 0, marginBottom: '24px' }}>
                <div className="glass-card" style={{ borderLeft: '4px solid var(--accent-glow)' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div>
                            <p style={{ color: 'var(--text-secondary)', marginBottom: '4px' }}>Download Speed</p>
                            <h3>{stats.rx.toFixed(2)} <span style={{ fontSize: '1rem' }}>Mbps</span></h3>
                        </div>
                        <ArrowDown color="var(--accent-glow)" size={32} />
                    </div>
                </div>

                <div className="glass-card" style={{ borderLeft: '4px solid var(--accent-color)' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div>
                            <p style={{ color: 'var(--text-secondary)', marginBottom: '4px' }}>Upload Speed</p>
                            <h3>{stats.tx.toFixed(2)} <span style={{ fontSize: '1rem' }}>Mbps</span></h3>
                        </div>
                        <ArrowUp color="var(--accent-color)" size={32} />
                    </div>
                </div>

                <div className="glass-card" style={{ borderLeft: `4px solid ${alerts.length > 0 ? 'var(--danger)' : 'var(--success)'}` }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div>
                            <p style={{ color: 'var(--text-secondary)', marginBottom: '4px' }}>System Status</p>
                            <h3 style={{ color: alerts.length > 0 ? 'var(--danger)' : 'var(--success)' }}>
                                {alerts.length > 0 ? `${alerts.length} Alert${alerts.length > 1 ? 's' : ''}` : 'Healthy'}
                            </h3>
                        </div>
                        {alerts.length > 0
                            ? <AlertTriangle color="var(--danger)" size={32} />
                            : <Activity color="var(--success)" size={32} />
                        }
                    </div>
                </div>
            </div>

            {/* High Utilization Alerts */}
            {alerts.length > 0 && (
                <div className="glass-card" style={{ marginBottom: '24px', borderLeft: '4px solid var(--danger)' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                        <h3 style={{ color: 'var(--danger)', display: 'flex', alignItems: 'center', gap: '8px' }}>
                            <AlertTriangle size={20} /> High Utilization Alerts
                            <span style={{ fontSize: '0.75rem', background: 'var(--danger)', color: 'white', borderRadius: '12px', padding: '2px 8px', marginLeft: '4px' }}>
                                {alerts.length} sites
                            </span>
                        </h3>
                        <button
                            className="btn"
                            onClick={() => onAlertClick && onAlertClick('')}
                            style={{ background: 'rgba(255,100,100,0.15)', border: '1px solid var(--danger)', color: 'var(--danger)', fontSize: '0.8rem', padding: '6px 12px' }}
                        >
                            View All in Logs →
                        </button>
                    </div>
                    <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                        <thead>
                            <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--border)' }}>
                                <th style={{ padding: '10px 16px', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>WAN ID</th>
                                <th style={{ padding: '10px 16px', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>Utilization</th>
                                <th style={{ padding: '10px 16px', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>RX (Mbps)</th>
                                <th style={{ padding: '10px 16px', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>TX (Mbps)</th>
                                <th style={{ padding: '10px 16px' }}></th>
                            </tr>
                        </thead>
                        <tbody>
                            {alerts.map(item => (
                                <tr
                                    key={item.wan_id}
                                    onClick={() => onAlertClick && onAlertClick(item.wan_id)}
                                    style={{
                                        borderBottom: '1px solid rgba(255,255,255,0.05)',
                                        cursor: 'pointer',
                                        transition: 'background 0.2s',
                                    }}
                                    onMouseEnter={e => e.currentTarget.style.background = 'rgba(255,100,100,0.08)'}
                                    onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
                                >
                                    <td style={{ padding: '12px 16px', fontWeight: 'bold' }}>{item.wan_id}</td>
                                    <td style={{ padding: '12px 16px' }}>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                            <div style={{ width: '80px', height: '6px', background: 'rgba(255,255,255,0.1)', borderRadius: '3px', overflow: 'hidden' }}>
                                                <div style={{ width: `${Math.min(item.utilization_percent, 100)}%`, height: '100%', background: 'var(--danger)' }}></div>
                                            </div>
                                            <span style={{ color: 'var(--danger)', fontWeight: 'bold' }}>{item.utilization_percent.toFixed(1)}%</span>
                                        </div>
                                    </td>
                                    <td style={{ padding: '12px 16px', color: 'var(--accent-glow)' }}>{item.rx_mbps.toFixed(2)}</td>
                                    <td style={{ padding: '12px 16px', color: 'var(--accent-color)' }}>{item.tx_mbps.toFixed(2)}</td>
                                    <td style={{ padding: '12px 16px', color: 'var(--text-secondary)' }}>
                                        <ChevronRight size={18} />
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}

            <div className="glass-card" style={{ height: '400px' }}>
                <h3 style={{ marginBottom: '24px' }}>Real-time Traffic History</h3>
                <ResponsiveContainer width="100%" height="80%">
                    <AreaChart data={trafficData}>
                        <defs>
                            <linearGradient id="colorRx" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="var(--accent-glow)" stopOpacity={0.3} />
                                <stop offset="95%" stopColor="var(--accent-glow)" stopOpacity={0} />
                            </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.05)" vertical={false} />
                        <XAxis dataKey="time" stroke="var(--text-secondary)" fontSize={12} tickLine={false} axisLine={false} />
                        <YAxis stroke="var(--text-secondary)" fontSize={12} tickLine={false} axisLine={false} mirror unit="M" />
                        <Tooltip
                            contentStyle={{ background: 'var(--bg-secondary)', border: '1px solid var(--border)', borderRadius: '8px' }}
                            itemStyle={{ color: 'white' }}
                        />
                        <Area type="monotone" dataKey="rx" stroke="var(--accent-glow)" fillOpacity={1} fill="url(#colorRx)" />
                        <Area type="monotone" dataKey="tx" stroke="var(--accent-color)" fill="transparent" />
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
};

export default DashboardPage;
