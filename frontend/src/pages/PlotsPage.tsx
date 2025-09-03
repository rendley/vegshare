import { useParams } from 'react-router-dom';
import { useGetPlotsByStructureQuery, useLeasePlotMutation } from '../features/api/apiSlice';

export const PlotsPage = () => {
  const { structureId } = useParams<{ structureId: string }>();
  if (!structureId) return <div>Ошибка: ID сооружения не указан.</div>;

  const { data: plots, error, isLoading } = useGetPlotsByStructureQuery(structureId);
  const [leasePlot, { isLoading: isLeasing }] = useLeasePlotMutation();

  const handleLeaseClick = async (plotId: string) => {
    try {
      await leasePlot({ unit_id: plotId, unit_type: 'plot' }).unwrap();
      // RTK Query автоматически обновит список благодаря тегам
      alert('Грядка успешно арендована!');
    } catch (err) {
      alert('Не удалось арендовать грядку. Ошибка: ' + JSON.stringify(err));
    }
  };

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
              {plot.status === 'available' && (
                <button onClick={() => handleLeaseClick(plot.id)} disabled={isLeasing}>
                  {isLeasing ? 'Аренда...' : 'Арендовать'}
                </button>
              )}
            </li>
          ))
        ) : (
          <li>Грядок в этом сооружении не найдено</li>
        )}
      </ul>
    </div>
  );
};