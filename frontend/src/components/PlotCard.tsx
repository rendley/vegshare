import React, { useEffect, useRef } from 'react';
import { useDispatch } from 'react-redux';
import VideoPlayer from './VideoPlayer';
import PlantingDisplay from '../features/plantings/PlantingDisplay';
import PlantingControl from '../features/plantings/PlantingControl';
import { useGetActionsForUnitQuery } from '../features/api/apiSlice';
import { type EnrichedLease, type Camera, type OperationLog } from '../features/api/apiSlice';
import { showNotification } from '../store/notificationSlice';
import {
    Box,
    Typography,
    Card,
    CardContent,
    CardHeader,
    Alert,
    Divider,
    List,
    ListItem,
    ListItemText,
    Chip,
    CircularProgress
} from '@mui/material';

interface PlotCardProps {
    lease: EnrichedLease;
}

const OperationsHistory: React.FC<{ plotId: string }> = ({ plotId }) => {
    const dispatch = useDispatch();
    const { data: actions, isLoading } = useGetActionsForUnitQuery(plotId, {
        pollingInterval: 5000, // Опрашивать каждые 5 секунд
    });

    const prevActionsRef = useRef<OperationLog[]>([]);

    useEffect(() => {
        if (prevActionsRef.current && actions) {
            // Ищем завершенные или проваленные задачи
            prevActionsRef.current.forEach(prevAction => {
                const newAction = actions.find(a => a.id === prevAction.id);
                if (newAction && (prevAction.status === 'in_progress' || prevAction.status === 'processing')) {
                    if (newAction.status === 'completed') {
                        dispatch(showNotification({ message: `Действие '${newAction.action_type}' успешно завершено!`, severity: 'success' }));
                    }
                    if (newAction.status === 'failed') {
                        dispatch(showNotification({ message: `Ошибка выполнения действия '${newAction.action_type}'`, severity: 'error' }));
                    }
                }
            });
        }
        prevActionsRef.current = actions || [];
    }, [actions, dispatch]);

    if (isLoading) return <CircularProgress size={20} />;

    const activeActions = actions?.filter(a => a.status === 'pending' || a.status === 'processing' || a.status === 'in_progress') || [];
    const completedActions = actions?.filter(a => a.status === 'completed' || a.status === 'failed').slice(0, 3) || [];

    return (
        <Box>
            <Typography variant="subtitle2" color="text.secondary" gutterBottom>Активные задачи:</Typography>
            {activeActions.length > 0 ? (
                <List dense>
                    {activeActions.map((action: OperationLog) => (
                        <ListItem key={action.id} disableGutters secondaryAction={<Chip label={action.status} size="small" color="primary" />}>
                            <ListItemText primary={`Действие: ${action.action_type}`} />
                        </ListItem>
                    ))}
                </List>
            ) : (
                <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>Нет активных задач.</Typography>
            )}

            <Divider sx={{ my: 1 }} />

            <Typography variant="subtitle2" color="text.secondary" gutterBottom>История операций:</Typography>
            {completedActions.length > 0 ? (
                <List dense>
                    {completedActions.map((action: OperationLog) => (
                        <ListItem key={action.id} disableGutters secondaryAction={<Chip label={action.status} size="small" color={action.status === 'completed' ? 'success' : 'error'} />}>
                            <ListItemText primary={`Действие: ${action.action_type}`} secondary={new Date(action.updated_at).toLocaleString()} />
                        </ListItem>
                    ))}
                </List>
            ) : (
                <Typography variant="body2" color="text.secondary">История операций пуста.</Typography>
            )}
        </Box>
    );
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
                    <Typography variant="h6" gutterBottom>Содержимое грядки</Typography>
                    <PlantingDisplay contents={plot.contents || []} />
                    <Divider sx={{ my: 2 }} />
                    <Typography variant="h6" gutterBottom>Операции</Typography>
                    <OperationsHistory plotId={plot.id} />
                    <Divider sx={{ my: 2 }} />
                    <Typography variant="h6" gutterBottom>Управление</Typography>
                    <PlantingControl plotId={plot.id} />
                </Box>
            </CardContent>
        </Card>
    );
};

export default PlotCard;