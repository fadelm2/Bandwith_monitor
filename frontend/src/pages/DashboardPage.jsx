import React, { useState, useEffect } from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { ArrowDown, ArrowUp, Activity } from 'lucide-react';
import { api } from '../api/client';

const DashboardPage = () => {
    const [trafficData, setTrafficData] = useState([]);
    const [stats, setStats] = useState({ rx: 0, tx: 0 });

    useEffect(() => {
        const fetchData = async () => {
            try {
                const res = await api.searchTraffic({ page: 1, size: 20 });
                const data = res.data.data || [];

                // Reverse for chart (oldest to newest)
                const chartData = data.slice().reverse().map(item => ({
                    time: new Date(item.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }),
                    rx: parseFloat(item.rx_mbps.toFixed(2)),
                    tx: parseFloat(item.tx_mbps.toFixed(2)),
                }));

                setTrafficData(chartData);

                // Update stats from latest
                if (data.length > 0) {
                    setStats({ rx: data[0].rx_mbps, tx: data[0].tx_mbps });
                }
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

                <div className="glass-card" style={{ borderLeft: '4px solid var(--success)' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div>
                            <p style={{ color: 'var(--text-secondary)', marginBottom: '4px' }}>System Status</p>
                            <h3>Healthy</h3>
                        </div>
                        <Activity color="var(--success)" size={32} />
                    </div>
                </div>
            </div>

            <div className="glass-card" style={{ height: '400px', marginTop: '24px' }}>
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
