
import React, { useState } from 'react';
import { useGetAvailableCropsQuery, usePlantCropMutation } from '../api/apiSlice';

interface PlantingControlProps {
  plotId: string;
}

const PlantingControl: React.FC<PlantingControlProps> = ({ plotId }) => {
  const { data: availableCrops, isLoading: isLoadingCrops } = useGetAvailableCropsQuery();
  const [plantCrop, { isLoading: isPlanting }] = usePlantCropMutation();
  const [selectedCrop, setSelectedCrop] = useState<string>('');

  const handlePlant = async () => {
    if (selectedCrop) {
      try {
        await plantCrop({ plotId, cropId: selectedCrop }).unwrap();
        // Optionally, show a success message or update UI
      } catch (error) {
        console.error('Failed to plant crop: ', error);
        // Optionally, show an error message
      }
    }
  };

  if (isLoadingCrops) {
    return <p>Loading crops...</p>;
  }

  return (
    <div>
      <select value={selectedCrop} onChange={(e) => setSelectedCrop(e.target.value)}>
        <option value="" disabled>Select a crop</option>
        {availableCrops?.map((crop) => (
          <option key={crop.id} value={crop.id}>
            {crop.name}
          </option>
        ))}
      </select>
      <button onClick={handlePlant} disabled={!selectedCrop || isPlanting}>
        {isPlanting ? 'Planting...' : 'Plant'}
      </button>
    </div>
  );
};

export default PlantingControl;
