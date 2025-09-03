import React, { useState, useEffect } from 'react';
import { 
  useCreateRegionMutation, 
  useGetRegionsQuery, 
  useCreateLandParcelMutation, 
  useGetLandParcelsByRegionQuery,
  useCreateStructureMutation,
  useGetStructuresByLandParcelQuery,
  useCreatePlotMutation,
  useGetPlotsByStructureQuery,
  useCreateCameraMutation,
  useCreateCropMutation
} from '../features/api/apiSlice';
import { 
  Box, TextField, Button, Typography, CircularProgress, Select, MenuItem, InputLabel, FormControl, Grid
} from '@mui/material';

// --- Form Components ---

const CreateRegionForm = () => {
  const [name, setName] = useState('');
  const [createRegion, { isLoading, isSuccess, isError }] = useCreateRegionMutation();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim()) {
      createRegion({ name }).unwrap().then(() => setName(''));
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2, p: 2, border: '1px solid grey', borderRadius: 1 }}>
      <Typography variant="h6">1. Создать регион</Typography>
      <TextField label="Название региона" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
      <Button type="submit" variant="contained" disabled={isLoading}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
      {isSuccess && <Typography color="green" sx={{ mt: 1 }}>Успешно!</Typography>}
      {isError && <Typography color="error" sx={{ mt: 1 }}>Ошибка.</Typography>}
    </Box>
  );
};

const CreateLandParcelForm = () => {
  const [name, setName] = useState('');
  const [regionId, setRegionId] = useState('');
  const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
  const [createLandParcel, { isLoading, isSuccess, isError }] = useCreateLandParcelMutation();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim() && regionId) {
      createLandParcel({ regionId, name }).unwrap().then(() => { setName(''); setRegionId(''); });
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2, p: 2, border: '1px solid grey', borderRadius: 1 }}>
      <Typography variant="h6">2. Создать участок</Typography>
      <FormControl fullWidth margin="normal" required>
        <InputLabel>Регион</InputLabel>
        <Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)} disabled={isLoadingRegions}>
          {regions?.map(region => <MenuItem key={region.id} value={region.id}>{region.name}</MenuItem>)}
        </Select>
      </FormControl>
      <TextField label="Название участка" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
      <Button type="submit" variant="contained" disabled={isLoading || !regionId}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
      {isSuccess && <Typography color="green" sx={{ mt: 1 }}>Успешно!</Typography>}
      {isError && <Typography color="error" sx={{ mt: 1 }}>Ошибка.</Typography>}
    </Box>
  );
}

const CreateStructureForm = () => {
    const [name, setName] = useState('');
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');

    const { data: regions, isLoading: isLoadingRegions } = useGetRegionsQuery();
    const { data: landParcels, isLoading: isLoadingParcels } = useGetLandParcelsByRegionQuery(regionId, { skip: !regionId });
    const [createStructure, { isLoading, isSuccess, isError }] = useCreateStructureMutation();

    useEffect(() => {
        setLandParcelId('');
    }, [regionId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && landParcelId) {
            createStructure({ landParcelId, name }).unwrap().then(() => { setName(''); setLandParcelId(''); setRegionId(''); });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2, p: 2, border: '1px solid grey', borderRadius: 1 }}>
            <Typography variant="h6">3. Создать сооружение</Typography>
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
            <Button type="submit" variant="contained" disabled={isLoading || !landParcelId}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
            {isSuccess && <Typography color="green" sx={{ mt: 1 }}>Успешно!</Typography>}
            {isError && <Typography color="error" sx={{ mt: 1 }}>Ошибка.</Typography>}
        </Box>
    );
};

const CreatePlotForm = () => {
    const [name, setName] = useState('');
    const [size, setSize] = useState('');
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');
    const [structureId, setStructureId] = useState('');

    const { data: regions } = useGetRegionsQuery();
    const { data: landParcels } = useGetLandParcelsByRegionQuery(regionId, { skip: !regionId });
    const { data: structures } = useGetStructuresByLandParcelQuery(landParcelId, { skip: !landParcelId });
    const [createPlot, { isLoading, isSuccess, isError }] = useCreatePlotMutation();

    useEffect(() => { setLandParcelId(''); setStructureId(''); }, [regionId]);
    useEffect(() => { setStructureId(''); }, [landParcelId]);

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

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2, p: 2, border: '1px solid grey', borderRadius: 1 }}>
            <Typography variant="h6">4. Создать грядку</Typography>
            <FormControl fullWidth margin="normal" required><InputLabel>Регион</InputLabel><Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)}>{regions?.map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!regionId}><InputLabel>Участок</InputLabel><Select value={landParcelId} label="Участок" onChange={(e) => setLandParcelId(e.target.value)}>{landParcels?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!landParcelId}><InputLabel>Сооружение</InputLabel><Select value={structureId} label="Сооружение" onChange={(e) => setStructureId(e.target.value)}>{structures?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select></FormControl>
            <TextField label="Название грядки" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <TextField label="Размер (например, 2x1.5м)" value={size} onChange={(e) => setSize(e.target.value)} fullWidth margin="normal" />
            <Button type="submit" variant="contained" disabled={isLoading || !structureId}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
            {isSuccess && <Typography color="green" sx={{ mt: 1 }}>Успешно!</Typography>}
            {isError && <Typography color="error" sx={{ mt: 1 }}>Ошибка.</Typography>}
        </Box>
    );
};

const CreateCameraForm = () => {
    const [name, setName] = useState('');
    const [rtspPath, setRtspPath] = useState('');
    const [plotId, setPlotId] = useState('');
    // State for cascading selects
    const [regionId, setRegionId] = useState('');
    const [landParcelId, setLandParcelId] = useState('');
    const [structureId, setStructureId] = useState('');

    const { data: regions } = useGetRegionsQuery();
    const { data: landParcels } = useGetLandParcelsByRegionQuery(regionId, { skip: !regionId });
    const { data: structures } = useGetStructuresByLandParcelQuery(landParcelId, { skip: !landParcelId });
    const { data: plots } = useGetPlotsByStructureQuery(structureId, { skip: !structureId });
    const [createCamera, { isLoading, isSuccess, isError }] = useCreateCameraMutation();

    useEffect(() => { setLandParcelId(''); setStructureId(''); setPlotId(''); }, [regionId]);
    useEffect(() => { setStructureId(''); setPlotId(''); }, [landParcelId]);
    useEffect(() => { setPlotId(''); }, [structureId]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (name.trim() && plotId && rtspPath.trim()) {
            createCamera({ plotId, name, rtsp_path_name: rtspPath }).unwrap().then(() => { 
                setName(''); setRtspPath(''); setPlotId(''); setStructureId(''); setLandParcelId(''); setRegionId(''); 
            });
        }
    };

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2, p: 2, border: '1px solid grey', borderRadius: 1 }}>
            <Typography variant="h6">5. Создать и привязать камеру</Typography>
            <FormControl fullWidth margin="normal" required><InputLabel>Регион</InputLabel><Select value={regionId} label="Регион" onChange={(e) => setRegionId(e.target.value)}>{regions?.map(r => <MenuItem key={r.id} value={r.id}>{r.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!regionId}><InputLabel>Участок</InputLabel><Select value={landParcelId} label="Участок" onChange={(e) => setLandParcelId(e.target.value)}>{landParcels?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!landParcelId}><InputLabel>Сооружение</InputLabel><Select value={structureId} label="Сооружение" onChange={(e) => setStructureId(e.target.value)}>{structures?.map(s => <MenuItem key={s.id} value={s.id}>{s.name}</MenuItem>)}</Select></FormControl>
            <FormControl fullWidth margin="normal" required disabled={!structureId}><InputLabel>Грядка</InputLabel><Select value={plotId} label="Грядка" onChange={(e) => setPlotId(e.target.value)}>{plots?.map(p => <MenuItem key={p.id} value={p.id}>{p.name}</MenuItem>)}</Select></FormControl>
            <TextField label="Название камеры" value={name} onChange={(e) => setName(e.target.value)} fullWidth margin="normal" required />
            <TextField label="Путь RTSP (rtsp_path_name)" value={rtspPath} onChange={(e) => setRtspPath(e.target.value)} fullWidth margin="normal" required helperText="Например, plot_1_cam_1" />
            <Button type="submit" variant="contained" disabled={isLoading || !plotId}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
            {isSuccess && <Typography color="green" sx={{ mt: 1 }}>Успешно!</Typography>}
            {isError && <Typography color="error" sx={{ mt: 1 }}>Ошибка.</Typography>}
        </Box>
    );
}

const CreateCropForm = () => {
    const [formState, setFormState] = useState({ name: '', description: '', planting_time: 0, harvest_time: 0 });
    const [createCrop, { isLoading, isSuccess, isError }] = useCreateCropMutation();

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
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2, p: 2, border: '1px solid grey', borderRadius: 1 }}>
            <Typography variant="h6">Создать культуру</Typography>
            <TextField name="name" label="Название культуры" value={formState.name} onChange={handleChange} fullWidth margin="normal" required />
            <TextField name="description" label="Описание" value={formState.description} onChange={handleChange} fullWidth margin="normal" multiline rows={2} />
            <TextField name="planting_time" label="Время посадки (дней)" type="number" value={formState.planting_time} onChange={handleChange} fullWidth margin="normal" />
            <TextField name="harvest_time" label="Время сбора (дней)" type="number" value={formState.harvest_time} onChange={handleChange} fullWidth margin="normal" />
            <Button type="submit" variant="contained" disabled={isLoading}>{isLoading ? <CircularProgress size={24} /> : 'Создать'}</Button>
            {isSuccess && <Typography color="green" sx={{ mt: 1 }}>Успешно!</Typography>}
            {isError && <Typography color="error" sx={{ mt: 1 }}>Ошибка.</Typography>}
        </Box>
    );
}

export const AdminPage = () => {
  return (
    <div>
      <h1>Панель администратора</h1>
      <Typography>Используйте эти формы, чтобы последовательно заполнить систему данными.</Typography>
      <Grid container spacing={2} sx={{mt: 1}}>
        <Grid size={{ xs: 12, md: 6 }}><CreateRegionForm /></Grid>
        <Grid size={{ xs: 12, md: 6 }}><CreateLandParcelForm /></Grid>
        <Grid size={{ xs: 12, md: 6 }}><CreateStructureForm /></Grid>
        <Grid size={{ xs: 12, md: 6 }}><CreatePlotForm /></Grid>
        <Grid size={{ xs: 12, md: 6 }}><CreateCameraForm /></Grid>
        <Grid size={{ xs: 12, md: 6 }}><CreateCropForm /></Grid>
      </Grid>
    </div>
  );
};