import React, { useState } from 'react';
import { useGetCatalogItemsQuery, useCreateActionMutation } from '../api/apiSlice';
import {
    Box,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Button,
    CircularProgress,
    TextField,
    ButtonGroup
} from '@mui/material';

interface PlantingControlProps {
    plotId: string;
}

const PlantingControl: React.FC<PlantingControlProps> = ({ plotId }) => {
    const { data: availableCrops, isLoading: isLoadingCrops } = useGetCatalogItemsQuery('crop');
    const [createAction, { isLoading: isCreatingAction }] = useCreateActionMutation();

    const [selectedCrop, setSelectedCrop] = useState<string>('');
    const [quantity, setQuantity] = useState(1);

    const handleCreateAction = async (actionType: string, parameters: any) => {
        try {
            await createAction({
                unit_id: plotId,
                unit_type: 'plot',
                action_type: actionType,
                parameters: parameters,
            }).unwrap();
        } catch (error) {
            console.error(`Failed to create ${actionType} action`, error);
        }
    };

    const handlePlant = () => {
        if (selectedCrop && quantity > 0) {
            handleCreateAction('plant', { item_id: selectedCrop, quantity });
        }
    };

    const handleWater = () => {
        handleCreateAction('water', { volume_liters: 5 });
    };

    if (isLoadingCrops) {
        return <CircularProgress size={24} />;
    }

    return (
        <Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                <FormControl fullWidth size="small">
                    <InputLabel>Выбрать культуру</InputLabel>
                    <Select
                        value={selectedCrop}
                        label="Выбрать культуру"
                        onChange={(e) => setSelectedCrop(e.target.value)}
                    >
                        {availableCrops?.map((crop) => (
                            <MenuItem key={crop.id} value={crop.id}>
                                {crop.name}
                            </MenuItem>
                        ))}
                    </Select>
                </FormControl>
                <TextField
                    type="number"
                    label="Кол-во"
                    size="small"
                    value={quantity}
                    onChange={(e) => setQuantity(parseInt(e.target.value, 10) || 1)}
                    sx={{ width: 150 }}
                />
                <Button
                    variant="contained"
                    onClick={handlePlant}
                    disabled={!selectedCrop || isCreatingAction}
                    sx={{ flexShrink: 0 }}
                >
                    Посадить
                </Button>
            </Box>
            <ButtonGroup variant="outlined" size="small">
                 <Button onClick={handleWater} disabled={isCreatingAction}>
                    Полить
                </Button>
                {/* Другие общие действия можно добавить сюда */}
            </ButtonGroup>
        </Box>
    );
};

export default PlantingControl;
