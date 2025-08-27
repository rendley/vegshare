import React, { useState } from 'react';
import type { Crop } from '../../features/api/apiSlice';
import {
    useGetAvailableCropsQuery,
    useCreateCropMutation,
    useUpdateCropMutation,
    useDeleteCropMutation
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
    DialogTitle
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

// --- Форма создания ---
const CreateCropForm = () => {
    const [formState, setFormState] = useState({ name: '', description: '', planting_time: 0, harvest_time: 0 });
    const [createCrop, { isLoading }] = useCreateCropMutation();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormState(prev => ({ ...prev, [name]: name.includes('time') ? parseInt(value, 10) || 0 : value }));
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (formState.name.trim()) {
            createCrop(formState).unwrap().then(() => {
                setFormState({ name: '', description: '', planting_time: 0, harvest_time: 0 });
            });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новую культуру</Typography>
            <TextField name="name" label="Название" value={formState.name} onChange={handleChange} fullWidth margin="normal" required />
            <TextField name="description" label="Описание" value={formState.description} onChange={handleChange} fullWidth margin="normal" multiline rows={2} />
            <TextField name="planting_time" label="Время посадки (дней)" type="number" value={formState.planting_time} onChange={handleChange} fullWidth margin="normal" />
            <TextField name="harvest_time" label="Время сбора (дней)" type="number" value={formState.harvest_time} onChange={handleChange} fullWidth margin="normal" />
            <Button type="submit" variant="contained" disabled={isLoading}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
        </Box>
    );
}

// --- Таблица управления ---
const CropManagementPage = () => {
    const { data: crops, isLoading, isError } = useGetAvailableCropsQuery();
    const [deleteCrop] = useDeleteCropMutation();
    const [updateCrop] = useUpdateCropMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedCrop, setSelectedCrop] = useState<Crop | null>(null);
    const [editedFormState, setEditedFormState] = useState({ name: '', description: '', planting_time: 0, harvest_time: 0 });

    const handleOpenEditModal = (crop: Crop) => {
        setSelectedCrop(crop);
        setEditedFormState({ 
            name: crop.name, 
            description: crop.description || '', 
            planting_time: crop.planting_time || 0, 
            harvest_time: crop.harvest_time || 0 
        });
        setEditModalOpen(true);
    };

    const handleOpenDeleteDialog = (crop: Crop) => {
        setSelectedCrop(crop);
        setDeleteDialogOpen(true);
    };

    const handleClose = () => {
        setEditModalOpen(false);
        setDeleteDialogOpen(false);
        setSelectedCrop(null);
    };

    const handleDelete = () => {
        if (selectedCrop) {
            deleteCrop(selectedCrop.id);
            handleClose();
        }
    };

    const handleUpdate = () => {
        if (selectedCrop && editedFormState.name.trim()) {
            updateCrop({ id: selectedCrop.id, ...editedFormState });
            handleClose();
        }
    };

    const handleEditChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setEditedFormState(prev => ({ ...prev, [name]: name.includes('time') ? parseInt(value, 10) || 0 : value }));
    };

    if (isLoading) return <CircularProgress />;
    if (isError) return <Typography color="error">Не удалось загрузить культуры.</Typography>;

    return (
        <Box>
            <CreateCropForm />
            <Typography variant="h6" sx={{ mb: 2 }}>Каталог культур</Typography>
            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>ID</TableCell>
                            <TableCell>Название</TableCell>
                            <TableCell>Описание</TableCell>
                            <TableCell>Время посадки</TableCell>
                            <TableCell>Время сбора</TableCell>
                            <TableCell align="right">Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {crops?.map((crop) => (
                            <TableRow key={crop.id}>
                                <TableCell>{crop.id}</TableCell>
                                <TableCell>{crop.name}</TableCell>
                                <TableCell>{crop.description}</TableCell>
                                <TableCell>{crop.planting_time}</TableCell>
                                <TableCell>{crop.harvest_time}</TableCell>
                                <TableCell align="right">
                                    <IconButton onClick={() => handleOpenEditModal(crop)}><EditIcon /></IconButton>
                                    <IconButton onClick={() => handleOpenDeleteDialog(crop)}><DeleteIcon /></IconButton>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>

            <Modal open={editModalOpen} onClose={handleClose}>
                <Box sx={{ ...modalStyle }}>
                    <Typography variant="h6">Редактировать культуру</Typography>
                    <TextField name="name" label="Название" value={editedFormState.name} onChange={handleEditChange} fullWidth margin="normal" required />
                    <TextField name="description" label="Описание" value={editedFormState.description} onChange={handleEditChange} fullWidth margin="normal" multiline rows={2} />
                    <TextField name="planting_time" label="Время посадки (дней)" type="number" value={editedFormState.planting_time} onChange={handleEditChange} fullWidth margin="normal" />
                    <TextField name="harvest_time" label="Время сбора (дней)" type="number" value={editedFormState.harvest_time} onChange={handleEditChange} fullWidth margin="normal" />
                    <Button onClick={handleUpdate} variant="contained">Сохранить</Button>
                </Box>
            </Modal>

            <Dialog open={deleteDialogOpen} onClose={handleClose}>
                <DialogTitle>Подтвердите удаление</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить культуру "{selectedCrop?.name}"?
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

export default CropManagementPage;