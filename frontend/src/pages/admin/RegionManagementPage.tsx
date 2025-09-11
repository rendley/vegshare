import React, { useState, useMemo } from 'react';
import { 
    useGetRegionsForAdminQuery, // Используем новый хук
    useCreateRegionMutation,
    useUpdateRegionMutation,
    useDeleteRegionMutation,
    useRestoreRegionMutation, // Используем новый хук
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
    Switch,
    FormControlLabel,
    Tooltip
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import RestoreFromTrashIcon from '@mui/icons-material/RestoreFromTrash';

// --- Форма создания (остается без изменений) ---
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

// --- Таблица управления (переписана) ---
const RegionManagementPage = () => {
    const [showDeleted, setShowDeleted] = useState(false);

    // Загружаем все регионы (активные и удаленные) один раз
    const { data: allRegions, isLoading, isError } = useGetRegionsForAdminQuery();

    const [deleteRegion] = useDeleteRegionMutation();
    const [updateRegion] = useUpdateRegionMutation();
    const [restoreRegion] = useRestoreRegionMutation();

    // Состояния для модальных окон
    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedRegion, setSelectedRegion] = useState<{ id: string; name: string } | null>(null);
    const [editedName, setEditedName] = useState('');

    // Фильтрация регионов на клиенте
    const { activeRegions, deletedRegions } = useMemo(() => {
        const active: any[] = [];
        const deleted: any[] = [];
        allRegions?.forEach(region => {
            if (region.deleted_at) {
                deleted.push(region);
            } else {
                active.push(region);
            }
        });
        return { activeRegions: active, deletedRegions: deleted };
    }, [allRegions]);

    const regionsToDisplay = showDeleted ? deletedRegions : activeRegions;

    // --- Обработчики модальных окон ---
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

    // --- Обработчики действий ---
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

    const handleRestore = (regionId: string) => {
        restoreRegion(regionId);
    };

    if (isLoading) return <CircularProgress />;
    if (isError) return <Typography color="error">Не удалось загрузить регионы.</Typography>;

    return (
        <Box>
            <CreateRegionForm />

            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6">{showDeleted ? 'Удаленные регионы' : 'Активные регионы'}</Typography>
                <FormControlLabel
                    control={<Switch checked={showDeleted} onChange={(e) => setShowDeleted(e.target.checked)} />}
                    label="Показать удаленные"
                />
            </Box>

            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>ID</TableCell>
                            <TableCell>Название</TableCell>
                            {showDeleted && <TableCell>Дата удаления</TableCell>}
                            <TableCell align="right">Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {regionsToDisplay.map((region) => (
                            <TableRow key={region.id} sx={{ backgroundColor: showDeleted ? '#f5f5f5' : 'inherit' }}>
                                <TableCell>{region.id}</TableCell>
                                <TableCell>{region.name}</TableCell>
                                {showDeleted && <TableCell>{new Date(region.deleted_at).toLocaleString()}</TableCell>}
                                <TableCell align="right">
                                    {showDeleted ? (
                                        <Tooltip title="Восстановить">
                                            <IconButton onClick={() => handleRestore(region.id)}><RestoreFromTrashIcon /></IconButton>
                                        </Tooltip>
                                    ) : (
                                        <>
                                            <Tooltip title="Редактировать">
                                                <IconButton onClick={() => handleOpenEditModal(region)}><EditIcon /></IconButton>
                                            </Tooltip>
                                            <Tooltip title="Удалить">
                                                <IconButton onClick={() => handleOpenDeleteDialog(region)}><DeleteIcon /></IconButton>
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
                        Вы уверены, что хотите удалить регион "{selectedRegion?.name}"?
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