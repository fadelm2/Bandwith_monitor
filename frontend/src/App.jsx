import React, { useState, useEffect } from 'react';
import DashboardPage from './pages/DashboardPage';
import WanCapacityPage from './pages/WanCapacityPage';
import TrafficPage from './pages/TrafficPage';
import LoginPage from './pages/LoginPage';
import { Layout } from './components/Layout';
import { authApi } from './api/auth';

function App() {
  const [currentPage, setCurrentPage] = useState('dashboard');
  const [trafficFilter, setTrafficFilter] = useState('');
  const [user, setUser] = useState(null);
  const [checkingAuth, setCheckingAuth] = useState(true);

  // Check existing token on mount
  useEffect(() => {
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    if (token && username) {
      setUser({ username });
    }
    setCheckingAuth(false);
  }, []);

  const handleLogin = (userData) => {
    setUser(userData);
    setCurrentPage('dashboard');
  };

  const handleLogout = async () => {
    try {
      await authApi.logout();
    } catch (_) {
      // silently ignore logout errors
    } finally {
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      setUser(null);
    }
  };

  const navigateToTraffic = (wanId = '') => {
    setTrafficFilter(wanId);
    setCurrentPage('traffic');
  };

  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard':
        return <DashboardPage onAlertClick={navigateToTraffic} />;
      case 'traffic':
        return <TrafficPage initialFilter={trafficFilter} />;
      case 'capacity':
        return <WanCapacityPage />;
      default:
        return <DashboardPage onAlertClick={navigateToTraffic} />;
    }
  };

  // Loading splash
  if (checkingAuth) {
    return (
      <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'var(--bg-color)' }}>
        <div style={{ width: 40, height: 40, borderRadius: '50%', border: '3px solid rgba(255,255,255,0.1)', borderTopColor: '#22d3ee', animation: 'spin 0.7s linear infinite' }} />
        <style>{`@keyframes spin { to { transform: rotate(360deg); } }`}</style>
      </div>
    );
  }

  if (!user) {
    return <LoginPage onLogin={handleLogin} />;
  }

  return (
    <Layout current={currentPage} setPage={setCurrentPage} user={user} onLogout={handleLogout}>
      {renderPage()}
    </Layout>
  );
}

export default App;

