import { useGetUsersQuery, useUpdateUserRoleMutation } from '../../features/api/apiSlice';
import type { User } from '../../features/api/apiSlice';
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
    Select,
    MenuItem,
    FormControl
} from '@mui/material';

const UserManagementPage = () => {
    const { data: users, isLoading, isError } = useGetUsersQuery();
    const [updateUserRole, { isLoading: isUpdating }] = useUpdateUserRoleMutation();

    const handleRoleChange = (userId: string, newRole: string) => {
        if (newRole === 'admin' || newRole === 'user') {
            updateUserRole({ userId, role: newRole });
        }
    };

    if (isLoading) return <CircularProgress />;
    if (isError) return <Typography color="error">Не удалось загрузить пользователей.</Typography>;

    return (
        <Box>
            <Typography variant="h4" sx={{ mb: 2 }}>Управление пользователями</Typography>
            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>ID</TableCell>
                            <TableCell>Имя</TableCell>
                            <TableCell>Email</TableCell>
                            <TableCell>Роль</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {users?.map((user: User) => (
                            <TableRow key={user.id}>
                                <TableCell>{user.id}</TableCell>
                                <TableCell>{user.name}</TableCell>
                                <TableCell>{user.email}</TableCell>
                                <TableCell>
                                    <FormControl size="small">
                                        <Select
                                            value={user.role}
                                            onChange={(e) => handleRoleChange(user.id, e.target.value)}
                                            disabled={isUpdating}
                                        >
                                            <MenuItem value="user">user</MenuItem>
                                            <MenuItem value="admin">admin</MenuItem>
                                        </Select>
                                    </FormControl>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </Box>
    );
};

export default UserManagementPage;
