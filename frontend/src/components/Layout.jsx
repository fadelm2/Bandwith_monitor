import { LayoutDashboard, Settings, Activity, LogOut } from 'lucide-react';

export const Layout = ({ children, current, setPage, user, onLogout }) => {
    const initial = user?.username ? user.username[0].toUpperCase() : '?';

    return (
        <div className="layout">
            {/* Sidebar */}
            <div className="sidebar glass-card">
                <div>
                    <h1 className="sidebar-logo">
                        WAN <span style={{ color: 'var(--accent-glow)' }}>System</span>
                    </h1>

                    <nav className="sidebar-nav">
                        <button
                            className={`btn ${current === 'dashboard' ? 'btn-primary' : ''}`}
                            style={{ width: '100%', background: current === 'dashboard' ? '' : 'transparent' }}
                            onClick={() => setPage('dashboard')}
                        >
                            <LayoutDashboard size={20} /> Dashboard
                        </button>
                        <button
                            className={`btn ${current === 'traffic' ? 'btn-primary' : ''}`}
                            style={{ width: '100%', background: current === 'traffic' ? '' : 'transparent' }}
                            onClick={() => setPage('traffic')}
                        >
                            <Activity size={20} /> Traffic Logs
                        </button>
                        <button
                            className={`btn ${current === 'capacity' ? 'btn-primary' : ''}`}
                            style={{ width: '100%', background: current === 'capacity' ? '' : 'transparent' }}
                            onClick={() => setPage('capacity')}
                        >
                            <Settings size={20} /> Capacity MGMT
                        </button>
                        <button
                            className={`btn ${current === 'telegraf' ? 'btn-primary' : ''}`}
                            style={{ width: '100%', background: current === 'telegraf' ? '' : 'transparent' }}
                            onClick={() => setPage('telegraf')}
                        >
                            <Settings size={20} /> Telegraf Settings
                        </button>
                    </nav>
                </div>

                {/* User info + Logout */}
                {user && (
                    <div style={userBoxStyle}>
                        <div style={avatarStyle}>{initial}</div>
                        <div style={{ flex: 1, minWidth: 0 }}>
                            <div style={{ fontWeight: 600, fontSize: '0.9rem', color: 'var(--text-primary)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                                {user.username}
                            </div>
                            <div style={{ fontSize: '0.75rem', color: 'var(--text-secondary)' }}>Online</div>
                        </div>
                        <button
                            onClick={onLogout}
                            title="Logout"
                            style={logoutBtnStyle}
                            onMouseEnter={e => e.currentTarget.style.background = 'rgba(239,68,68,0.2)'}
                            onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
                        >
                            <LogOut size={18} />
                        </button>
                    </div>
                )}
            </div>

            {/* Content */}
            <main className="main-content">
                {children}
            </main>
        </div>
    );
};

const userBoxStyle = {
    display: 'flex', alignItems: 'center', gap: 10,
    padding: '12px 14px',
    background: 'rgba(255,255,255,0.05)',
    borderRadius: 12,
    border: '1px solid rgba(255,255,255,0.08)',
    marginTop: 12,
};

const avatarStyle = {
    width: 36, height: 36, borderRadius: '50%',
    background: 'linear-gradient(135deg, #3b82f6, #22d3ee)',
    display: 'flex', alignItems: 'center', justifyContent: 'center',
    fontWeight: 700, fontSize: '1rem', color: 'white', flexShrink: 0,
};

const logoutBtnStyle = {
    background: 'transparent', border: 'none', cursor: 'pointer',
    color: '#ef4444', borderRadius: 8, padding: 6,
    display: 'flex', alignItems: 'center', justifyContent: 'center',
    transition: 'background 0.2s',
    flexShrink: 0,
};

