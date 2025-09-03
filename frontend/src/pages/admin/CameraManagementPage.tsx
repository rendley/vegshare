import React, { useState, useEffect } from 'react';
import type { Camera } from '../../features/api/apiSlice';
import {
    useGetRegionsQuery,
    useGetLandParcelsByRegionQuery,
    useGetStructuresByLandParcelQuery,
    useGetPlotsByStructureQuery,
    useGetCamerasByPlotQuery,
    useCreateCameraMutation,
    useUpdateCameraMutation,
    useDeleteCameraMutation
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
    MenuItem
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

// --- Форма создания ---
const CreateCameraForm = () => {
    const [name, setName] = useState('');
    const [rtspPath, setRtspPath] = useState('');
    const [plotId, setPlotId] = useState('');
    const [structureId, setStructureId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');
    const [regionId, setRegionId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(regionId, { skip: !regionId });
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresByLandParcelQuery(landParcelId, { skip: !landParcelId });
    const { data: plots, isLoading: isLoadingPlots } = useGetPlotsByStructureQuery(structureId, { skip: !structureId });
    const [createCamera, { isLoading }] = useCreateCameraMutation();

    useEffect(() => { setLandParcelId(''); setStructureId(''); setPlotId(''); }, [regionId]);
    useEffect(() => { setStructureId(''); setPlotId(''); }, [landParcelId]);
    useEffect(() => { setPlotId(''); }, [structureId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && plotId && rtspPath.trim()) {
            createCamera({ plotId, name, rtsp_path_name: rtspPath }).unwrap().then(() => {
                setName('');
                setRtspPath('');
                setPlotId('');
                setStructureId('');
                setLandParcelId('');
                setRegionId('');
            });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новую камеру</Typography>
            <FormControl fullWidth margin="normal" required><InputLabel>Регион</InputLabel><Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)} disabled={isLoadingRegions}>{regions?.map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!regionId || isLoadingParcels}><InputLabel>Участок</InputLabel><Select value={landParcelId} label="Участок" onChange={(e) => setLandParcelId(e.target.value)}>{landParcels?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!landParcelId || isLoadingStructures}><InputLabel>Сооружение</InputLabel><Select value={structureId} label="Сооружение" onChange={(e) => setStructureId(e.target.value)}>{structures?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!structureId || isLoadingPlots}><InputLabel>Грядка</InputLabel><Select value={plotId} label="Грядка" onChange={(e) => setPlotId(e.target.value)}>{plots?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
            <TextField label="Название камеры" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <TextField label="Путь RTSP" value={rtspPath} onChange={(e) => setRtspPath(e.target.value)} fullWidth margin="normal" required />
            <Button type="submit" variant="contained" disabled={isLoading || !plotId}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
        </Box>
    );
}

// --- Таблица управления ---
const CameraManagementPage = () => {
    const [selectedRegion, setSelectedRegion] = useState('');
    const [selectedParcel, setSelectedParcel] = useState('');
    const [selectedStructure, setSelectedStructure] = useState('');
    const [selectedPlot, setSelectedPlot] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: parcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(selectedRegion, { skip: !selectedRegion });
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresByLandParcelQuery(selectedParcel, { skip: !selectedParcel });
    const { data: plots, isLoading: isLoadingPlots } = useGetPlotsByStructureQuery(selectedStructure, { skip: !selectedStructure });
    const { data: cameras, isLoading: isLoadingCameras } = useGetCamerasByPlotQuery(selectedPlot, { skip: !selectedPlot });

    const [deleteCamera] = useDeleteCameraMutation();
    const [updateCamera] = useUpdateCameraMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedCamera, setSelectedCamera] = useState<Camera | null>(null);
    const [editedName, setEditedName] = useState('');
    const [editedRtspPath, setEditedRtspPath] = useState('');

    useEffect(() => { setSelectedParcel(''); setSelectedStructure(''); setSelectedPlot(''); }, [selectedRegion]);
    useEffect(() => { setSelectedStructure(''); setSelectedPlot(''); }, [selectedParcel]);
    useEffect(() => { setSelectedPlot(''); }, [selectedStructure]);

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

    const handleUpdate = () => {
        if (selectedCamera && editedName.trim() && editedRtspPath.trim()) {
            updateCamera({ id: selectedCamera.id, name: editedName.trim(), rtsp_path_name: editedRtspPath.trim() });
            handleClose();
        }
    };

    return (
        <Box>
            <CreateCameraForm />
            <Typography variant="h6" sx={{ mb: 2 }}>Список камер</Typography>
            <FormControl fullWidth margin="normal"><InputLabel>Фильтр по региону</InputLabel><Select value={selectedRegion} label="Фильтр по региону" onChange={(e) => setSelectedRegion(e.target.value)} disabled={isLoadingRegions}>{regions?.map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedRegion || isLoadingParcels}><InputLabel>Фильтр по участку</InputLabel><Select value={selectedParcel} label="Фильтр по участку" onChange={(e) => setSelectedParcel(e.target.value)}>{parcels?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedParcel || isLoadingStructures}><InputLabel>Фильтр по сооружению</InputLabel><Select value={selectedStructure} label="Фильтр по сооружению" onChange={(e) => setSelectedStructure(e.target.value)}>{structures?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedStructure || isLoadingPlots}><InputLabel>Фильтр по грядке</InputLabel><Select value={selectedPlot} label="Фильтр по грядке" onChange={(e) => setSelectedPlot(e.target.value)}>{plots?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>

            {isLoadingCameras && <CircularProgress />}

            {selectedPlot && cameras && (
                <TableContainer component={Paper} sx={{mt: 2}}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>ID</TableCell>
                                <TableCell>Название</TableCell>
                                <TableCell>Путь RTSP</TableCell>
                                <TableCell align="right">Действия</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {cameras.map((camera) => (
                                <TableRow key={camera.id}>
                                    <TableCell>{camera.id}</TableCell>
                                    <TableCell>{camera.name}</TableCell>
                                    <TableCell>{camera.rtsp_path_name}</TableCell>
                                    <TableCell align="right">
                                        <IconButton onClick={() => handleOpenEditModal(camera)}><EditIcon /></IconButton>
                                        <IconButton onClick={() => handleOpenDeleteDialog(camera)}><DeleteIcon /></IconButton>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}

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