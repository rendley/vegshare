import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from 'react-router-dom';
import { RegionsPage } from './pages/RegionsPage';
import { LandParcelsPage } from './pages/LandParcelsPage';
import { GreenhousesPage } from './pages/GreenhousesPage';
import { PlotsPage } from './pages/PlotsPage';
import { MyPlotsPage } from './pages/MyPlotsPage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';

function App() {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token');
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
            {token ? (
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
        </Routes>
      </div>
  );
}

function Home() {
  return <h2>Главная страница</h2>;
}

const Root = () => <Router><App /></Router>;

export default Root;
