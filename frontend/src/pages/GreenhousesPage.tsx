import { useParams, Link } from 'react-router-dom';
import { useGetGreenhousesByLandParcelQuery } from '../features/api/apiSlice';

export const GreenhousesPage = () => {
  const { parcelId } = useParams<{ parcelId: string }>();
  if (!parcelId) return <div>Ошибка: ID участка не указан.</div>;

  const { data: greenhouses, error, isLoading } = useGetGreenhousesByLandParcelQuery(parcelId);

  if (isLoading) return <div>Загрузка теплиц...</div>;
  if (error) return <div>Ошибка при загрузке теплиц.</div>;

  return (
    <div>
      <h1>Теплицы</h1>
      <ul>
        {greenhouses && greenhouses.length > 0 ? (
          greenhouses.map((greenhouse) => (
            <li key={greenhouse.id}>
              <Link to={`/greenhouses/${greenhouse.id}/plots`}>{greenhouse.name}</Link>
            </li>
          ))
        ) : (
          <li>Теплиц на этом участке не найдено</li>
        )}
      </ul>
    </div>
  );
};
