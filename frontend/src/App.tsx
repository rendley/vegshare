import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { selectCurrentUserRole, selectIsLoggedIn } from './features/auth/authSlice';

import { Layout } from './components/Layout';
import { MarketplacePage } from './pages/MarketplacePage';
import { DashboardPage } from './pages/DashboardPage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import AdminLayout from './pages/admin/AdminLayout';
import UserManagementPage from './pages/admin/UserManagementPage';
import RegionManagementPage from './pages/admin/RegionManagementPage';
import LandParcelManagementPage from './pages/admin/LandParcelManagementPage';
import GreenhouseManagementPage from './pages/admin/GreenhouseManagementPage';
import PlotManagementPage from './pages/admin/PlotManagementPage';
import CameraManagementPage from './pages/admin/CameraManagementPage';
import CropManagementPage from './pages/admin/CropManagementPage';

// Компонент для защиты роутов, который перенаправляет на указанный путь
const ProtectedRoute = ({ isAllowed, redirectTo = "/", children }: { isAllowed: boolean, redirectTo?: string, children: React.ReactNode }) => {
  if (!isAllowed) {
    return <Navigate to={redirectTo} replace />;
  }
  return <>{children}</>;
};

function App() {
  const isLoggedIn = useSelector(selectIsLoggedIn);
  const userRole = useSelector(selectCurrentUserRole);

  return (
    <Router>
      <Routes>
        <Route element={<Layout />}>
          {/* Public Routes */}
          <Route index element={<MarketplacePage />} />
          <Route path="login" element={<LoginPage />} />
          <Route path="register" element={<RegisterPage />} />

          {/* Protected User Routes */}
          <Route 
            path="dashboard" 
            element={
              <ProtectedRoute isAllowed={isLoggedIn} redirectTo="/login">
                <DashboardPage />
              </ProtectedRoute>
            } 
          />

          {/* Protected Admin Routes */}
          <Route 
            path="admin" 
            element={
              <ProtectedRoute isAllowed={userRole === 'admin'}>
                <AdminLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Navigate to="users" replace />} />
            <Route path="users" element={<UserManagementPage />} />
            <Route path="regions" element={<RegionManagementPage />} />
            <Route path="parcels" element={<LandParcelManagementPage />} />
            <Route path="greenhouses" element={<GreenhouseManagementPage />} />
            <Route path="plots" element={<PlotManagementPage />} />
            <Route path="cameras" element={<CameraManagementPage />} />
            <Route path="crops" element={<CropManagementPage />} />
          </Route>
          
          {/* Fallback for any other route */}
          <Route path="*" element={<Navigate to="/" />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;