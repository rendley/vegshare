import React, { useState } from 'react';
import { 
    useGetRegionsQuery,
    useCreateRegionMutation,
    useUpdateRegionMutation,
    useDeleteRegionMutation
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
const CreateRegionForm = () => {
    const [name, setName] = useState('');
    const [createRegion, { isLoading }] = useCreateRegionMutation();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim()) {
            createRegion({ name }).unwrap().then(() => setName(''));
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mb: 4 }}>
            <Typography variant="h6">Создать новый регион</Typography>
            <TextField label="Название региона" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <Button type="submit" variant="contained" disabled={isLoading}>
                {isLoading ? <CircularProgress size={24} /> : 'Создать'}
            </Button>
        </Box>
    );
};

// --- Таблица управления ---
const RegionManagementPage = () => {
    const { data: regions, isLoading, isError } = useGetRegionsQuery();
    const [deleteRegion] = useDeleteRegionMutation();
    const [updateRegion] = useUpdateRegionMutation();

    // Состояние для модальных окон
    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedRegion, setSelectedRegion] = useState<{ id: string; name: string } | null>(null);
    const [editedName, setEditedName] = useState('');

    // Обработчики для открытия/закрытия окон
    const handleOpenEditModal = (region: { id: string; name: string }) => {
        setSelectedRegion(region);
        setEditedName(region.name);
        setEditModalOpen(true);
    };

    const handleOpenDeleteDialog = (region: { id: string; name: string }) => {
        setSelectedRegion(region);
        setDeleteDialogOpen(true);
    };

    const handleClose = () => {
        setEditModalOpen(false);
        setDeleteDialogOpen(false);
        setSelectedRegion(null);
    };

    // Обработчики действий
    const handleDelete = () => {
        if (selectedRegion) {
            deleteRegion(selectedRegion.id);
            handleClose();
        }
    };

    const handleUpdate = () => {
        if (selectedRegion && editedName.trim()) {
            updateRegion({ id: selectedRegion.id, name: editedName.trim() });
            handleClose();
        }
    };

    if (isLoading) return <CircularProgress />;
    if (isError) return <Typography color="error">Не удалось загрузить регионы.</Typography>;

    return (
        <Box>
            <CreateRegionForm />

            <Typography variant="h6" sx={{ mb: 2 }}>Список регионов</Typography>
            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>ID</TableCell>
                            <TableCell>Название</TableCell>
                            <TableCell align="right">Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {regions?.map((region) => (
                            <TableRow key={region.id}>
                                <TableCell>{region.id}</TableCell>
                                <TableCell>{region.name}</TableCell>
                                <TableCell align="right">
                                    <IconButton onClick={() => handleOpenEditModal(region)}><EditIcon /></IconButton>
                                    <IconButton onClick={() => handleOpenDeleteDialog(region)}><DeleteIcon /></IconButton>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>

            {/* Модальное окно редактирования */}
            <Modal open={editModalOpen} onClose={handleClose}>
                <Box sx={{ ...modalStyle }}>
                    <Typography variant="h6">Редактировать регион</Typography>
                    <TextField
                        label="Новое название"
                        value={editedName}
                        onChange={(e) => setEditedName(e.target.value)}
                        fullWidth
                        margin="normal"
                    />
                    <Button onClick={handleUpdate} variant="contained">Сохранить</Button>
                </Box>
            </Modal>

            {/* Диалог подтверждения удаления */}
            <Dialog open={deleteDialogOpen} onClose={handleClose}>
                <DialogTitle>Подтвердите удаление</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить регион "{selectedRegion?.name}"? Это действие необратимо.
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

export default RegionManagementPage;
