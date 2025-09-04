import React, { useState } from 'react';
import { type CatalogItem } from '../../features/api/apiSlice';
import {
    useGetCatalogItemsQuery,
    useCreateCatalogItemMutation,
    useUpdateCatalogItemMutation,
    useDeleteCatalogItemMutation
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
    Select,
    MenuItem,
    FormControl,
    InputLabel,
    Alert
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import AddIcon from '@mui/icons-material/Add';

const modalStyle = {
    position: 'absolute' as 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 600,
    bgcolor: 'background.paper',
    border: '2px solid #000',
    boxShadow: 24,
    p: 4,
};

const CatalogItemForm: React.FC<{ item?: CatalogItem | null, onSave: () => void, onCancel: () => void }> = ({ item, onSave, onCancel }) => {
    const [name, setName] = useState(item?.name || '');
    const [itemType, setItemType] = useState(item?.item_type || 'crop');
    const [description, setDescription] = useState(item?.description || '');
    const [attributes, setAttributes] = useState(JSON.stringify(item?.attributes || {}, null, 2));
    const [jsonError, setJsonError] = useState('');

    const [createItem, { isLoading: isCreating }] = useCreateCatalogItemMutation();
    const [updateItem, { isLoading: isUpdating }] = useUpdateCatalogItemMutation();

    const validateAndGetAttributes = () => {
        try {
            const parsed = JSON.parse(attributes);
            setJsonError('');
            return parsed;
        } catch (e) {
            setJsonError('Invalid JSON format.');
            return null;
        }
    };

    const handleSubmit = () => {
        const parsedAttributes = validateAndGetAttributes();
        if (!parsedAttributes) return;

        const itemData = { name, item_type: itemType, description, attributes: parsedAttributes };

        if (item && item.id) {
            updateItem({ id: item.id, ...itemData }).unwrap().then(onSave);
        } else {
            createItem(itemData).unwrap().then(onSave);
        }
    };

    return (
        <Box>
            <Typography variant="h6">{item ? 'Edit' : 'Create'} Catalog Item</Typography>
            <TextField label="Name" value={name} onChange={e => setName(e.target.value)} fullWidth margin="normal" required />
            <FormControl fullWidth margin="normal">
                <InputLabel>Item Type</InputLabel>
                <Select value={itemType} label="Item Type" onChange={e => setItemType(e.target.value)}>
                    <MenuItem value="crop">Crop</MenuItem>
                    <MenuItem value="poultry">Poultry</MenuItem>
                    <MenuItem value="fish">Fish</MenuItem>
                </Select>
            </FormControl>
            <TextField label="Description" value={description} onChange={e => setDescription(e.target.value)} fullWidth margin="normal" multiline rows={3} />
            <TextField
                label="Attributes (JSON)"
                value={attributes}
                onChange={e => setAttributes(e.target.value)}
                fullWidth
                margin="normal"
                multiline
                rows={6}
                error={!!jsonError}
                helperText={jsonError || 'Enter attributes as a JSON object.'}
            />
            <DialogActions>
                <Button onClick={onCancel}>Cancel</Button>
                <Button onClick={handleSubmit} variant="contained" disabled={isCreating || isUpdating}>Save</Button>
            </DialogActions>
        </Box>
    );
};

const CatalogManagementPage = () => {
    const [itemTypeFilter, setItemTypeFilter] = useState('crop');
    const { data: items, isLoading, isError, refetch } = useGetCatalogItemsQuery(itemTypeFilter, { refetchOnMountOrArgChange: true });
    const [deleteItem] = useDeleteCatalogItemMutation();

    const [modalOpen, setModalOpen] = useState(false);
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [selectedItem, setSelectedItem] = useState<CatalogItem | null>(null);

    const handleOpenCreateModal = () => {
        setSelectedItem(null);
        setModalOpen(true);
    };

    const handleOpenEditModal = (item: CatalogItem) => {
        setSelectedItem(item);
        setModalOpen(true);
    };

    const handleOpenDeleteDialog = (item: CatalogItem) => {
        setSelectedItem(item);
        setDeleteDialogOpen(true);
    };

    const handleClose = () => {
        setModalOpen(false);
        setDeleteDialogOpen(false);
        setSelectedItem(null);
    };

    const handleDelete = () => {
        if (selectedItem) {
            deleteItem(selectedItem.id).unwrap().then(() => handleClose());
        }
    };

    const handleSave = () => {
        handleClose();
        refetch(); // Refetch data after saving
    }

    return (
        <Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h5">Catalog Management</Typography>
                <Button variant="contained" startIcon={<AddIcon />} onClick={handleOpenCreateModal}>Create Item</Button>
            </Box>

            <FormControl sx={{ minWidth: 200, mb: 2 }}>
                <InputLabel>Filter by Type</InputLabel>
                <Select value={itemTypeFilter} label="Filter by Type" onChange={e => setItemTypeFilter(e.target.value)}>
                    <MenuItem value="crop">Crops</MenuItem>
                    <MenuItem value="poultry">Poultry</MenuItem>
                    <MenuItem value="fish">Fish</MenuItem>
                </Select>
            </FormControl>

            {isLoading && <CircularProgress />}
            {isError && <Alert severity="error">Failed to load catalog items.</Alert>}

            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>Name</TableCell>
                            <TableCell>Type</TableCell>
                            <TableCell>Description</TableCell>
                            <TableCell>Attributes</TableCell>
                            <TableCell align="right">Actions</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {items?.map((item) => (
                            <TableRow key={item.id}>
                                <TableCell>{item.name}</TableCell>
                                <TableCell>{item.item_type}</TableCell>
                                <TableCell>{item.description}</TableCell>
                                <TableCell><pre>{JSON.stringify(item.attributes, null, 2)}</pre></TableCell>
                                <TableCell align="right">
                                    <IconButton onClick={() => handleOpenEditModal(item)}><EditIcon /></IconButton>
                                    <IconButton onClick={() => handleOpenDeleteDialog(item)}><DeleteIcon /></IconButton>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>

            <Modal open={modalOpen} onClose={handleClose}>
                <Box sx={modalStyle}>
                    <CatalogItemForm item={selectedItem} onSave={handleSave} onCancel={handleClose} />
                </Box>
            </Modal>

            <Dialog open={deleteDialogOpen} onClose={handleClose}>
                <DialogTitle>Confirm Deletion</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Are you sure you want to delete "{selectedItem?.name}"?
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>Cancel</Button>
                    <Button onClick={handleDelete} color="error">Delete</Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};

export default CatalogManagementPage;