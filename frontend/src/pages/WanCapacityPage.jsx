import React, { useState, useEffect } from 'react';
import { Plus, Trash2, Edit2, Download, Save, X, ChevronLeft, ChevronRight } from 'lucide-react';
import { api } from '../api/client';

const PAGE_SIZE = 10;

const WanCapacityPage = () => {
    const [capacities, setCapacities] = useState([]);
    const [showModal, setShowModal] = useState(false);
    const [showBulkModal, setShowBulkModal] = useState(false);
    const [bulkData, setBulkData] = useState('');
    const [editing, setEditing] = useState(null);
    const [formData, setFormData] = useState({ wan_id: '', capacity_mbps: 0, threshold_percent: 0, description: '' });
    const [page, setPage] = useState(1);
    const [search, setSearch] = useState('');

    const fetchCapacities = async () => {
        try {
            const res = await api.listCapacity();
            setCapacities(res.data.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    useEffect(() => { fetchCapacities(); }, []);

    const filtered = capacities.filter(item =>
        item.wan_id.toLowerCase().includes(search.toLowerCase()) ||
        (item.description || '').toLowerCase().includes(search.toLowerCase())
    );
    const totalPages = Math.ceil(filtered.length / PAGE_SIZE);
    const paginated = filtered.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (editing) {
                await api.updateCapacity(editing.wan_id, formData);
            } else {
                await api.createCapacity(formData);
            }
            setShowModal(false);
            setEditing(null);
            fetchCapacities();
        } catch (err) { alert("Failed to save capacity"); }
    };

    const handleDelete = async (id) => {
        if (window.confirm("Delete this WAN capacity?")) {
            try {
                await api.deleteCapacity(id);
                fetchCapacities();
            } catch (err) { console.error(err); }
        }
    };

    return (
        <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px' }}>
                <h1>WAN <span style={{ color: 'var(--accent-glow)' }}>Capacity</span> MGMT</h1>
                <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', background: 'var(--glass)', border: '1px solid var(--border)', borderRadius: '8px', padding: '8px 12px' }}>
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="var(--text-secondary)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8" /><line x1="21" y1="21" x2="16.65" y2="16.65" /></svg>
                        <input
                            type="text"
                            placeholder="Filter by WAN ID or Description..."
                            value={search}
                            onChange={(e) => { setSearch(e.target.value); setPage(1); }}
                            style={{ background: 'transparent', border: 'none', color: 'white', outline: 'none', width: '240px', fontSize: '0.875rem' }}
                        />
                        {search && (
                            <button onClick={() => { setSearch(''); setPage(1); }} style={{ background: 'none', border: 'none', color: 'var(--text-secondary)', cursor: 'pointer', padding: '0' }}>✕</button>
                        )}
                    </div>
                    <button className="btn" style={{ background: 'var(--bg-secondary)', border: '1px solid var(--border)' }} onClick={() => setShowBulkModal(true)}>
                        <Download size={20} /> Bulk Update
                    </button>
                    <button className="btn btn-primary" onClick={() => { setEditing(null); setFormData({ wan_id: '', capacity_mbps: 0, threshold_percent: 0, description: '' }); setShowModal(true); }}>
                        <Plus size={20} /> Add Capacity
                    </button>
                </div>
            </div>

            <div className="glass-card">
                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                    <thead>
                        <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--border)' }}>
                            <th style={{ padding: '16px' }}>WAN ID</th>
                            <th style={{ padding: '16px' }}>Capacity (Mbps)</th>
                            <th style={{ padding: '16px' }}>Threshold (%)</th>
                            <th style={{ padding: '16px' }}>Description</th>
                            <th style={{ padding: '16px', textAlign: 'right' }}>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {paginated.map((item) => (
                            <tr key={item.wan_id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                                <td style={{ padding: '16px', fontWeight: 'bold' }}>{item.wan_id}</td>
                                <td style={{ padding: '16px' }}>{item.capacity_mbps}</td>
                                <td style={{ padding: '16px' }}>{item.threshold_percent}%</td>
                                <td style={{ padding: '16px', color: 'var(--text-secondary)', maxWidth: '220px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                                    {item.description || <span style={{ fontStyle: 'italic', opacity: 0.4 }}>—</span>}
                                </td>
                                <td style={{ padding: '16px', textAlign: 'right', display: 'flex', justifyContent: 'flex-end', gap: '8px' }}>
                                    <button className="btn" style={{ background: 'rgba(59, 130, 246, 0.1)', color: 'var(--accent-glow)' }}
                                        onClick={() => { setEditing(item); setFormData(item); setShowModal(true); }}>
                                        <Edit2 size={16} />
                                    </button>
                                    <button className="btn" style={{ background: 'rgba(239, 68, 68, 0.1)', color: 'var(--danger)' }}
                                        onClick={() => handleDelete(item.wan_id)}>
                                        <Trash2 size={16} />
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>

                {/* Pagination */}
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '16px', borderTop: '1px solid var(--border)' }}>
                    <p style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>
                        Showing {filtered.length === 0 ? 0 : (page - 1) * PAGE_SIZE + 1}–{Math.min(page * PAGE_SIZE, filtered.length)} of {filtered.length} entries
                        {search && <span style={{ color: 'var(--accent-glow)', marginLeft: '4px' }}>(filtered from {capacities.length} total)</span>}
                    </p>
                    <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                        <button className="btn" disabled={page <= 1} onClick={() => setPage(p => p - 1)}
                            style={{ background: 'var(--glass)', border: '1px solid var(--border)', opacity: page <= 1 ? 0.4 : 1 }}>
                            <ChevronLeft size={18} /> Prev
                        </button>
                        {Array.from({ length: totalPages }, (_, i) => i + 1).map(p => (
                            <button key={p} className="btn" onClick={() => setPage(p)}
                                style={{
                                    background: p === page ? 'var(--accent-glow)' : 'var(--glass)',
                                    border: '1px solid var(--border)',
                                    color: p === page ? '#000' : 'white',
                                    fontWeight: p === page ? 'bold' : 'normal',
                                    minWidth: '36px',
                                }}>
                                {p}
                            </button>
                        ))}
                        <button className="btn" disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}
                            style={{ background: 'var(--glass)', border: '1px solid var(--border)', opacity: page >= totalPages ? 0.4 : 1 }}>
                            Next <ChevronRight size={18} />
                        </button>
                    </div>
                </div>
            </div>

            {showModal && (
                <div style={{
                    position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
                    background: 'rgba(0,0,0,0.8)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 100
                }}>
                    <div className="glass-card" style={{ width: '100%', maxWidth: '450px', position: 'relative' }}>
                        <button style={{ position: 'absolute', top: '16px', right: '16px', background: 'transparent', border: 'none', color: 'white', cursor: 'pointer' }}
                            onClick={() => setShowModal(false)}><X size={24} /></button>

                        <h2 style={{ marginBottom: '24px' }}>{editing ? 'Edit' : 'Add'} WAN Capacity</h2>

                        <form onSubmit={handleSubmit}>
                            <label style={{ display: 'block', marginBottom: '8px' }}>WAN ID</label>
                            <input type="text" className="input-field" value={formData.wan_id} disabled={!!editing}
                                onChange={(e) => setFormData({ ...formData, wan_id: e.target.value })} required />

                            <label style={{ display: 'block', marginBottom: '8px', marginTop: '16px' }}>Capacity (Mbps)</label>
                            <input type="number" className="input-field" value={formData.capacity_mbps}
                                onChange={(e) => setFormData({ ...formData, capacity_mbps: parseFloat(e.target.value) })} required />

                            <label style={{ display: 'block', marginBottom: '8px', marginTop: '16px' }}>Threshold (%)</label>
                            <input type="number" className="input-field" value={formData.threshold_percent}
                                onChange={(e) => setFormData({ ...formData, threshold_percent: parseFloat(e.target.value) })} required />

                            <label style={{ display: 'block', marginBottom: '8px', marginTop: '16px' }}>Description</label>
                            <textarea className="input-field" rows={3} value={formData.description || ''}
                                placeholder="Optional description for this WAN link..."
                                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                                style={{ resize: 'vertical', fontFamily: 'inherit' }} />

                            <button type="submit" className="btn btn-primary" style={{ width: '100%', marginTop: '16px' }}>
                                <Save size={20} /> {editing ? 'Update' : 'Create'} Capacity
                            </button>
                        </form>
                    </div>
                </div>
            )}

            {showBulkModal && (
                <div style={{
                    position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
                    background: 'rgba(0,0,0,0.8)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 100
                }}>
                    <div className="glass-card" style={{ width: '100%', maxWidth: '600px', position: 'relative' }}>
                        <button style={{ position: 'absolute', top: '16px', right: '16px', background: 'transparent', border: 'none', color: 'white', cursor: 'pointer' }}
                            onClick={() => setShowBulkModal(false)}><X size={24} /></button>

                        <h2 style={{ marginBottom: '24px' }}>Bulk Update Capacity</h2>
                        <p style={{ color: 'var(--text-secondary)', marginBottom: '16px' }}>Paste JSON array of capacities:</p>

                        <textarea
                            className="input-field"
                            style={{ height: '200px', fontFamily: 'monospace' }}
                            placeholder='[{"wan_id": "WAN1", "capacity_mbps": 100, "threshold_percent": 80, "description": "Main link"}]'
                            value={bulkData}
                            onChange={(e) => setBulkData(e.target.value)}
                        />

                        <button className="btn btn-primary" style={{ width: '100%', marginTop: '16px' }}
                            onClick={async () => {
                                try {
                                    const data = JSON.parse(bulkData);
                                    await api.bulkUpdateCapacity(data);
                                    setShowBulkModal(false);
                                    fetchCapacities();
                                } catch (err) { alert("Invalid JSON or Update Failed"); }
                            }}>
                            <Save size={20} /> Execute Bulk Update
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default WanCapacityPage;
