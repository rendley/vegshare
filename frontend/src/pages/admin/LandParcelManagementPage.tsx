import React, { useState, useMemo } from 'react';
import {
    useGetRegionsQuery,
    useGetLandParcelsForAdminQuery, // Используем новый хук
    useCreateLandParcelMutation,
    useUpdateLandParcelMutation,
    useDeleteLandParcelMutation,
    useRestoreLandParcelMutation, // Используем новый хук
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
    Tooltip,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import RestoreFromTrashIcon from '@mui/icons-material/RestoreFromTrash';

// --- Типы ---
type LandParcel = { id: string; name: string; region_id: string; deleted_at?: string };

// --- Форма создания (без изменений) ---
const CreateLandParcelForm = () => {
    const [name, setName] = useState('');
    const [regionId, setRegionId] = useState('');
    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const [createLandParcel, { isLoading }] = useCreateLandParcelMutation();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && regionId) {
            createLandParcel({ regionId, name }).unwrap().then(() => {
                setName('');
                setRegionId('');
            });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новый участок</Typography>
            <FormControl fullWidth margin="normal" required>
                <InputLabel>Регион</InputLabel>
                <Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)} disabled={isLoadingRegions}>
                    {regions?.filter(r => !r.deleted_at).map(region => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
                </Select>
            </FormControl>
            <TextField label="Название участка" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <Button type="submit" variant="contained" disabled={isLoading || !regionId}>
                {isLoading ? <CircularProgress size={24} /> : 'Создать'}
            </Button>
        </Box>
    );
}

// --- Таблица управления ---
const LandParcelManagementPage = () => {
    const [selectedRegion, setSelectedRegion] = useState('');
    const [showDeleted, setShowDeleted] = useState(false);

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: allParcels, isLoading: isLoadingParcels, isError } = useGetLandParcelsForAdminQuery();
    
    const [deleteLandParcel] = useDeleteLandParcelMutation();
    const [updateLandParcel] = useUpdateLandParcelMutation();
    const [restoreLandParcel] = useRestoreLandParcelMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedParcel, setSelectedParcel] = useState<LandParcel | null>(null);
    const [editedName, setEditedName] = useState('');

    const regionMap = useMemo(() => {
        const map = new Map<string, string>();
        regions?.forEach(region => {
            map.set(region.id, region.name);
        });
        return map;
    }, [regions]);

    const { activeParcels, deletedParcels } = useMemo(() => {
        const filteredByRegion = allParcels?.filter(p => !selectedRegion || p.region_id === selectedRegion) || [];
        const active: LandParcel[] = [];
        const deleted: LandParcel[] = [];
        filteredByRegion.forEach(parcel => {
            if (parcel.deleted_at) {
                deleted.push(parcel);
            } else {
                active.push(parcel);
            }
        });
        return { activeParcels: active, deletedParcels: deleted };
    }, [allParcels, selectedRegion]);

    const parcelsToDisplay = showDeleted ? deletedParcels : activeParcels;

    const handleOpenEditModal = (parcel: LandParcel) => {
        setSelectedParcel(parcel);
        setEditedName(parcel.name);
        setEditModalOpen(true);
    };

    const handleOpenDeleteDialog = (parcel: LandParcel) => {
        setSelectedParcel(parcel);
        setDeleteDialogOpen(true);
    };

    const handleClose = () => {
        setEditModalOpen(false);
        setDeleteDialogOpen(false);
        setSelectedParcel(null);
    };

    const handleDelete = () => {
        if (selectedParcel) {
            deleteLandParcel(selectedParcel.id);
            handleClose();
        }
    };

    const handleUpdate = () => {
        if (selectedParcel && editedName.trim()) {
            updateLandParcel({ id: selectedParcel.id, name: editedName.trim() });
            handleClose();
        }
    };

    const handleRestore = (parcelId: string) => {
        restoreLandParcel(parcelId);
    };

    if (isLoadingParcels || isLoadingRegions) return <CircularProgress />;
    if (isError) return <Typography color="error">Не удалось загрузить данные.</Typography>;

    return (
        <Box>
            <CreateLandParcelForm />

            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, mt: 4 }}>
                <Typography variant="h6">{showDeleted ? 'Удаленные участки' : 'Активные участки'}</Typography>
                <FormControlLabel
                    control={<Switch checked={showDeleted} onChange={(e) => setShowDeleted(e.target.checked)} />}
                    label="Показать удаленные"
                />
            </Box>

            <FormControl fullWidth margin="normal">
                <InputLabel>Фильтр по региону</InputLabel>
                <Select value={selectedRegion} label="Фильтр по региону" onChange={(e) => setSelectedRegion(e.target.value)}>
                    <MenuItem value=""><em>Все регионы</em></MenuItem>
                    {regions?.filter(r => !r.deleted_at).map(region => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
                </Select>
            </FormControl>

            <TableContainer component={Paper} sx={{mt: 2}}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>ID</TableCell>
                            <TableCell>Название</TableCell>
                            <TableCell>Регион</TableCell>
                            {showDeleted && <TableCell>Дата удаления</TableCell>}
                            <TableCell align="right">Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {parcelsToDisplay.map((parcel) => (
                            <TableRow key={parcel.id} sx={{ backgroundColor: showDeleted ? '#f5f5f5' : 'inherit' }}>
                                <TableCell>{parcel.id}</TableCell>
                                <TableCell>{parcel.name}</TableCell>
                                <TableCell>{regionMap.get(parcel.region_id) || 'N/A'}</TableCell>
                                {showDeleted && <TableCell>{new Date(parcel.deleted_at!).toLocaleString()}</TableCell>}
                                <TableCell align="right">
                                    {showDeleted ? (
                                        <Tooltip title="Восстановить">
                                            <IconButton onClick={() => handleRestore(parcel.id)}><RestoreFromTrashIcon /></IconButton>
                                        </Tooltip>
                                    ) : (
                                        <>
                                            <Tooltip title="Редактировать">
                                                <IconButton onClick={() => handleOpenEditModal(parcel)}><EditIcon /></IconButton>
                                            </Tooltip>
                                            <Tooltip title="Удалить">
                                                <IconButton onClick={() => handleOpenDeleteDialog(parcel)}><DeleteIcon /></IconButton>
                                            </Tooltip>
                                        </>
                                    )}
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>

            {/* Модальное окно редактирования */}
            <Modal open={editModalOpen} onClose={handleClose}>
                <Box sx={{ ...modalStyle }}>
                    <Typography variant="h6">Редактировать участок</Typography>
                    <TextField label="Новое название" value={editedName} onChange={(e) => setEditedName(e.target.value)} fullWidth margin="normal" />
                    <Button onClick={handleUpdate} variant="contained">Сохранить</Button>
                </Box>
            </Modal>

            {/* Диалог подтверждения удаления */}
            <Dialog open={deleteDialogOpen} onClose={handleClose}>
                <DialogTitle>Подтвердите удаление</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить участок "{selectedParcel?.name}"?
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

export default LandParcelManagementPage;