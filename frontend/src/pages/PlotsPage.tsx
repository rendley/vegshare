import { useParams } from 'react-router-dom';
import { useGetPlotsByGreenhouseQuery } from '../features/api/apiSlice';

export const PlotsPage = () => {
  const { greenhouseId } = useParams<{ greenhouseId: string }>();
  if (!greenhouseId) return <div>Ошибка: ID теплицы не указан.</div>;

  const { data: plots, error, isLoading } = useGetPlotsByGreenhouseQuery(greenhouseId);

  if (isLoading) return <div>Загрузка грядок...</div>;
  if (error) return <div>Ошибка при загрузке грядок.</div>;

  return (
    <div>
      <h1>Грядки</h1>
      <ul>
        {plots && plots.length > 0 ? (
          plots.map((plot) => (
            <li key={plot.id}>
              {plot.name} - <span style={{ color: plot.status === 'available' ? 'green' : 'red' }}>{plot.status}</span>
            </li>
          ))
        ) : (
          <li>Грядок в этой теплице не найдено</li>
        )}
      </ul>
    </div>
  );
};
