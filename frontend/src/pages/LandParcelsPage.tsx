import { useParams, Link } from 'react-router-dom';
import { useGetLandParcelsByRegionQuery } from '../features/api/apiSlice';

export const LandParcelsPage = () => {
  const { regionId } = useParams<{ regionId: string }>();
  if (!regionId) return <div>Ошибка: ID региона не указан.</div>;

  const { data: parcels, error, isLoading } = useGetLandParcelsByRegionQuery(regionId);

  if (isLoading) return <div>Загрузка участков...</div>;
  if (error) return <div>Ошибка при загрузке участков.</div>;

  return (
    <div>
      <h1>Земельные участки</h1>
      <ul>
        {parcels && parcels.length > 0 ? (
          parcels.map((parcel) => (
            <li key={parcel.id}>
              <Link to={`/land-parcels/${parcel.id}/structures`}>{parcel.name}</Link>
            </li>
          ))
        ) : (
          <li>Участков в этом регионе не найдено</li>
        )}
      </ul>
    </div>
  );
};