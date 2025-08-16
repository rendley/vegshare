import { configureStore } from '@reduxjs/toolkit';
import { apiSlice } from '../features/api/apiSlice';

export const store = configureStore({
  reducer: {
    // Добавляем сгенерированный редюсер от apiSlice
    [apiSlice.reducerPath]: apiSlice.reducer,
  },
  // Добавляем middleware от RTK Query. Это необходимо для кеширования, инвалидации и т.д.
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(apiSlice.middleware),
});

// Типы для всего стора, чтобы использовать их в приложении
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
