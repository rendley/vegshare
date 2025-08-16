import { useGetMyLeasesQuery } from '../features/api/apiSlice';

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
            <li key={lease.id}>
              ID Грядки: {lease.plot_id} (Арендована до: {new Date(lease.end_date).toLocaleDateString()})
            </li>
          ))}
        </ul>
      ) : (
        <p>У вас пока нет арендованных грядок.</p>
      )}
    </div>
  );
};
