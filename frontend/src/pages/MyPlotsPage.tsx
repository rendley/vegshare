import { useGetMyLeasesQuery } from '../features/api/apiSlice';
import { PlantingControl } from '../features/plantings/PlantingControl';

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
              <p><b>Арендована до:</b> {new Date(lease.end_date).toLocaleDateString()}</p>
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