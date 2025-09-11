import { combineReducers, configureStore, type PayloadAction } from '@reduxjs/toolkit';
import { apiSlice } from '../features/api/apiSlice';
import authReducer, { logout } from '../features/auth/authSlice';
import notificationReducer from './notificationSlice';

// Комбинируем все редьюсеры в один корневой редьюсер
const rootReducer = combineReducers({
  [apiSlice.reducerPath]: apiSlice.reducer,
  auth: authReducer,
  notification: notificationReducer,
});

// Создаем обертку над корневым редьюсером, которая будет сбрасывать состояние при выходе
const resettableRootReducer = (state: ReturnType<typeof rootReducer> | undefined, action: PayloadAction) => {
  if (action.type === logout.type) {
    // При вызове logout, передаем undefined, чтобы сбросить состояние к начальному
    return rootReducer(undefined, action);
  }
  return rootReducer(state, action);
};

export const store = configureStore({
  reducer: resettableRootReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(apiSlice.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
