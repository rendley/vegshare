import { Box, Drawer, List, ListItem, ListItemButton, ListItemText, Toolbar } from '@mui/material';
import { NavLink, Outlet } from 'react-router-dom';

const drawerWidth = 240;

const AdminLayout = () => {
  return (
    <Box sx={{ display: 'flex' }}>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box' },
        }}
      >
        <Toolbar />
        <Box sx={{ overflow: 'auto' }}>
          <List>
            <ListItem disablePadding>
              <ListItemButton component={NavLink} to="/admin/regions">
                <ListItemText primary="Регионы" />
              </ListItemButton>
            </ListItem>
            {/* Ссылки на другие разделы будут здесь */}
          </List>
        </Box>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Toolbar />
        {/* Здесь будут отображаться дочерние компоненты (страницы управления) */}
        <Outlet />
      </Box>
    </Box>
  );
};

export default AdminLayout;
