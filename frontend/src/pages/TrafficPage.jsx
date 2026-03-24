import React, { useState, useEffect } from 'react';
import { ChevronLeft, ChevronRight, Search } from 'lucide-react';
import { api } from '../api/client';

const TrafficPage = () => {
    const [traffic, setTraffic] = useState([]);
    const [paging, setPaging] = useState({ page: 1, size: 10, total_item: 0, total_page: 0 });
    const [wanFilter, setWanFilter] = useState('');

    const fetchTraffic = async (page = 1) => {
        try {
            const res = await api.searchTraffic({ page, size: 10, wan_id: wanFilter });
            setTraffic(res.data.data || []);
            if (res.data.paging) {
                setPaging(res.data.paging);
            }
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => { fetchTraffic(1); }, [wanFilter]);

    return (
        <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
                <h1>Network <span style={{ color: 'var(--accent-glow)' }}>Traffic</span> Logs</h1>
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center', background: 'var(--glass)', padding: '4px 12px', borderRadius: '12px', border: '1px solid var(--border)' }}>
                    <Search size={18} color="var(--text-secondary)" />
                    <input
                        type="text"
                        placeholder="Filter WAN ID... (e.g. 003 or wan-003)"
                        value={wanFilter}
                        onChange={(e) => {
                            let val = e.target.value.toUpperCase();
                            // If user types only digits, auto-prefix WAN-
                            if (/^\d+$/.test(val)) {
                                val = 'WAN-' + val;
                            }
                            setWanFilter(val);
                        }}
                        style={{ background: 'transparent', border: 'none', color: 'white', padding: '8px', outline: 'none', width: '200px' }}
                    />
                </div>
            </div>

            <div className="glass-card">
                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                    <thead>
                        <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--border)' }}>
                            <th style={{ padding: '16px' }}>Time</th>
                            <th style={{ padding: '16px' }}>WAN ID</th>
                            <th style={{ padding: '16px' }}>RX (Mbps)</th>
                            <th style={{ padding: '16px' }}>TX (Mbps)</th>
                            <th style={{ padding: '16px' }}>Utilization</th>
                        </tr>
                    </thead>
                    <tbody>
                        {traffic.map((item) => (
                            <tr key={item.id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                                <td style={{ padding: '16px', color: 'var(--text-secondary)' }}>{new Date(item.created_at).toLocaleString()}</td>
                                <td style={{ padding: '16px', fontWeight: 'bold' }}>{item.wan_id}</td>
                                <td style={{ padding: '16px', color: 'var(--accent-glow)' }}>{item.rx_mbps.toFixed(2)}</td>
                                <td style={{ padding: '16px', color: 'var(--accent-color)' }}>{item.tx_mbps.toFixed(2)}</td>
                                <td style={{ padding: '16px' }}>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                        <div style={{ width: '60px', height: '6px', background: 'rgba(255,255,255,0.1)', borderRadius: '3px', overflow: 'hidden' }}>
                                            <div style={{ width: `${Math.min(item.utilization_percent, 100)}%`, height: '100%', background: item.utilization_percent > 80 ? 'var(--danger)' : 'var(--success)' }}></div>
                                        </div>
                                        <span>{item.utilization_percent.toFixed(1)}%</span>
                                    </div>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>

                {/* Pagination */}
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '24px', padding: '0 8px' }}>
                    <p style={{ color: 'var(--text-secondary)' }}>Showing {traffic.length} of {paging.total_item} records</p>
                    <div style={{ display: 'flex', gap: '8px' }}>
                        <button className="btn" disabled={paging.page <= 1} onClick={() => fetchTraffic(paging.page - 1)}
                            style={{ background: 'var(--glass)', border: '1px solid var(--border)' }}>
                            <ChevronLeft size={20} /> Prev
                        </button>
                        <div className="btn" style={{ background: 'var(--accent-color)' }}>{paging.page}</div>
                        <button className="btn" disabled={paging.page >= paging.total_page} onClick={() => fetchTraffic(paging.page + 1)}
                            style={{ background: 'var(--glass)', border: '1px solid var(--border)' }}>
                            Next <ChevronRight size={20} />
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default TrafficPage;
