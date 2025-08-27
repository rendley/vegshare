import React, { useState, useEffect } from 'react';
import type { Camera } from '../../features/api/apiSlice';
import {
    useGetRegionsQuery,
    useGetLandParcelsByRegionQuery,
    useGetGreenhousesByLandParcelQuery,
    useGetPlotsByGreenhouseQuery,
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
    const [greenhouseId, setGreenhouseId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');
    const [regionId, setRegionId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(regionId, { skip: !regionId });
    const { data: greenhouses, isLoading: isLoadingGreenhouses } = useGetGreenhousesByLandParcelQuery(landParcelId, { skip: !landParcelId });
    const { data: plots, isLoading: isLoadingPlots } = useGetPlotsByGreenhouseQuery(greenhouseId, { skip: !greenhouseId });
    const [createCamera, { isLoading }] = useCreateCameraMutation();

    useEffect(() => { setLandParcelId(''); setGreenhouseId(''); setPlotId(''); }, [regionId]);
    useEffect(() => { setGreenhouseId(''); setPlotId(''); }, [landParcelId]);
    useEffect(() => { setPlotId(''); }, [greenhouseId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && plotId && rtspPath.trim()) {
            createCamera({ plotId, name, rtsp_path_name: rtspPath }).unwrap().then(() => {
                setName('');
                setRtspPath('');
                setPlotId('');
                setGreenhouseId('');
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
            <FormControl fullWidth margin="normal" required disabled={!landParcelId || isLoadingGreenhouses}><InputLabel>Теплица</InputLabel><Select value={greenhouseId} label="Теплица" onChange={(e) => setGreenhouseId(e.target.value)}>{greenhouses?.map(g => <MenuItem key={g.id} value={g.id}>{g.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!greenhouseId || isLoadingPlots}><InputLabel>Грядка</InputLabel><Select value={plotId} label="Грядка" onChange={(e) => setPlotId(e.target.value)}>{plots?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
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
    const [selectedGreenhouse, setSelectedGreenhouse] = useState('');
    const [selectedPlot, setSelectedPlot] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: parcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(selectedRegion, { skip: !selectedRegion });
    const { data: greenhouses, isLoading: isLoadingGreenhouses } = useGetGreenhousesByLandParcelQuery(selectedParcel, { skip: !selectedParcel });
    const { data: plots, isLoading: isLoadingPlots } = useGetPlotsByGreenhouseQuery(selectedGreenhouse, { skip: !selectedGreenhouse });
    const { data: cameras, isLoading: isLoadingCameras } = useGetCamerasByPlotQuery(selectedPlot, { skip: !selectedPlot });

    const [deleteCamera] = useDeleteCameraMutation();
    const [updateCamera] = useUpdateCameraMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedCamera, setSelectedCamera] = useState<Camera | null>(null);
    const [editedName, setEditedName] = useState('');
    const [editedRtspPath, setEditedRtspPath] = useState('');

    useEffect(() => { setSelectedParcel(''); setSelectedGreenhouse(''); setSelectedPlot(''); }, [selectedRegion]);
    useEffect(() => { setSelectedGreenhouse(''); setSelectedPlot(''); }, [selectedParcel]);
    useEffect(() => { setSelectedPlot(''); }, [selectedGreenhouse]);

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
            <FormControl fullWidth margin="normal" disabled={!selectedParcel || isLoadingGreenhouses}><InputLabel>Фильтр по теплице</InputLabel><Select value={selectedGreenhouse} label="Фильтр по теплице" onChange={(e) => setSelectedGreenhouse(e.target.value)}>{greenhouses?.map(g => <MenuItem key={g.id} value={g.id}>{g.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedGreenhouse || isLoadingPlots}><InputLabel>Фильтр по грядке</InputLabel><Select value={selectedPlot} label="Фильтр по грядке" onChange={(e) => setSelectedPlot(e.target.value)}>{plots?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>

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