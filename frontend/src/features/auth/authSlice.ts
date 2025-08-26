import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import type { RootState } from '../../store/store';

interface AuthState {
  token: string | null;
  role: string | null;
}

// Helper to decode JWT without external libraries
const decodeJwt = (token: string): { role: string; exp: number } | null => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    const decoded = JSON.parse(jsonPayload);
    return { role: decoded.role, exp: decoded.exp };
  } catch (error) {
    console.error("Failed to decode JWT:", error);
    return null;
  }
};

const initialState: AuthState = {
  token: null,
  role: null,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setCredentials(state, action: PayloadAction<{ accessToken: string }>) {
      const { accessToken } = action.payload;
      const decoded = decodeJwt(accessToken);
      if (decoded && decoded.exp * 1000 > Date.now()) {
        state.token = accessToken;
        state.role = decoded.role;
        localStorage.setItem('token', accessToken);
      } else {
        state.token = null;
        state.role = null;
        localStorage.removeItem('token');
      }
    },
    logout(state) {
      state.token = null;
      state.role = null;
      localStorage.removeItem('token');
    },
    checkInitialAuthState(state) {
        const token = localStorage.getItem('token');
        if (token) {
            const decoded = decodeJwt(token);
            if (decoded && decoded.exp * 1000 > Date.now()) {
                state.token = token;
                state.role = decoded.role;
            } else {
                localStorage.removeItem('token');
            }
        }
    }
  },
});

export const { setCredentials, logout, checkInitialAuthState } = authSlice.actions;

export default authSlice.reducer;

export const selectCurrentUserRole = (state: RootState) => state.auth.role;
export const selectIsLoggedIn = (state: RootState) => !!state.auth.token;
