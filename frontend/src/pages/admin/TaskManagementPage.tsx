import { useMemo } from 'react';
import { useGetTasksQuery, useAcceptTaskMutation, useCompleteTaskMutation, useFailTaskMutation } from '../../features/api/apiSlice';
import type { Task } from '../../features/api/apiSlice';
import {
    Box,
    Typography,
    CircularProgress,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    Button,
    ButtonGroup
} from '@mui/material';

const TaskManagementPage = () => {
    const { data: tasks, isLoading, isError } = useGetTasksQuery();
    const [acceptTask, { isLoading: isAccepting }] = useAcceptTaskMutation();
    const [completeTask, { isLoading: isCompleting }] = useCompleteTaskMutation();
    const [failTask, { isLoading: isFailing }] = useFailTaskMutation();

    const isMutating = isAccepting || isCompleting || isFailing;

    const sortedTasks = useMemo(() => {
        if (!tasks) return [];
        return [...tasks].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    }, [tasks]);

    if (isLoading) return <CircularProgress />;
    if (isError) return <Typography color="error">Не удалось загрузить задачи.</Typography>;

    return (
        <Box>
            <Typography variant="h4" sx={{ mb: 2 }}>Управление задачами</Typography>
            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>ID</TableCell>
                            <TableCell>Название</TableCell>
                            <TableCell>Статус</TableCell>
                            <TableCell>Время создания</TableCell>
                            <TableCell>ID операции</TableCell>
                            <TableCell>Исполнитель</TableCell>
                            <TableCell>Действия</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {sortedTasks.map((task: Task) => (
                            <TableRow key={task.id}>
                                <TableCell>{task.id}</TableCell>
                                <TableCell>{task.title}</TableCell>
                                <TableCell>{task.status}</TableCell>
                                <TableCell>{new Date(task.created_at).toLocaleString()}</TableCell>
                                <TableCell>{task.operation_id}</TableCell>
                                <TableCell>{task.assignee_id || '-'}</TableCell>
                                <TableCell>
                                    <ButtonGroup variant="contained" size="small" disabled={isMutating}>
                                        <Button onClick={() => acceptTask(task.id)} disabled={task.status !== 'new'}>
                                            Принять
                                        </Button>
                                        <Button color="success" onClick={() => completeTask(task.id)} disabled={task.status !== 'in_progress'}>
                                            Завершить
                                        </Button>
                                        <Button color="error" onClick={() => failTask(task.id)} disabled={task.status !== 'in_progress'}>
                                            Провалить
                                        </Button>
                                    </ButtonGroup>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </Box>
    );
};

export default TaskManagementPage;