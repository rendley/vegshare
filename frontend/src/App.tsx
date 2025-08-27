import { BrowserRouter as Router, Routes, Route, Link, useNavigate, Navigate } from 'react-router-dom';
import { RegionsPage } from './pages/RegionsPage';
import { LandParcelsPage } from './pages/LandParcelsPage';
import { GreenhousesPage } from './pages/GreenhousesPage';
import { PlotsPage } from './pages/PlotsPage';
import { MyPlotsPage } from './pages/MyPlotsPage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import AdminLayout from './pages/admin/AdminLayout';
import RegionManagementPage from './pages/admin/RegionManagementPage';
import { useDispatch, useSelector } from 'react-redux';
import { logout, selectCurrentUserRole, selectIsLoggedIn } from './features/auth/authSlice';

// Компонент для защиты роутов
const ProtectedRoute = ({ isAllowed, children }: { isAllowed: boolean, children: React.ReactNode }) => {
  if (!isAllowed) {
    return <Navigate to="/" replace />;
  }
  return <>{children}</>;
};

function App() {
  const isLoggedIn = useSelector(selectIsLoggedIn);
  const userRole = useSelector(selectCurrentUserRole);
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  return (
      <div>
        <nav>
          <ul>
            <li>
              <Link to="/">Главная</Link>
            </li>
            <li>
              <Link to="/regions">Регионы</Link>
            </li>
            <li>
              <Link to="/my-plots">Мои грядки</Link>
            </li>
            {userRole === 'admin' && (
              <li>
                <Link to="/admin/regions">Админка</Link>
              </li>
            )}
            {isLoggedIn ? (
              <li>
                <button onClick={handleLogout}>Выйти</button>
              </li>
            ) : (
              <li>
                <Link to="/login">Войти</Link>
              </li>
            )}
          </ul>
        </nav>

        <hr />

        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/regions" element={<RegionsPage />} />
          <Route path="/regions/:regionId/land-parcels" element={<LandParcelsPage />} />
          <Route path="/land-parcels/:parcelId/greenhouses" element={<GreenhousesPage />} />
          <Route path="/greenhouses/:greenhouseId/plots" element={<PlotsPage />} />
          <Route path="/my-plots" element={<MyPlotsPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          {/* Админские роуты */}
          <Route 
            path="/admin" 
            element={
              <ProtectedRoute isAllowed={userRole === 'admin'}>
                <AdminLayout />
              </ProtectedRoute>
            }
          >
            <Route path="regions" element={<RegionManagementPage />} />
            {/* Другие админские страницы будут здесь */}
          </Route>

        </Routes>
      </div>
  );
}

function Home() {
  return <h2>Главная страница</h2>;
}

const Root = () => <Router><App /></Router>;

export default Root;
