import React from 'react';
import VideoPlayer from './VideoPlayer';
import PlantingControl from '../features/plantings/PlantingControl';
import { type EnrichedLease, type Camera } from '../features/api/apiSlice';
import {
    Box,
    Typography,
    Card,
    CardContent,
    CardHeader,
    Alert
} from '@mui/material';

interface PlotCardProps {
    lease: EnrichedLease;
}

const PlotCard: React.FC<PlotCardProps> = ({ lease }) => {
    if (!lease.plot) {
        return null; // Or some fallback UI
    }

    const { plot } = lease;

    return (
        <Card variant="outlined">
            <CardHeader
                title={`Грядка: ${plot.name}`}
                subheader={`Статус аренды: ${lease.status}`}
            />
            <CardContent>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                    ID: {plot.id}
                </Typography>

                {plot.cameras && plot.cameras.length > 0 ? (
                    plot.cameras.map((camera: Camera) => (
                        <VideoPlayer key={camera.id} camera={camera} />
                    ))
                ) : (
                    <Alert severity="info">Камер для этой грядки не найдено.</Alert>
                )}

                <Box sx={{ mt: 2, p: 2, border: '1px dashed grey', borderRadius: 1 }}>
                    <Typography variant="h6" gutterBottom>Управление посадкой</Typography>
                    <PlantingControl plotId={plot.id} />
                </Box>
            </CardContent>
        </Card>
    );
};

export default PlotCard;
