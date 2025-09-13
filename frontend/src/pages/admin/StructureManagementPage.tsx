import React, { useState, useMemo } from 'react';
import { 
    useGetRegionsForAdminQuery,
    useGetLandParcelsForAdminQuery,
    useGetStructuresForAdminQuery, 
    useRestoreStructureMutation, 
    useCreateStructureMutation, 
    useUpdateStructureMutation, 
    useDeleteStructureMutation, 
    useGetStructureTypesQuery 
} from '../../features/api/apiSlice';
import type { Structure, Region, LandParcel } from '../../features/api/apiSlice';
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
    Autocomplete,
    Switch,
    FormControlLabel,
    Tooltip
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import RestoreFromTrashIcon from '@mui/icons-material/RestoreFromTrash';

// --- Форма создания (остается почти без изменений) ---
const CreateStructureForm = () => {
    const [name, setName] = useState('');
    const [type, setType] = useState<string | null>(null);
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsForAdminQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsForAdminQuery();
    const { data: structureTypes, isLoading: isLoadingStructureTypes } = useGetStructureTypesQuery();
    const [createStructure, { isLoading }] = useCreateStructureMutation();

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
                    {regions?.filter((r: Region) => !r.deleted_at).map((region: Region) => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
                </Select>
            </FormControl>
            <FormControl fullWidth margin="normal" required disabled={!regionId || isLoadingParcels}>
                <InputLabel>Участок</InputLabel>
                <Select value={landParcelId} label="Участок" onChange={(e) => setLandParcelId(e.target.value)}>
                    {landParcels?.filter((p: LandParcel) => p.region_id === regionId && !p.deleted_at).map((parcel: LandParcel) => <MenuItem key={parcel.id} value={parcel.id}>{parcel.name}</MenuItem>)}
                </Select>
            </FormControl>
            <TextField label="Название сооружения" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <Autocomplete
                freeSolo
                options={structureTypes || []}
                value={type}
                onChange={(_event, newValue) => setType(newValue)}
                onInputChange={(_event, newInputValue) => setType(newInputValue)}
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
    const [showDeleted, setShowDeleted] = useState(false);
    const { data: allStructures, isLoading: isLoadingStructures, isError } = useGetStructuresForAdminQuery();
    const { data: allParcels, isLoading: isLoadingParcels } = useGetLandParcelsForAdminQuery();
    const { data: allRegions, isLoading: isLoadingRegions } = useGetRegionsForAdminQuery();

    const [deleteStructure] = useDeleteStructureMutation();
    const [updateStructure] = useUpdateStructureMutation();
    const [restoreStructure] = useRestoreStructureMutation();

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedStructure, setSelectedStructure] = useState<Structure | null>(null);
    const [editedName, setEditedName] = useState('');
    const [editedType, setEditedType] = useState('');

    const { parcelMap, regionMap } = useMemo(() => {
        const pMap = new Map<string, { name: string, region_id: string }>();
        allParcels?.forEach((p: LandParcel) => pMap.set(p.id, { name: p.name, region_id: p.region_id }));
        const rMap = new Map<string, string>();
        allRegions?.forEach((r: Region) => rMap.set(r.id, r.name));
        return { parcelMap: pMap, regionMap: rMap };
    }, [allParcels, allRegions]);

    const { activeStructures, deletedStructures } = useMemo(() => {
        const active: Structure[] = [];
        const deleted: Structure[] = [];
        allStructures?.forEach(structure => {
            if (structure.deleted_at) {
                deleted.push(structure);
            } else {
                active.push(structure);
            }
        });
        return { activeStructures: active, deletedStructures: deleted };
    }, [allStructures]);

    const structuresToDisplay = showDeleted ? deletedStructures : activeStructures;

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

    const handleRestore = (id: string) => {
        restoreStructure(id);
    };

    const handleUpdate = () => {
        if (selectedStructure && editedName.trim()) {
            updateStructure({ id: selectedStructure.id, name: editedName.trim(), type: editedType });
            handleClose();
        }
    };

    if (isLoadingStructures || isLoadingParcels || isLoadingRegions) return <CircularProgress />;
    if (isError) return <Typography color="error">Не удалось загрузить данные.</Typography>;

    return (
        <Box>
            <CreateStructureForm />
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2, mt: 4 }}>
                <Typography variant="h6">{showDeleted ? 'Удаленные сооружения' : 'Активные сооружения'}</Typography>
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
                            <TableCell>Тип</TableCell>
                            <TableCell>Участок</TableCell>
                            <TableCell>Регион</TableCell>
                            {showDeleted && <TableCell>Дата удаления</TableCell>}
                            <TableCell align="right">Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {structuresToDisplay.map((structure) => {
                            const parcel = parcelMap.get(structure.land_parcel_id);
                            const region = parcel ? regionMap.get(parcel.region_id) : undefined;
                            return (
                                <TableRow key={structure.id} sx={{ backgroundColor: showDeleted ? '#f5f5f5' : 'inherit' }}>
                                    <TableCell>{structure.name}</TableCell>
                                    <TableCell>{structure.type}</TableCell>
                                    <TableCell>{parcel?.name || 'N/A'}</TableCell>
                                    <TableCell>{region || 'N/A'}</TableCell>
                                    {showDeleted && <TableCell>{structure.deleted_at ? new Date(structure.deleted_at).toLocaleString() : '-'}</TableCell>}
                                    <TableCell align="right">
                                        {showDeleted ? (
                                            <Tooltip title="Восстановить">
                                                <IconButton onClick={() => handleRestore(structure.id)}><RestoreFromTrashIcon /></IconButton>
                                            </Tooltip>
                                        ) : (
                                            <>
                                                <Tooltip title="Редактировать">
                                                    <IconButton onClick={() => handleOpenEditModal(structure)}><EditIcon /></IconButton>
                                                </Tooltip>
                                                <Tooltip title="Удалить">
                                                    <IconButton onClick={() => handleOpenDeleteDialog(structure)}><DeleteIcon /></IconButton>
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