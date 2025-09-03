import React, { useState, useEffect } from 'react';
import type { Plot } from '../../features/api/apiSlice';
import {
    useGetRegionsQuery,
    useGetLandParcelsByRegionQuery,
    useGetStructuresByLandParcelQuery,
    useGetPlotsByStructureQuery,
    useCreatePlotMutation,
    useUpdatePlotMutation,
    useDeletePlotMutation
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
const CreatePlotForm = () => {
    const [name, setName] = useState('');
    const [size, setSize] = useState('');
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');
    const [structureId, setStructureId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(regionId, { skip: !regionId });
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresByLandParcelQuery(landParcelId, { skip: !landParcelId });
    const [createPlot, { isLoading }] = useCreatePlotMutation();

    useEffect(() => { setLandParcelId(''); setStructureId(''); }, [regionId]);
    useEffect(() => { setStructureId(''); }, [landParcelId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && structureId) {
            createPlot({ structure_id: structureId, name, size }).unwrap().then(() => {
                setName('');
                setSize('');
                setStructureId('');
                setLandParcelId('');
                setRegionId('');
            });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новую грядку</Typography>
            <FormControl fullWidth margin="normal" required><InputLabel>Регион</InputLabel><Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)} disabled={isLoadingRegions}>{regions?.map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!regionId || isLoadingParcels}><InputLabel>Участок</InputLabel><Select value={landParcelId} label="Участок" onChange={(e) => setLandParcelId(e.target.value)}>{landParcels?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!landParcelId || isLoadingStructures}><InputLabel>Сооружение</InputLabel><Select value={structureId} label="Сооружение" onChange={(e) => setStructureId(e.target.value)}>{structures?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select></FormControl>
            <TextField label="Название грядки" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <TextField label="Размер (опционально)" value={size} onChange={(e) => setSize(e.target.value)} fullWidth margin="normal" />
            <Button type="submit" variant="contained" disabled={isLoading || !structureId}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
        </Box>
    );
}

// --- Таблица управления ---
const PlotManagementPage = () => {
    const [selectedRegion, setSelectedRegion] = useState('');
    const [selectedParcel, setSelectedParcel] = useState('');
    const [selectedStructure, setSelectedStructure] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: parcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(selectedRegion, { skip: !selectedRegion });
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresByLandParcelQuery(selectedParcel, { skip: !selectedParcel });
    const { data: plots, isLoading: isLoadingPlots } = useGetPlotsByStructureQuery(selectedStructure, { skip: !selectedStructure });

    const [deletePlot] = useDeletePlotMutation();
    const [updatePlot] = useUpdatePlotMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedPlot, setSelectedPlot] = useState<Plot | null>(null);
    const [editedName, setEditedName] = useState('');
    const [editedSize, setEditedSize] = useState('');

    useEffect(() => { setSelectedParcel(''); setSelectedStructure(''); }, [selectedRegion]);
    useEffect(() => { setSelectedStructure(''); }, [selectedParcel]);

    const handleOpenEditModal = (plot: Plot) => {
        setSelectedPlot(plot);
        setEditedName(plot.name);
        setEditedSize(plot.size || '');
        setEditModalOpen(true);
    };

    const handleOpenDeleteDialog = (plot: Plot) => {
        setSelectedPlot(plot);
        setDeleteDialogOpen(true);
    };

    const handleClose = () => {
        setEditModalOpen(false);
        setDeleteDialogOpen(false);
        setSelectedPlot(null);
    };

    const handleDelete = () => {
        if (selectedPlot) {
            deletePlot(selectedPlot.id);
            handleClose();
        }
    };

    const handleUpdate = () => {
        if (selectedPlot && editedName.trim()) {
            updatePlot({ id: selectedPlot.id, name: editedName.trim(), size: editedSize });
            handleClose();
        }
    };

    return (
        <Box>
            <CreatePlotForm />
            <Typography variant="h6" sx={{ mb: 2 }}>Список грядок</Typography>
            <FormControl fullWidth margin="normal">
                <InputLabel>Фильтр по региону</InputLabel>
                <Select value={selectedRegion} label="Фильтр по региону" onChange={(e) => setSelectedRegion(e.target.value)} disabled={isLoadingRegions}>{regions?.map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select>
            </FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedRegion || isLoadingParcels}>
                <InputLabel>Фильтр по участку</InputLabel>
                <Select value={selectedParcel} label="Фильтр по участку" onChange={(e) => setSelectedParcel(e.target.value)}>{parcels?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select>
            </FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedParcel || isLoadingStructures}>
                <InputLabel>Фильтр по сооружению</InputLabel>
                <Select value={selectedStructure} label="Фильтр по сооружению" onChange={(e) => setSelectedStructure(e.target.value)}>{structures?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select>
            </FormControl>

            {isLoadingPlots && <CircularProgress />}

            {selectedStructure && plots && (
                <TableContainer component={Paper} sx={{mt: 2}}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>ID</TableCell>
                                <TableCell>Название</TableCell>
                                <TableCell>Размер</TableCell>
                                <TableCell>Статус</TableCell>
                                <TableCell align="right">Действия</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {plots.map((plot) => (
                                <TableRow key={plot.id}>
                                    <TableCell>{plot.id}</TableCell>
                                    <TableCell>{plot.name}</TableCell>
                                    <TableCell>{plot.size}</TableCell>
                                    <TableCell>{plot.status}</TableCell>
                                    <TableCell align="right">
                                        <IconButton onClick={() => handleOpenEditModal(plot)}><EditIcon /></IconButton>
                                        <IconButton onClick={() => handleOpenDeleteDialog(plot)}><DeleteIcon /></IconButton>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}

            <Modal open={editModalOpen} onClose={handleClose}>
                <Box sx={{ ...modalStyle }}>
                    <Typography variant="h6">Редактировать грядку</Typography>
                    <TextField label="Новое название" value={editedName} onChange={(e) => setEditedName(e.target.value)} fullWidth margin="normal" />
                    <TextField label="Новый размер" value={editedSize} onChange={(e) => setEditedSize(e.target.value)} fullWidth margin="normal" />
                    <Button onClick={handleUpdate} variant="contained">Сохранить</Button>
                </Box>
            </Modal>

            <Dialog open={deleteDialogOpen} onClose={handleClose}>
                <DialogTitle>Подтвердите удаление</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить грядку "{selectedPlot?.name}"?
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

export default PlotManagementPage;