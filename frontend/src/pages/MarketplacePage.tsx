import { useState } from 'react';
import { useGetRegionsQuery, useGetLandParcelsByRegionQuery, useGetStructuresByLandParcelQuery, useGetPlotsByStructureQuery, useLeasePlotMutation } from '../features/api/apiSlice';
import { Box, Typography, FormControl, InputLabel, Select, MenuItem, Button, CircularProgress, Grid, Card, CardContent, CardActions } from '@mui/material';

export const MarketplacePage = () => {
    const [selectedRegion, setSelectedRegion] = useState('');
    const [selectedParcel, setSelectedParcel] = useState('');
    const [selectedStructure, setSelectedStructure] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: parcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(selectedRegion, { skip: !selectedRegion });
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresByLandParcelQuery(selectedParcel, { skip: !selectedParcel });
    const { data: plots, isLoading: isLoadingPlots } = useGetPlotsByStructureQuery(selectedStructure, { skip: !selectedStructure });

    const [leasePlot, { isLoading: isLeasing }] = useLeasePlotMutation();

    const handleLeaseClick = async (plotId: string) => {
        try {
            await leasePlot({ unit_id: plotId, unit_type: 'plot' }).unwrap();
            alert('Грядка успешно арендована!');
        } catch (err) {
            alert('Не удалось арендовать грядку. Ошибка: ' + JSON.stringify(err));
        }
    };

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h4" gutterBottom>Торговая площадка (Тест)</Typography>
            <Typography variant="body1" gutterBottom>Выберите участок для аренды.</Typography>

            <Grid container spacing={2} sx={{ mb: 3 }}>
                <Grid size={{ xs: 12, md: 4 }}>
                    <FormControl fullWidth disabled={isLoadingRegions}>
                        <InputLabel>Регион</InputLabel>
                        <Select value={selectedRegion} label="Регион" onChange={(e) => { setSelectedRegion(e.target.value); setSelectedParcel(''); setSelectedStructure(''); }}>
                            {regions?.map(region => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
                        </Select>
                    </FormControl>
                </Grid>
                <Grid size={{ xs: 12, md: 4 }}>
                    <FormControl fullWidth disabled={!selectedRegion || isLoadingParcels}>
                        <InputLabel>Участок</InputLabel>
                        <Select value={selectedParcel} label="Участок" onChange={(e) => { setSelectedParcel(e.target.value); setSelectedStructure(''); }}>
                            {parcels?.map(parcel => <MenuItem key={parcel.id} value={parcel.id}>{parcel.name}</MenuItem>)}
                        </Select>
                    </FormControl>
                </Grid>
                <Grid size={{ xs: 12, md: 4 }}>
                    <FormControl fullWidth disabled={!selectedParcel || isLoadingStructures}>
                        <InputLabel>Сооружение</InputLabel>
                        <Select value={selectedStructure} label="Сооружение" onChange={(e) => setSelectedStructure(e.target.value)}>
                            {structures?.map(structure => <MenuItem key={structure.id} value={structure.id}>{structure.name}</MenuItem>)}
                        </Select>
                    </FormControl>
                </Grid>
            </Grid>

            {isLoadingPlots && <CircularProgress />}

            <Grid container spacing={3}>
                {plots?.map(plot => (
                    <Grid key={plot.id} size={{ xs: 12, sm: 6, md: 4 }}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6">{plot.name}</Typography>
                                <Typography color="text.secondary">Размер: {plot.size || 'не указан'}</Typography>
                                <Typography color={plot.status === 'available' ? 'green' : 'red'}>Статус: {plot.status}</Typography>
                            </CardContent>
                            <CardActions>
                                {plot.status === 'available' && (
                                    <Button 
                                        size="small" 
                                        onClick={() => handleLeaseClick(plot.id)} 
                                        disabled={isLeasing}
                                    >
                                        {isLeasing ? <CircularProgress size={20} /> : 'Арендовать'}
                                    </Button>
                                )}
                            </CardActions>
                        </Card>
                    </Grid>
                ))}
            </Grid>
        </Box>
    );
};