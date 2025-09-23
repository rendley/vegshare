import React, { useState, useEffect, useMemo } from 'react';
import type { Camera, Plot } from '../../features/api/apiSlice';
import {
    useGetRegionsForAdminQuery,
    useGetLandParcelsForAdminQuery,
    useGetStructuresForAdminQuery,
    useGetPlotsForAdminQuery,
    useGetCamerasForAdminQuery,
    useCreateCameraMutation,
    useUpdateCameraMutation,
    useDeleteCameraMutation,
    useRestoreCameraMutation
} from '../../features/api/apiSlice';
import {
    Box,
    Button,
    TextField,
    Typography,
    CircularProgress,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    IconButton,
    Modal,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Switch,
    FormControlLabel,
    Tooltip
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import RestoreFromTrashIcon from '@mui/icons-material/RestoreFromTrash';

const unitTypes = ['plot', 'coop', 'beehive']; // В будущем можно получать с бэкенда

// --- Форма создания ---
const CreateCameraForm = () => {
    const [name, setName] = useState('');
    const [rtspPath, setRtspPath] = useState('');
    const [unitId, setUnitId] = useState('');
    const [unitType, setUnitType] = useState('plot');
    
    // Состояния для каскадных селектов (только для plot)
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');
    const [structureId, setStructureId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsForAdminQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsForAdminQuery();
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresForAdminQuery();
    const { data: plots, isLoading: isLoadingPlots } = useGetPlotsForAdminQuery();
    const [createCamera, { isLoading }] = useCreateCameraMutation();

    useEffect(() => { setLandParcelId(''); setStructureId(''); setUnitId(''); }, [regionId]);
    useEffect(() => { setStructureId(''); setUnitId(''); }, [landParcelId]);
    useEffect(() => { setUnitId(''); }, [structureId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && unitId && rtspPath.trim()) {
            createCamera({ unit_id: unitId, unit_type: unitType, name, rtsp_path_name: rtspPath }).unwrap().then(() => {
                setName('');
                setRtspPath('');
                setUnitId('');
                setUnitType('plot');
                setStructureId('');
                setLandParcelId('');
                setRegionId('');
            });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новую камеру</Typography>
            <FormControl fullWidth margin="normal" required>
                <InputLabel>Тип Юнита</InputLabel>
                <Select value={unitType} label="Тип Юнита" onChange={(e) => setUnitType(e.target.value)}>
                    {unitTypes.map(type => <MenuItem key={type} value={type}>{type}</MenuItem>)}
                </Select>
            </FormControl>

            {unitType === 'plot' ? (
                <>
                    <FormControl fullWidth margin="normal" required><InputLabel>Регион</InputLabel><Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)} disabled={isLoadingRegions}>{regions?.filter(r => !r.deleted_at).map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select></FormControl>
                    <FormControl fullWidth margin="normal" required disabled={!regionId || isLoadingParcels}><InputLabel>Участок</InputLabel><Select value={landParcelId} label="Участок" onChange={(e) => setLandParcelId(e.target.value)}>{landParcels?.filter(p => p.region_id === regionId && !p.deleted_at).map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
                    <FormControl fullWidth margin="normal" required disabled={!landParcelId || isLoadingStructures}><InputLabel>Сооружение</InputLabel><Select value={structureId} label="Сооружение" onChange={(e) => setStructureId(e.target.value)}>{structures?.filter(s => s.land_parcel_id === landParcelId && !s.deleted_at).map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select></FormControl>
                    <FormControl fullWidth margin="normal" required disabled={!structureId || isLoadingPlots}><InputLabel>Грядка (Юнит)</InputLabel><Select value={unitId} label="Грядка (Юнит)" onChange={(e) => setUnitId(e.target.value)}>{plots?.filter(p => p.structure_id === structureId && !p.deleted_at).map((p: Plot) => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
                </>
            ) : (
                <TextField label="Unit ID" value={unitId} onChange={(e) => setUnitId(e.target.value)} fullWidth margin="normal" required />
            )}

            <TextField label="Название камеры" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <TextField label="Путь RTSP" value={rtspPath} onChange={(e) => setRtspPath(e.target.value)} fullWidth margin="normal" required />
            <Button type="submit" variant="contained" disabled={isLoading || !unitId}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
        </Box>
    );
}

// --- Таблица управления ---
const CameraManagementPage = () => {
    const [showDeleted, setShowDeleted] = useState(false);
    const { data: allCameras, isLoading: isLoadingCameras } = useGetCamerasForAdminQuery();
    const { data: allPlots, isLoading: isLoadingPlots } = useGetPlotsForAdminQuery();
    const { data: allStructures, isLoading: isLoadingStructures } = useGetStructuresForAdminQuery();
    const { data: allParcels, isLoading: isLoadingParcels } = useGetLandParcelsForAdminQuery();
    const { data: allRegions, isLoading: isLoadingRegions } = useGetRegionsForAdminQuery();

    const [deleteCamera] = useDeleteCameraMutation();
    const [updateCamera] = useUpdateCameraMutation();
    const [restoreCamera] = useRestoreCameraMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedCamera, setSelectedCamera] = useState<Camera | null>(null);
    const [editedName, setEditedName] = useState('');
    const [editedRtspPath, setEditedRtspPath] = useState('');

    const { plotMap, structureMap, parcelMap, regionMap } = useMemo(() => {
        const plotMap = new Map<string, { name: string, structure_id: string }>();
        allPlots?.forEach(p => plotMap.set(p.id, { name: p.name, structure_id: p.structure_id }));
        const structureMap = new Map<string, { name: string, land_parcel_id: string }>();
        allStructures?.forEach(s => structureMap.set(s.id, { name: s.name, land_parcel_id: s.land_parcel_id }));
        const pMap = new Map<string, { name: string, region_id: string }>();
        allParcels?.forEach(p => pMap.set(p.id, { name: p.name, region_id: p.region_id }));
        const rMap = new Map<string, string>();
        allRegions?.forEach(r => rMap.set(r.id, r.name));
        return { plotMap, structureMap, parcelMap: pMap, regionMap: rMap };
    }, [allPlots, allStructures, allParcels, allRegions]);

    const { activeCameras, deletedCameras } = useMemo(() => {
        const active: Camera[] = [];
        const deleted: Camera[] = [];
        allCameras?.forEach((camera: Camera) => {
            if (camera.deleted_at) {
                deleted.push(camera);
            } else {
                active.push(camera);
            }
        });
        return { activeCameras: active, deletedCameras: deleted };
    }, [allCameras]) as { activeCameras: Camera[]; deletedCameras: Camera[] };

    const camerasToDisplay = showDeleted ? deletedCameras : activeCameras;

    const handleOpenEditModal = (camera: Camera) => {
        setSelectedCamera(camera);
        setEditedName(camera.name);
        setEditedRtspPath(camera.rtsp_path_name);
        setEditModalOpen(true);
    };

    const handleOpenDeleteDialog = (camera: Camera) => {
        setSelectedCamera(camera);
        setDeleteDialogOpen(true);
    };

    const handleClose = () => {
        setEditModalOpen(false);
        setDeleteDialogOpen(false);
        setSelectedCamera(null);
    };

    const handleDelete = () => {
        if (selectedCamera) {
            deleteCamera(selectedCamera.id);
            handleClose();
        }
    };

    const handleRestore = (id: string) => {
        restoreCamera(id);
    };

    const handleUpdate = () => {
        if (selectedCamera && editedName.trim() && editedRtspPath.trim()) {
            updateCamera({ id: selectedCamera.id, name: editedName.trim(), rtsp_path_name: editedRtspPath.trim() });
            handleClose();
        }
    };

    if (isLoadingCameras || isLoadingPlots || isLoadingStructures || isLoadingParcels || isLoadingRegions) return <CircularProgress />;

    return (
        <Box>
            <CreateCameraForm />
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, mt: 4 }}>
                <Typography variant="h6">{showDeleted ? 'Удаленные камеры' : 'Активные камеры'}</Typography>
                <FormControlLabel
                    control={<Switch checked={showDeleted} onChange={(e) => setShowDeleted(e.target.checked)} />}
                    label="Показать удаленные"
                />
            </Box>

            <TableContainer component={Paper} sx={{mt: 2}}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>Название</TableCell>
                            <TableCell>Путь RTSP</TableCell>
                            <TableCell>Юнит</TableCell>
                            <TableCell>Тип юнита</TableCell>
                            <TableCell>Сооружение</TableCell>
                            <TableCell>Участок</TableCell>
                            <TableCell>Регион</TableCell>
                            {showDeleted && <TableCell>Дата удаления</TableCell>}
                            <TableCell align="right">Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {camerasToDisplay.map((camera: Camera) => {
                            const plot = plotMap.get(camera.unit_id);
                            const structure = plot ? structureMap.get(plot.structure_id) : undefined;
                            const parcel = structure ? parcelMap.get(structure.land_parcel_id) : undefined;
                            const region = parcel ? regionMap.get(parcel.region_id) : undefined;
                            return (
                                <TableRow key={camera.id} sx={{ backgroundColor: showDeleted ? '#f5f5f5' : 'inherit' }}>
                                    <TableCell>{camera.name}</TableCell>
                                    <TableCell>{camera.rtsp_path_name}</TableCell>
                                    <TableCell>{plot?.name || camera.unit_id}</TableCell>
                                    <TableCell>{camera.unit_type}</TableCell>
                                    <TableCell>{structure?.name || 'N/A'}</TableCell>
                                    <TableCell>{parcel?.name || 'N/A'}</TableCell>
                                    <TableCell>{region || 'N/A'}</TableCell>
                                    {showDeleted && <TableCell>{camera.deleted_at ? new Date(camera.deleted_at).toLocaleString() : '-'}</TableCell>}
                                    <TableCell align="right">
                                        {showDeleted ? (
                                            <Tooltip title="Восстановить">
                                                <IconButton onClick={() => handleRestore(camera.id)}><RestoreFromTrashIcon /></IconButton>
                                            </Tooltip>
                                        ) : (
                                            <>
                                                <Tooltip title="Редактировать">
                                                    <IconButton onClick={() => handleOpenEditModal(camera)}><EditIcon /></IconButton>
                                                </Tooltip>
                                                <Tooltip title="Удалить">
                                                    <IconButton onClick={() => handleOpenDeleteDialog(camera)}><DeleteIcon /></IconButton>
                                                </Tooltip>
                                            </>
                                        )}
                                    </TableCell>
                                </TableRow>
                            );
                        })}
                    </TableBody>
                </Table>
            </TableContainer>

            <Modal open={editModalOpen} onClose={handleClose}>
                <Box sx={{ ...modalStyle }}>
                    <Typography variant="h6">Редактировать камеру</Typography>
                    <TextField label="Новое название" value={editedName} onChange={(e) => setEditedName(e.target.value)} fullWidth margin="normal" />
                    <TextField label="Новый путь RTSP" value={editedRtspPath} onChange={(e) => setEditedRtspPath(e.target.value)} fullWidth margin="normal" />
                    <Button onClick={handleUpdate} variant="contained">Сохранить</Button>
                </Box>
            </Modal>

            <Dialog open={deleteDialogOpen} onClose={handleClose}>
                <DialogTitle>Подтвердите удаление</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить камеру "{selectedCamera?.name}"?
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>Отмена</Button>
                    <Button onClick={handleDelete} color="error">Удалить</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};

const modalStyle = {
    position: 'absolute' as 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 400,
    bgcolor: 'background.paper',
    border: '2px solid #000',
    boxShadow: 24,
    p: 4,
};

export default CameraManagementPage;