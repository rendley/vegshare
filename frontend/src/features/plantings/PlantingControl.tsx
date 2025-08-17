
import React, { useState } from 'react';
import { useGetAvailableCropsQuery, useGetPlotCropsQuery, usePlantCropMutation, useRemoveCropMutation } from '../api/apiSlice';

interface PlantingControlProps {
  plotId: string;
}

const PlantingControl: React.FC<PlantingControlProps> = ({ plotId }) => {
  const { data: availableCrops, isLoading: isLoadingCrops } = useGetAvailableCropsQuery();
  const { data: plotCrops, isLoading: isLoadingPlotCrops } = useGetPlotCropsQuery(plotId);
  const [plantCrop, { isLoading: isPlanting }] = usePlantCropMutation();
  const [removeCrop, { isLoading: isRemoving }] = useRemoveCropMutation();
  const [selectedCrop, setSelectedCrop] = useState<string>('');

  const handlePlant = async () => {
    if (selectedCrop) {
      try {
        await plantCrop({ plotId, cropId: selectedCrop }).unwrap();
      } catch (error) {
        console.error('Failed to plant crop: ', error);
      }
    }
  };

  const handleRemove = async (plantingId: string) => {
    try {
      await removeCrop({ plotId, plantingId }).unwrap();
    } catch (error) {
      console.error('Failed to remove crop: ', error);
    }
  };

  if (isLoadingCrops || isLoadingPlotCrops) {
    return <p>Loading...</p>;
  }

  if (plotCrops && plotCrops.length > 0) {
    return (
      <div>
        <p>Planted crop: {plotCrops[0].crop_id}</p>
        <button onClick={() => handleRemove(plotCrops[0].id)} disabled={isRemoving}>
          {isRemoving ? 'Removing...' : 'Remove'}
        </button>
      </div>
    );
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
