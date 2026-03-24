import React, { useState } from 'react';
import DashboardPage from './pages/DashboardPage';
import WanCapacityPage from './pages/WanCapacityPage';
import TrafficPage from './pages/TrafficPage';
import { Layout } from './components/Layout';

function App() {
  const [currentPage, setCurrentPage] = useState('dashboard');
  const [trafficFilter, setTrafficFilter] = useState('');

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

  return (
    <Layout current={currentPage} setPage={setCurrentPage}>
      {renderPage()}
    </Layout>
  );
}

export default App;
