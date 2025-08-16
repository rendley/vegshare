import { Link } from 'react-router-dom';
import { useGetRegionsQuery } from '../features/api/apiSlice';

export const RegionsPage = () => {
  const { data: regions, error, isLoading } = useGetRegionsQuery();

  if (isLoading) return <div>Загрузка...</div>;
  if (error) return <div>Ошибка при загрузке регионов.</div>;

  return (
    <div>
      <h1>Регионы</h1>
      <ul>
        {regions && regions.length > 0 ? (
          regions.map((region) => (
            <li key={region.id}>
              <Link to={`/regions/${region.id}/land-parcels`}>{region.name}</Link>
            </li>
          ))
        ) : (
          <li>Регионов не найдено</li>
        )}
      </ul>
    </div>
  );
};