import React, { useState } from 'react';
import {
    useGetRegionsQuery,
    useGetLandParcelsByRegionQuery,
    useCreateLandParcelMutation,
    useUpdateLandParcelMutation,
    useDeleteLandParcelMutation
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

// --- Типы ---
type LandParcel = { id: string; name: string; region_id: string };

// --- Форма создания ---
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
                    {regions?.map(region => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
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
    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: parcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(selectedRegion, { skip: !selectedRegion });
    
    const [deleteLandParcel] = useDeleteLandParcelMutation();
    const [updateLandParcel] = useUpdateLandParcelMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedParcel, setSelectedParcel] = useState<LandParcel | null>(null);
    const [editedName, setEditedName] = useState('');

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

    return (
        <Box>
            <CreateLandParcelForm />

            <Typography variant="h6" sx={{ mb: 2 }}>Список участков</Typography>

            <FormControl fullWidth margin="normal">
                <InputLabel>Фильтр по региону</InputLabel>
                <Select value={selectedRegion} label="Фильтр по региону" onChange={(e) => setSelectedRegion(e.target.value)} disabled={isLoadingRegions}>
                    <MenuItem value=""><em>Все регионы (не реализовано)</em></MenuItem>
                    {regions?.map(region => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
                </Select>
            </FormControl>

            {isLoadingParcels && <CircularProgress />}

            {selectedRegion && parcels && (
                <TableContainer component={Paper} sx={{mt: 2}}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>ID</TableCell>
                                <TableCell>Название</TableCell>
                                <TableCell align="right">Действия</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {parcels.map((parcel) => (
                                <TableRow key={parcel.id}>
                                    <TableCell>{parcel.id}</TableCell>
                                    <TableCell>{parcel.name}</TableCell>
                                    <TableCell align="right">
                                        <IconButton onClick={() => handleOpenEditModal(parcel as LandParcel)}><EditIcon /></IconButton>
                                        <IconButton onClick={() => handleOpenDeleteDialog(parcel as LandParcel)}><DeleteIcon /></IconButton>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}

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
