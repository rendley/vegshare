import { useGetMyLeasesQuery, useGetCamerasByPlotQuery } from '../features/api/apiSlice';
import PlantingControl from '../features/plantings/PlantingControl';
import VideoPlayer from '../components/VideoPlayer';

// Новый компонент для отображения камер грядки
const PlotCameraViewer = ({ plotId }: { plotId: string }) => {
  const { data: cameras, isLoading, error } = useGetCamerasByPlotQuery(plotId);

  if (isLoading) return <p>Загрузка камер...</p>;
  if (error) return <p style={{ color: 'orange' }}>Не удалось загрузить камеры.</p>;

  return (
    <div>
      {cameras && cameras.length > 0 ? (
        cameras.map(camera => <VideoPlayer key={camera.id} camera={camera} />)
      ) : (
        <p>Камеры не найдены.</p>
      )}
    </div>
  );
}

export const MyPlotsPage = () => {
  const { data: leases, error, isLoading } = useGetMyLeasesQuery();

  if (isLoading) return <div>Загрузка ваших аренд...</div>;
  if (error) return <div>Ошибка при загрузке аренд.</div>;

  return (
    <div>
      <h1>Мои грядки</h1>
      {leases && leases.length > 0 ? (
        <ul>
          {leases.map((lease) => (
            <li key={lease.id} style={{ border: '1px solid #ccc', padding: '10px', marginBottom: '10px' }}>
              <p><b>ID Грядки:</b> {lease.plot_id}</p>
              <hr />
              <PlotCameraViewer plotId={lease.plot_id} />
              <hr />
              <div>
                <b>Управление посадкой:</b>
                <PlantingControl plotId={lease.plot_id} />
              </div>
            </li>
          ))}
        </ul>
      ) : (
        <p>У вас пока нет арендованных грядок.</p>
      )}
    </div>
  );
};