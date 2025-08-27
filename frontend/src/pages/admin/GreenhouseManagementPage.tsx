import React, { useState, useEffect } from 'react';
import type { Greenhouse } from '../../features/api/apiSlice';
import {
    useGetRegionsQuery,
    useGetLandParcelsByRegionQuery,
    useGetGreenhousesByLandParcelQuery,
    useCreateGreenhouseMutation,
    useUpdateGreenhouseMutation,
    useDeleteGreenhouseMutation
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
const CreateGreenhouseForm = () => {
    const [name, setName] = useState('');
    const [type, setType] = useState('');
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(regionId, { skip: !regionId });
    const [createGreenhouse, { isLoading }] = useCreateGreenhouseMutation();

    useEffect(() => {
        setLandParcelId('');
    }, [regionId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && landParcelId) {
            createGreenhouse({ landParcelId, name, type }).unwrap().then(() => {
                setName('');
                setType('');
                setLandParcelId('');
                setRegionId('');
            });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новую теплицу</Typography>
            <FormControl fullWidth margin="normal" required>
                <InputLabel>Регион</InputLabel>
                <Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)} disabled={isLoadingRegions}>
                    {regions?.map(region => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
                </Select>
            </FormControl>
            <FormControl fullWidth margin="normal" required disabled={!regionId || isLoadingParcels}>
                <InputLabel>Участок</InputLabel>
                <Select value={landParcelId} label="Участок" onChange={(e) => setLandParcelId(e.target.value)}>
                    {landParcels?.map(parcel => <MenuItem key={parcel.id} value={parcel.id}>{parcel.name}</MenuItem>)}
                </Select>
            </FormControl>
            <TextField label="Название теплицы" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <TextField label="Тип (опционально)" value={type} onChange={(e) => setType(e.target.value)} fullWidth margin="normal" />
            <Button type="submit" variant="contained" disabled={isLoading || !landParcelId}>
                {isLoading ? <CircularProgress size={24} /> : 'Создать'}
            </Button>
        </Box>
    );
}

// --- Таблица управления ---
const GreenhouseManagementPage = () => {
    const [selectedRegion, setSelectedRegion] = useState('');
    const [selectedParcel, setSelectedParcel] = useState('');
    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: parcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(selectedRegion, { skip: !selectedRegion });
    const { data: greenhouses, isLoading: isLoadingGreenhouses } = useGetGreenhousesByLandParcelQuery(selectedParcel, { skip: !selectedParcel });

    const [deleteGreenhouse] = useDeleteGreenhouseMutation();
    const [updateGreenhouse] = useUpdateGreenhouseMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedGreenhouse, setSelectedGreenhouse] = useState<Greenhouse | null>(null);
    const [editedName, setEditedName] = useState('');
    const [editedType, setEditedType] = useState('');

    useEffect(() => { setSelectedParcel(''); }, [selectedRegion]);

    const handleOpenEditModal = (greenhouse: Greenhouse) => {
        setSelectedGreenhouse(greenhouse);
        setEditedName(greenhouse.name);
        setEditedType(greenhouse.type || '');
        setEditModalOpen(true);
    };

    const handleOpenDeleteDialog = (greenhouse: Greenhouse) => {
        setSelectedGreenhouse(greenhouse);
        setDeleteDialogOpen(true);
    };

    const handleClose = () => {
        setEditModalOpen(false);
        setDeleteDialogOpen(false);
        setSelectedGreenhouse(null);
    };

    const handleDelete = () => {
        if (selectedGreenhouse) {
            deleteGreenhouse(selectedGreenhouse.id);
            handleClose();
        }
    };

    const handleUpdate = () => {
        if (selectedGreenhouse && editedName.trim()) {
            updateGreenhouse({ id: selectedGreenhouse.id, name: editedName.trim(), type: editedType });
            handleClose();
        }
    };

    return (
        <Box>
            <CreateGreenhouseForm />
            <Typography variant="h6" sx={{ mb: 2 }}>Список теплиц</Typography>
            <FormControl fullWidth margin="normal">
                <InputLabel>Фильтр по региону</InputLabel>
                <Select value={selectedRegion} label="Фильтр по региону" onChange={(e) => setSelectedRegion(e.target.value)} disabled={isLoadingRegions}>
                    {regions?.map(region => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
                </Select>
            </FormControl>
            <FormControl fullWidth margin="normal" disabled={!selectedRegion || isLoadingParcels}>
                <InputLabel>Фильтр по участку</InputLabel>
                <Select value={selectedParcel} label="Фильтр по участку" onChange={(e) => setSelectedParcel(e.target.value)}>
                    {parcels?.map(parcel => <MenuItem key={parcel.id} value={parcel.id}>{parcel.name}</MenuItem>)}
                </Select>
            </FormControl>

            {isLoadingGreenhouses && <CircularProgress />}

            {selectedParcel && greenhouses && (
                <TableContainer component={Paper} sx={{mt: 2}}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableCell>ID</TableCell>
                                <TableCell>Название</TableCell>
                                <TableCell>Тип</TableCell>
                                <TableCell align="right">Действия</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {greenhouses.map((greenhouse) => (
                                <TableRow key={greenhouse.id}>
                                    <TableCell>{greenhouse.id}</TableCell>
                                    <TableCell>{greenhouse.name}</TableCell>
                                    <TableCell>{greenhouse.type}</TableCell>
                                    <TableCell align="right">
                                        <IconButton onClick={() => handleOpenEditModal(greenhouse)}><EditIcon /></IconButton>
                                        <IconButton onClick={() => handleOpenDeleteDialog(greenhouse)}><DeleteIcon /></IconButton>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}

            <Modal open={editModalOpen} onClose={handleClose}>
                <Box sx={{ ...modalStyle }}>
                    <Typography variant="h6">Редактировать теплицу</Typography>
                    <TextField label="Новое название" value={editedName} onChange={(e) => setEditedName(e.target.value)} fullWidth margin="normal" />
                    <TextField label="Новый тип" value={editedType} onChange={(e) => setEditedType(e.target.value)} fullWidth margin="normal" />
                    <Button onClick={handleUpdate} variant="contained">Сохранить</Button>
                </Box>
            </Modal>

            <Dialog open={deleteDialogOpen} onClose={handleClose}>
                <DialogTitle>Подтвердите удаление</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить теплицу "{selectedGreenhouse?.name}"?
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

export default GreenhouseManagementPage;