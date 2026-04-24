import React, { useState, useEffect } from 'react';
import { Plus, Trash2, Save, X, Settings, RefreshCw, FileCode, ArrowRightLeft } from 'lucide-react';
import { api } from '../api/client';

const TelegrafPage = () => {
    const [agents, setAgents] = useState([]);
    const [configPreview, setConfigPreview] = useState('');
    const [showModal, setShowModal] = useState(false);
    const [loading, setLoading] = useState(false);
    const [formData, setFormData] = useState({ ip_address: '', port: 161, protocol: 'udp', description: '' });

    const fetchData = async () => {
        setLoading(true);
        try {
            const [agentsRes, configRes] = await Promise.all([
                api.listTelegrafAgents(),
                api.getTelegrafConfig()
            ]);
            setAgents(agentsRes.data.data || []);
            setConfigPreview(configRes.data.data || '');
        } catch (err) {
            console.error("Failed to fetch Telegraf data", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { fetchData(); }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await api.createTelegrafAgent(formData);
            setShowModal(false);
            setFormData({ ip_address: '', port: 161, protocol: 'udp', description: '' });
            fetchData();
        } catch (err) {
            alert("Failed to add Telegraf agent");
        }
    };

    const handleImport = async () => {
        if (!window.confirm("Import agents from /etc/telegraf/telegraf.conf? Duplicate IPs will be skipped.")) return;
        setLoading(true);
        try {
            const res = await api.importTelegrafAgents();
            alert(`Import successful! Added ${res.data.data} new agents.`);
            fetchData();
        } catch (err) {
            alert("Failed to import agents from file. Check permissions.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
                <h1>Telegraf <span style={{ color: 'var(--accent-glow)' }}>Settings</span> MGMT</h1>
                <div style={{ display: 'flex', gap: '12px' }}>
                    <button className="btn" style={{ background: 'var(--bg-secondary)', border: '1px solid var(--border)' }} onClick={fetchData}>
                        <RefreshCw size={20} className={loading ? 'animate-spin' : ''} /> Refresh
                    </button>
                    <button className="btn" style={{ background: 'rgba(59, 130, 246, 0.1)', color: 'var(--accent-glow)', border: '1px solid var(--border)' }} onClick={handleImport}>
                        <ArrowRightLeft size={20} /> Sync from File
                    </button>
                    <button className="btn btn-primary" onClick={() => setShowModal(true)}>
                        <Plus size={20} /> Add SNMP Agent
                    </button>
                </div>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '24px' }}>
                {/* Agents Table */}
                <div className="glass-card">
                    <div style={{ padding: '16px', borderBottom: '1px solid var(--border)', display: 'flex', alignItems: 'center', gap: '8px' }}>
                        <Settings size={20} style={{ color: 'var(--accent-glow)' }} />
                        <h3 style={{ margin: 0 }}>SNMP Agents List</h3>
                    </div>
                    <div style={{ padding: '0px' }}>
                        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                            <thead>
                                <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--border)', fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
                                    <th style={{ padding: '12px 16px' }}>Address</th>
                                    <th style={{ padding: '12px 16px' }}>Protocol</th>
                                    <th style={{ padding: '12px 16px' }}>Description</th>
                                </tr>
                            </thead>
                            <tbody>
                                {agents.map((agent) => (
                                    <tr key={agent.id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                                        <td style={{ padding: '12px 16px', fontWeight: 'bold', color: 'var(--accent-glow)' }}>
                                            {agent.ip_address}:{agent.port}
                                        </td>
                                        <td style={{ padding: '12px 16px' }}>
                                            <span style={{ fontSize: '0.75rem', background: 'rgba(255,255,255,0.1)', padding: '2px 8px', borderRadius: '4px' }}>
                                                {agent.protocol.toUpperCase()}
                                            </span>
                                        </td>
                                        <td style={{ padding: '12px 16px', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>
                                            {agent.description || '—'}
                                        </td>
                                    </tr>
                                ))}
                                {agents.length === 0 && (
                                    <tr>
                                        <td colSpan="3" style={{ padding: '40px', textAlign: 'center', color: 'var(--text-secondary)' }}>
                                            No agents configured yet.
                                        </td>
                                    </tr>
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Config Preview */}
                <div className="glass-card" style={{ display: 'flex', flexDirection: 'column' }}>
                    <div style={{ padding: '16px', borderBottom: '1px solid var(--border)', display: 'flex', alignItems: 'center', gap: '8px' }}>
                        <FileCode size={20} style={{ color: '#fbbf24' }} />
                        <h3 style={{ margin: 0 }}>telegraf.conf Preview</h3>
                    </div>
                    <div style={{ flex: 1, padding: '16px', background: 'rgba(0,0,0,0.3)', margin: '16px', borderRadius: '8px', border: '1px solid var(--border)' }}>
                        <pre style={{ margin: 0, fontSize: '0.8rem', fontFamily: 'monospace', color: '#d1d5db', lineHeight: '1.5', whiteSpace: 'pre-wrap' }}>
                            {configPreview || '# Loading config...'}
                        </pre>
                    </div>
                    <div style={{ padding: '16px', borderTop: '1px solid var(--border)', fontSize: '0.75rem', color: 'var(--text-secondary)' }}>
                        Note: This is an automatically generated block based on the agents list.
                    </div>
                </div>
            </div>

            {/* Modal */}
            {showModal && (
                <div style={{
                    position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
                    background: 'rgba(0,0,0,0.8)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 100
                }}>
                    <div className="glass-card" style={{ width: '100%', maxWidth: '450px', position: 'relative' }}>
                        <button style={{ position: 'absolute', top: '16px', right: '16px', background: 'transparent', border: 'none', color: 'white', cursor: 'pointer' }}
                            onClick={() => setShowModal(false)}><X size={24} /></button>

                        <h2 style={{ marginBottom: '24px' }}>Add SNMP Agent</h2>

                        <form onSubmit={handleSubmit}>
                            <label style={{ display: 'block', marginBottom: '8px' }}>IP Address</label>
                            <input type="text" className="input-field" placeholder="e.g. 103.140.78.233" 
                                value={formData.ip_address}
                                onChange={(e) => setFormData({ ...formData, ip_address: e.target.value })} required />

                            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px', marginTop: '16px' }}>
                                <div>
                                    <label style={{ display: 'block', marginBottom: '8px' }}>Port</label>
                                    <input type="number" className="input-field" value={formData.port}
                                        onChange={(e) => setFormData({ ...formData, port: parseInt(e.target.value) })} required />
                                </div>
                                <div>
                                    <label style={{ display: 'block', marginBottom: '8px' }}>Protocol</label>
                                    <select className="input-field" value={formData.protocol}
                                        onChange={(e) => setFormData({ ...formData, protocol: e.target.value })}>
                                        <option value="udp">UDP</option>
                                        <option value="tcp">TCP</option>
                                    </select>
                                </div>
                            </div>

                            <label style={{ display: 'block', marginBottom: '8px', marginTop: '16px' }}>Description</label>
                            <textarea className="input-field" rows={2} value={formData.description}
                                placeholder="Optional name or location..."
                                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                style={{ resize: 'vertical', fontFamily: 'inherit' }} />

                            <button type="submit" className="btn btn-primary" style={{ width: '100%', marginTop: '24px' }}>
                                <Save size={20} /> Add SNMP Agent
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default TelegrafPage;
