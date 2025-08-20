import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

const BASE_URL = '/api/v1';

// --- Типы данных ---
// TODO: Вынести в /types
interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user_id: string;
}

interface Region {
  id: string;
  name: string;
}

interface LandParcel {
  id: string;
  name: string;
}

interface Greenhouse {
  id: string;
  name: string;
}

export interface Plot {
  id: string;
  name: string;
  status: string;
}

export interface PlotLease {
  id: string;
  plot_id: string;
}

export interface Crop {
  id: string;
  name: string;
}

export interface PlotCrop {
  id: string;
  plot_id: string;
  crop_id: string;
  status: string;
}

// --- Слайс API ---
export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: BASE_URL,
    prepareHeaders: (headers) => {
      const token = localStorage.getItem('token');
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
      return headers;
    },
  }),
  tagTypes: ['Plot', 'Lease', 'PlotCrop'],
  endpoints: (builder) => ({
    // QUERIES
    getRegions: builder.query<Region[], void>({ query: () => 'farm/regions' }),
    getLandParcelsByRegion: builder.query<LandParcel[], string>({ query: (regionId) => `farm/regions/${regionId}/land-parcels` }),
    getGreenhousesByLandParcel: builder.query<Greenhouse[], string>({ query: (landParcelId) => `farm/land-parcels/${landParcelId}/greenhouses` }),
    getPlotsByGreenhouse: builder.query<Plot[], string>({
      query: (greenhouseId) => `farm/greenhouses/${greenhouseId}/plots`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Plot' as const, id })), { type: 'Plot', id: 'LIST' }] : [{ type: 'Plot', id: 'LIST' }],
    }),
    getMyLeases: builder.query<PlotLease[], void>({
      query: () => 'leasing/me/leases',
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Lease' as const, id })), { type: 'Lease', id: 'LIST' }] : [{ type: 'Lease', id: 'LIST' }],
    }),
    getAvailableCrops: builder.query<Crop[], void>({ query: () => 'catalog/crops' }),
    getPlotCrops: builder.query<PlotCrop[], string>({ 
      query: (plotId) => `operations/plots/${plotId}/plantings`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'PlotCrop' as const, id })), { type: 'PlotCrop', id: 'LIST' }] : [{ type: 'PlotCrop', id: 'LIST' }],
    }),

    // MUTATIONS
    login: builder.mutation<AuthResponse, { email, password }>({ 
      query: (credentials) => ({
        url: 'auth/login',
        method: 'POST',
        body: credentials,
      }),
    }),
    register: builder.mutation<AuthResponse, { email, password }>({ 
      query: (credentials) => ({
        url: 'auth/register',
        method: 'POST',
        body: credentials,
      }),
    }),
    leasePlot: builder.mutation<PlotLease, string>({
      query: (plotId) => ({ url: `leasing/plots/${plotId}/lease`, method: 'POST' }),
      invalidatesTags: [{ type: 'Plot', id: 'LIST' }, { type: 'Lease', id: 'LIST' }],
    }),
    plantCrop: builder.mutation<PlotCrop, { plotId: string; cropId: string }>({ 
      query: ({ plotId, cropId }) => ({
        url: `operations/plots/${plotId}/plantings`,
        method: 'POST',
        body: { crop_id: cropId },
      }),
      invalidatesTags: ['PlotCrop'],
    }),
    removeCrop: builder.mutation<void, { plotId: string; plantingId: string }>({ 
      query: ({ plotId, plantingId }) => ({
        url: `operations/plots/${plotId}/plantings/${plantingId}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['PlotCrop'],
    }),
    performAction: builder.mutation<void, { plotId: string; action: string }>({ 
      query: ({ plotId, action }) => ({
        url: `operations/plots/${plotId}/actions`,
        method: 'POST',
        body: { action },
      }),
    }),
  }),
});

// --- Экспорт хуков ---
export const {
  useGetRegionsQuery,
  useGetLandParcelsByRegionQuery,
  useGetGreenhousesByLandParcelQuery,
  useGetPlotsByGreenhouseQuery,
  useGetMyLeasesQuery,
  useGetAvailableCropsQuery,
  useGetPlotCropsQuery,
  useLoginMutation,
  useRegisterMutation,
  useLeasePlotMutation,
  usePlantCropMutation,
  useRemoveCropMutation,
  usePerformActionMutation,
} = apiSlice;