import React, { useState } from 'react';
import { type CatalogItem, useGetCatalogItemsQuery, useGetActionsForUnitQuery, useCreateActionMutation, useCancelActionMutation } from '../api/apiSlice';

interface PlantingControlProps {
  plotId: string;
}

const PlantingControl: React.FC<PlantingControlProps> = ({ plotId }) => {
  const { data: availableCrops, isLoading: isLoadingCrops } = useGetCatalogItemsQuery('crop');
  const { data: actions, isLoading: isLoadingActions } = useGetActionsForUnitQuery(plotId);
  const [createAction, { isLoading: isCreatingAction }] = useCreateActionMutation();
  const [cancelAction, { isLoading: isCancellingAction }] = useCancelActionMutation();
  const [selectedCrop, setSelectedCrop] = useState<string>('');

  const handlePlant = async () => {
    if (selectedCrop) {
      try {
        await createAction({
          unit_id: plotId,
          unit_type: 'plot',
          action_type: 'plant',
          parameters: { crop_id: selectedCrop },
        }).unwrap();
      } catch (error) {
        console.error('Failed to plant crop: ', error);
      }
    }
  };

  const handleCancel = async (actionId: string) => {
    try {
      await cancelAction(actionId).unwrap();
    } catch (error) {
      console.error('Failed to cancel action: ', error);
    }
  };

  const handleWater = async () => {
    try {
      await createAction({
        unit_id: plotId,
        unit_type: 'plot',
        action_type: 'water',
        parameters: { volume_liters: 5 },
      }).unwrap();
    } catch (error) {
      console.error('Failed to water crop: ', error);
    }
  };

  if (isLoadingCrops || isLoadingActions) {
    return <p>Loading...</p>;
  }

  const plantAction = actions?.find(a => a.action_type === 'plant' && a.status === 'completed');

  if (plantAction) {
    const plantedCrop = availableCrops?.find((c: CatalogItem) => c.id === plantAction.parameters.crop_id);
    return (
      <div>
        <p>Planted: {plantedCrop ? plantedCrop.name : 'Unknown Crop'}</p>
        <p>Harvest in: {plantedCrop?.attributes.harvest_time_days || 'N/A'} days</p>
        <button onClick={() => handleCancel(plantAction.id)} disabled={isCancellingAction}>
          {isCancellingAction ? 'Removing...' : 'Remove'}
        </button>
        <button onClick={handleWater} disabled={isCreatingAction}>
          {isCreatingAction ? 'Watering...' : 'Water'}
        </button>
      </div>
    );
  }

  return (
    <div>
      <select value={selectedCrop} onChange={(e) => setSelectedCrop(e.target.value)}>
        <option value="" disabled>Select a crop</option>
        {availableCrops?.map((crop: CatalogItem) => (
          <option key={crop.id} value={crop.id}>
            {crop.name}
          </option>
        ))}
      </select>
      <button onClick={handlePlant} disabled={!selectedCrop || isCreatingAction}>
        {isCreatingAction ? 'Planting...' : 'Plant'}
      </button>
    </div>
  );
};

export default PlantingControl;
