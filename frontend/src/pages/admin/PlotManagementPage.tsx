import React, { useState, useMemo } from 'react';
import type { Plot } from '../../features/api/apiSlice';
import {
    useGetRegionsForAdminQuery,
    useGetLandParcelsForAdminQuery,
    useGetStructuresForAdminQuery,
    useGetPlotsForAdminQuery,
    useCreatePlotMutation,
    useUpdatePlotMutation,
    useDeletePlotMutation,
    useRestorePlotMutation
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

// --- Форма создания ---
const CreatePlotForm = () => {
    const [name, setName] = useState('');
    const [size, setSize] = useState('');
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');
    const [structureId, setStructureId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsForAdminQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsForAdminQuery();
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresForAdminQuery();
    const [createPlot, { isLoading }] = useCreatePlotMutation();

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

    const filteredParcels = useMemo(() => landParcels?.filter(p => p.region_id === regionId && !p.deleted_at) || [], [landParcels, regionId]);
    const filteredStructures = useMemo(() => structures?.filter(s => s.land_parcel_id === landParcelId && !s.deleted_at) || [], [structures, landParcelId]);

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новую грядку</Typography>
            <FormControl fullWidth margin="normal" required><InputLabel>Регион</InputLabel><Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)} disabled={isLoadingRegions}>{regions?.filter(r => !r.deleted_at).map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!regionId || isLoadingParcels}><InputLabel>Участок</InputLabel><Select value={landParcelId} label="Участок" onChange={(e) => setLandParcelId(e.target.value)}>{filteredParcels.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!landParcelId || isLoadingStructures}><InputLabel>Сооружение</InputLabel><Select value={structureId} label="Сооружение" onChange={(e) => setStructureId(e.target.value)}>{filteredStructures.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select></FormControl>
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
    const [showDeleted, setShowDeleted] = useState(false);

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsForAdminQuery();
    const { data: parcels, isLoading: isLoadingParcels } = useGetLandParcelsForAdminQuery();
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresForAdminQuery();
    const { data: allPlots, isLoading: isLoadingPlots } = useGetPlotsForAdminQuery();

    const [deletePlot] = useDeletePlotMutation();
    const [updatePlot] = useUpdatePlotMutation();
    const [restorePlot] = useRestorePlotMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedPlot, setSelectedPlot] = useState<Plot | null>(null);
    const [editedName, setEditedName] = useState('');
    const [editedSize, setEditedSize] = useState('');

    const { structureMap, parcelMap, regionMap } = useMemo(() => {
        const sMap = new Map<string, { name: string, land_parcel_id: string }>();
        structures?.forEach(s => sMap.set(s.id, { name: s.name, land_parcel_id: s.land_parcel_id }));
        const pMap = new Map<string, { name: string, region_id: string }>();
        parcels?.forEach(p => pMap.set(p.id, { name: p.name, region_id: p.region_id }));
        const rMap = new Map<string, string>();
        regions?.forEach(r => rMap.set(r.id, r.name));
        return { structureMap: sMap, parcelMap: pMap, regionMap: rMap };
    }, [structures, parcels, regions]);

    const filteredPlots = useMemo(() => {
        let plots = allPlots || [];
        if (selectedStructure) {
            return plots.filter(p => p.structure_id === selectedStructure);
        }
        if (selectedParcel) {
            const structureIds = structures?.filter(s => s.land_parcel_id === selectedParcel).map(s => s.id) || [];
            return plots.filter(p => structureIds.includes(p.structure_id));
        }
        if (selectedRegion) {
            const parcelIds = parcels?.filter(p => p.region_id === selectedRegion).map(p => p.id) || [];
            const structureIds = structures?.filter(s => parcelIds.includes(s.land_parcel_id)).map(s => s.id) || [];
            return plots.filter(p => structureIds.includes(p.structure_id));
        }
        return plots;
    }, [allPlots, selectedRegion, selectedParcel, selectedStructure, parcels, structures]);

    const { activePlots, deletedPlots }: { activePlots: Plot[], deletedPlots: Plot[] } = useMemo(() => {
        const active: Plot[] = [];
        const deleted: Plot[] = [];
        filteredPlots?.forEach((plot: Plot) => {
            if (plot.deleted_at) {
                deleted.push(plot);
            } else {
                active.push(plot);
            }
        });
        return { activePlots: active, deletedPlots: deleted };
    }, [filteredPlots]);

    const plotsToDisplay = showDeleted ? deletedPlots : activePlots;

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

    const handleRestore = (id: string) => {
        restorePlot(id);
    };

    const handleUpdate = () => {
        if (selectedPlot && editedName.trim()) {
            updatePlot({ id: selectedPlot.id, name: editedName.trim(), size: editedSize });
            handleClose();
        }
    };

    if (isLoadingRegions || isLoadingParcels || isLoadingStructures || isLoadingPlots) return <CircularProgress />;

    return (
        <Box>
            <CreatePlotForm />
            <Typography variant="h6" sx={{ mb: 2, mt: 4 }}>Список грядок</Typography>
            <FormControl fullWidth margin="normal">
                <InputLabel>Фильтр по региону</InputLabel>
                <Select value={selectedRegion} label="Фильтр по региону" onChange={(e) => setSelectedRegion(e.target.value)}>{regions?.filter(r => !r.deleted_at).map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select>
            </FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedRegion}>
                <InputLabel>Фильтр по участку</InputLabel>
                <Select value={selectedParcel} label="Фильтр по участку" onChange={(e) => setSelectedParcel(e.target.value)}>{parcels?.filter(p => p.region_id === selectedRegion && !p.deleted_at).map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select>
            </FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedParcel}>
                <InputLabel>Фильтр по сооружению</InputLabel>
                <Select value={selectedStructure} label="Фильтр по сооружению" onChange={(e) => setSelectedStructure(e.target.value)}>{structures?.filter(s => s.land_parcel_id === selectedParcel && !s.deleted_at).map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select>
            </FormControl>

            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, mt: 2 }}>
                <Typography variant="h6">{showDeleted ? 'Удаленные грядки' : 'Активные грядки'}</Typography>
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
                            <TableCell>Размер</TableCell>
                            <TableCell>Статус</TableCell>
                            <TableCell>Сооружение</TableCell>
                            <TableCell>Участок</TableCell>
                            <TableCell>Регион</TableCell>
                            {showDeleted && <TableCell>Дата удаления</TableCell>}
                            <TableCell align="right">Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {plotsToDisplay.map((plot: Plot) => {
                            const structure = structureMap.get(plot.structure_id);
                            const parcel = structure ? parcelMap.get(structure.land_parcel_id) : undefined;
                            const region = parcel ? regionMap.get(parcel.region_id) : undefined;
                            return (
                                <TableRow key={plot.id} sx={{ backgroundColor: showDeleted ? '#f5f5f5' : 'inherit' }}>
                                    <TableCell>{plot.name}</TableCell>
                                    <TableCell>{plot.size}</TableCell>
                                    <TableCell>{plot.status}</TableCell>
                                    <TableCell>{structure?.name || 'N/A'}</TableCell>
                                    <TableCell>{parcel?.name || 'N/A'}</TableCell>
                                    <TableCell>{region || 'N/A'}</TableCell>
                                    {showDeleted && <TableCell>{plot.deleted_at ? new Date(plot.deleted_at).toLocaleString() : '-'}</TableCell>}
                                    <TableCell align="right">
                                        {showDeleted ? (
                                            <Tooltip title="Восстановить">
                                                <IconButton onClick={() => handleRestore(plot.id)}><RestoreFromTrashIcon /></IconButton>
                                            </Tooltip>
                                        ) : (
                                            <>
                                                <Tooltip title="Редактировать">
                                                    <IconButton onClick={() => handleOpenEditModal(plot)}><EditIcon /></IconButton>
                                                </Tooltip>
                                                <Tooltip title="Удалить">
                                                    <IconButton onClick={() => handleOpenDeleteDialog(plot)}><DeleteIcon /></IconButton>
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