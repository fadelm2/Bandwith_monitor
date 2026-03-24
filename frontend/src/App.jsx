import React, { useState } from 'react';
import DashboardPage from './pages/DashboardPage';
import WanCapacityPage from './pages/WanCapacityPage';
import TrafficPage from './pages/TrafficPage';
import { Layout } from './components/Layout';

function App() {
  const [currentPage, setCurrentPage] = useState('dashboard');

  const renderPage = () => {
    switch (currentPage) {
      case 'dashboard':
        return <DashboardPage />;
      case 'traffic':
        return <TrafficPage />;
      case 'capacity':
        return <WanCapacityPage />;
      default:
        return <DashboardPage />;
    }
  };

  return (
    <Layout current={currentPage} setPage={setCurrentPage}>
      {renderPage()}
    </Layout>
  );
}

export default App;
