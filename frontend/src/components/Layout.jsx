import { LayoutDashboard, Settings, Activity } from 'lucide-react';

export const Layout = ({ children, current, setPage }) => {
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
                    </nav>
                </div>
            </div>

            {/* Content */}
            <main className="main-content">
                {children}
            </main>
        </div>
    );
};
