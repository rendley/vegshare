import { useParams, Link } from 'react-router-dom';
import { useGetStructuresByLandParcelQuery } from '../features/api/apiSlice';

export const StructuresPage = () => {
  const { parcelId } = useParams<{ parcelId: string }>();
  if (!parcelId) return <div>Ошибка: ID участка не указан.</div>;

  const { data: structures, error, isLoading } = useGetStructuresByLandParcelQuery(parcelId);

  if (isLoading) return <div>Загрузка сооружений...</div>;
  if (error) return <div>Ошибка при загрузке сооружений.</div>;

  return (
    <div>
      <h1>Сооружения</h1>
      <ul>
        {structures && structures.length > 0 ? (
          structures.map((structure) => (
            <li key={structure.id}>
              <Link to={`/structures/${structure.id}/plots`}>{structure.name}</Link>
            </li>
          ))
        ) : (
          <li>Сооружений на этом участке не найдено</li>
        )}
      </ul>
    </div>
  );
};