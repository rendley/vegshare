import { useGetMyLeasesQuery, type EnrichedLease } from '../features/api/apiSlice';
import PlotCard from '../components/PlotCard';
import { Box, Typography, CircularProgress, Grid, Alert } from '@mui/material';

export const DashboardPage = () => {
    const { data: leases, isLoading, error } = useGetMyLeasesQuery();

    if (isLoading) {
        return <CircularProgress />;
    }

    if (error) {
        return <Alert severity="error">Ошибка при загрузке данных.</Alert>;
    }

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h4" gutterBottom>Моя ферма</Typography>
            {leases && leases.length > 0 ? (
                <Grid container spacing={3}>
                    {leases.map((lease: EnrichedLease) => (
                        <Grid size={{ xs: 12, md: 6 }} key={lease.id}>
                            {lease.unit_type === 'plot' && lease.plot && (
                                <PlotCard lease={lease} />
                            )}
                            {/* 
                                В будущем здесь будут другие типы юнитов:
                                {lease.unit_type === 'coop' && lease.coop && (
                                    <CoopCard lease={lease} />
                                )}
                            */}
                        </Grid>
                    ))}
                </Grid>
            ) : (
                <Typography>У вас пока нет арендованных участков.</Typography>
            )}
        </Box>
    );
};