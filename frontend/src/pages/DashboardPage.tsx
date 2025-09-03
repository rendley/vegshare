import { useGetMyLeasesQuery } from '../features/api/apiSlice';
import PlantingControl from '../features/plantings/PlantingControl';
import { Box, Typography, CircularProgress, Grid, Card, CardContent } from '@mui/material';

export const DashboardPage = () => {
    const { data: leases, isLoading, error } = useGetMyLeasesQuery();

    if (isLoading) {
        return <CircularProgress />;
    }

    if (error) {
        return <Typography color="error">Ошибка при загрузке данных.</Typography>;
    }

    const plotLeases = leases?.filter(lease => lease.unit_type === 'plot');

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h4" gutterBottom>Моя ферма</Typography>
            {plotLeases && plotLeases.length > 0 ? (
                <Grid container spacing={3}>
                    {plotLeases.map(lease => (
                        <Grid key={lease.id} size={{ xs: 12, md: 6 }}>
                            <Card>
                                <CardContent>
                                    <Typography variant="h6">Грядка ID: {lease.unit_id}</Typography>
                                    <Typography color="text.secondary">Статус аренды: {lease.status}</Typography>
                                    <hr />
                                    <PlantingControl plotId={lease.unit_id} />
                                </CardContent>
                            </Card>
                        </Grid>
                    ))}
                </Grid>
            ) : (
                <Typography>У вас пока нет арендованных участков.</Typography>
            )}
        </Box>
    );
};