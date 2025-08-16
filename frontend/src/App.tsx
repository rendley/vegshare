import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { RegionsPage } from './pages/RegionsPage';
import { LandParcelsPage } from './pages/LandParcelsPage';
import { GreenhousesPage } from './pages/GreenhousesPage';
import { PlotsPage } from './pages/PlotsPage';
import { MyPlotsPage } from './pages/MyPlotsPage';

function App() {
  return (
    <Router>
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
        </Routes>
      </div>
    </Router>
  );
}

function Home() {
  return <h2>Главная страница</h2>;
}

export default App;
