import { useState } from 'react';
import { useGetAvailableCropsQuery, usePlantCropMutation, useGetPlotCropsQuery } from './apiSlice';

interface PlantingControlProps {
  plotId: string;
}

export const PlantingControl = ({ plotId }: PlantingControlProps) => {
  const [selectedCrop, setSelectedCrop] = useState('');
  
  // Получаем список доступных культур для посадки
  const { data: availableCrops, isLoading: isLoadingCrops } = useGetAvailableCropsQuery();
  // Получаем информацию о том, что уже посажено на этой грядке
  const { data: plantedCrops, isLoading: isLoadingPlanted } = useGetPlotCropsQuery(plotId);
  
  const [plantCrop, { isLoading: isPlanting }] = usePlantCropMutation();

  const handlePlant = async () => {
    if (!selectedCrop) {
      alert('Пожалуйста, выберите культуру для посадки.');
      return;
    }
    try {
      await plantCrop({ plotId, cropId: selectedCrop }).unwrap();
      alert('Культура успешно посажена!');
    } catch (err) {
      alert('Не удалось посадить культуру: ' + JSON.stringify(err));
    }
  };

  if (isLoadingCrops || isLoadingPlanted) return <p>Загрузка информации о посадках...</p>;

  // Если что-то уже растет, не показываем контрол
  if (plantedCrops && plantedCrops.length > 0) {
    return <p>Посажено: {plantedCrops[0].status}</p>; // Упрощенно, показываем статус первой посадки
  }

  return (
    <div>
      <select value={selectedCrop} onChange={(e) => setSelectedCrop(e.target.value)}>
        <option value="" disabled>Выберите культуру</option>
        {availableCrops?.map((crop) => (
          <option key={crop.id} value={crop.id}>
            {crop.name}
          </option>
        ))}
      </select>
      <button onClick={handlePlant} disabled={isPlanting || !selectedCrop}>
        {isPlanting ? 'Посадка...' : 'Посадить'}
      </button>
    </div>
  );
};
