import React, { useState, useEffect } from 'react';
import type { Structure } from '../../features/api/apiSlice';
import {
    useGetRegionsQuery,
    useGetLandParcelsByRegionQuery,
    useGetStructuresByLandParcelQuery,
    useGetStructureTypesQuery, // Импортируем новый хук
    useCreateStructureMutation,
    useUpdateStructureMutation,
    useDeleteStructureMutation
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
    Autocomplete // Импортируем Autocomplete
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

// --- Форма создания ---
const CreateStructureForm = () => {
    const [name, setName] = useState('');
    const [type, setType] = useState<string | null>(null);
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(regionId, { skip: !regionId });
    const { data: structureTypes, isLoading: isLoadingStructureTypes } = useGetStructureTypesQuery(); // Получаем типы
    const [createStructure, { isLoading }] = useCreateStructureMutation();

    useEffect(() => {
        setLandParcelId('');
    }, [regionId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && landParcelId && type && type.trim()) {
            createStructure({ landParcelId, name, type }).unwrap().then(() => {
                setName('');
                setType(null);
                setLandParcelId('');
                setRegionId('');
            });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новое сооружение</Typography>
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
            <TextField label="Название сооружения" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <Autocomplete
                freeSolo
                options={structureTypes || []}
                value={type}
                onChange={(_event, newValue) => {
                    setType(newValue);
                }}
                onInputChange={(_event, newInputValue) => {
                    setType(newInputValue);
                }}
                renderInput={(params) => (
                    <TextField
                        {...params}
                        label="Тип"
                        margin="normal"
                        required
                        InputProps={{
                            ...params.InputProps,
                            endAdornment: (
                                <>
                                    {isLoadingStructureTypes ? <CircularProgress color="inherit" size={20} /> : null}
                                    {params.InputProps.endAdornment}
                                </>
                            ),
                        }}
                    />
                )}
            />
            <Button type="submit" variant="contained" disabled={isLoading || !landParcelId || !type}>
                {isLoading ? <CircularProgress size={24} /> : 'Создать'}
            </Button>
        </Box>
    );
}

// --- Таблица управления ---
const StructureManagementPage = () => {
    const [selectedRegion, setSelectedRegion] = useState('');
    const [selectedParcel, setSelectedParcel] = useState('');
    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: parcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(selectedRegion, { skip: !selectedRegion });
    const { data: structures, isLoading: isLoadingStructures } = useGetStructuresByLandParcelQuery(selectedParcel, { skip: !selectedParcel });

    const [deleteStructure] = useDeleteStructureMutation();
    const [updateStructure] = useUpdateStructureMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedStructure, setSelectedStructure] = useState<Structure | null>(null);
    const [editedName, setEditedName] = useState('');
    const [editedType, setEditedType] = useState('');

    useEffect(() => { setSelectedParcel(''); }, [selectedRegion]);

    const handleOpenEditModal = (structure: Structure) => {
        setSelectedStructure(structure);
        setEditedName(structure.name);
        setEditedType(structure.type || '');
        setEditModalOpen(true);
    };

    const handleOpenDeleteDialog = (structure: Structure) => {
        setSelectedStructure(structure);
        setDeleteDialogOpen(true);
    };

    const handleClose = () => {
        setEditModalOpen(false);
        setDeleteDialogOpen(false);
        setSelectedStructure(null);
    };

    const handleDelete = () => {
        if (selectedStructure) {
            deleteStructure(selectedStructure.id);
            handleClose();
        }
    };

    const handleUpdate = () => {
        if (selectedStructure && editedName.trim()) {
            updateStructure({ id: selectedStructure.id, name: editedName.trim(), type: editedType });
            handleClose();
        }
    };

    return (
        <Box>
            <CreateStructureForm />
            <Typography variant="h6" sx={{ mb: 2 }}>Список сооружений</Typography>
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

            {isLoadingStructures && <CircularProgress />}

            {selectedParcel && structures && (
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
                            {structures.map((structure) => (
                                <TableRow key={structure.id}>
                                    <TableCell>{structure.id}</TableCell>
                                    <TableCell>{structure.name}</TableCell>
                                    <TableCell>{structure.type}</TableCell>
                                    <TableCell align="right">
                                        <IconButton onClick={() => handleOpenEditModal(structure)}><EditIcon /></IconButton>
                                        <IconButton onClick={() => handleOpenDeleteDialog(structure)}><DeleteIcon /></IconButton>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}

            <Modal open={editModalOpen} onClose={handleClose}>
                <Box sx={{ ...modalStyle }}>
                    <Typography variant="h6">Редактировать сооружение</Typography>
                    <TextField label="Новое название" value={editedName} onChange={(e) => setEditedName(e.target.value)} fullWidth margin="normal" />
                    <TextField label="Новый тип" value={editedType} onChange={(e) => setEditedType(e.target.value)} fullWidth margin="normal" />
                    <Button onClick={handleUpdate} variant="contained">Сохранить</Button>
                </Box>
            </Modal>

            <Dialog open={deleteDialogOpen} onClose={handleClose}>
                <DialogTitle>Подтвердите удаление</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить сооружение "{selectedStructure?.name}"?
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

export default StructureManagementPage;