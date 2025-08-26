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

export interface Camera {
  id: string;
  name: string;
  plot_id: string;
  rtsp_path_name: string;
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

interface AuthRequest {
  email: string;
  password: string;
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
  tagTypes: ['Region', 'LandParcel', 'Greenhouse', 'Plot', 'Lease', 'PlotCrop', 'Camera', 'Crop'],
  endpoints: (builder) => ({
    // QUERIES
    getRegions: builder.query<Region[], void>({
      query: () => 'farm/regions',
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Region' as const, id })), { type: 'Region', id: 'LIST' }] : [{ type: 'Region', id: 'LIST' }],
    }),
    getLandParcelsByRegion: builder.query<LandParcel[], string>({
      query: (regionId) => `farm/regions/${regionId}/land-parcels`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'LandParcel' as const, id })), { type: 'LandParcel', id: 'LIST' }] : [{ type: 'LandParcel', id: 'LIST' }],
    }),
    getGreenhousesByLandParcel: builder.query<Greenhouse[], string>({
      query: (landParcelId) => `farm/land-parcels/${landParcelId}/greenhouses`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Greenhouse' as const, id })), { type: 'Greenhouse', id: 'LIST' }] : [{ type: 'Greenhouse', id: 'LIST' }],
    }),
    getPlotsByGreenhouse: builder.query<Plot[], string>({
      query: (greenhouseId) => `plots?greenhouse_id=${greenhouseId}`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Plot' as const, id })), { type: 'Plot', id: 'LIST' }] : [{ type: 'Plot', id: 'LIST' }],
    }),
    getCamerasByPlot: builder.query<Camera[], string>({
      query: (plotId) => `plots/${plotId}/cameras`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Camera' as const, id })), { type: 'Camera', id: 'LIST' }] : [{ type: 'Camera', id: 'LIST' }],
    }),
    getMyLeases: builder.query<PlotLease[], void>({
      query: () => 'leasing/me/leases',
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Lease' as const, id })), { type: 'Lease', id: 'LIST' }] : [{ type: 'Lease', id: 'LIST' }],
    }),
    getAvailableCrops: builder.query<Crop[], void>({
      query: () => 'catalog/crops',
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Crop' as const, id })), { type: 'Crop', id: 'LIST' }] : [{ type: 'Crop', id: 'LIST' }],
    }),
    getPlotCrops: builder.query<PlotCrop[], string>({
      query: (plotId) => `operations/plots/${plotId}/plantings`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'PlotCrop' as const, id })), { type: 'PlotCrop', id: 'LIST' }] : [{ type: 'PlotCrop', id: 'LIST' }],
    }),

    // MUTATIONS
    login: builder.mutation<AuthResponse, AuthRequest>({
      query: (credentials) => ({
        url: 'auth/login',
        method: 'POST',
        body: credentials,
      }),
    }),
    register: builder.mutation<AuthResponse, AuthRequest>({
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

    // Admin Mutations
    createRegion: builder.mutation<Region, { name: string }>({
      query: (body) => ({
        url: 'farm/regions',
        method: 'POST',
        body,
      }),
      invalidatesTags: [{ type: 'Region', id: 'LIST' }],
    }),
    createLandParcel: builder.mutation<LandParcel, { regionId: string; name: string }>({
      query: ({ regionId, ...body }) => ({
        url: `farm/regions/${regionId}/land-parcels`,
        method: 'POST',
        body,
      }),
      invalidatesTags: [{ type: 'LandParcel', id: 'LIST' }],
    }),
    createGreenhouse: builder.mutation<Greenhouse, { landParcelId: string; name: string }>({
      query: ({ landParcelId, ...body }) => ({
        url: `farm/land-parcels/${landParcelId}/greenhouses`,
        method: 'POST',
        body,
      }),
      invalidatesTags: [{ type: 'Greenhouse', id: 'LIST' }],
    }),
    createPlot: builder.mutation<Plot, { greenhouseId: string; name: string }>({
      query: (body) => ({
        url: 'plots',
        method: 'POST',
        body,
      }),
      invalidatesTags: [{ type: 'Plot', id: 'LIST' }],
    }),
    createCamera: builder.mutation<Camera, { plotId: string; name: string; rtsp_path_name: string }>({
      query: ({ plotId, ...body }) => ({
        url: `plots/${plotId}/cameras`,
        method: 'POST',
        body,
      }),
      invalidatesTags: [{ type: 'Camera', id: 'LIST' }],
    }),
    createCrop: builder.mutation<Crop, Partial<Crop>>({
      query: (body) => ({
        url: 'catalog/crops',
        method: 'POST',
        body,
      }),
      invalidatesTags: [{ type: 'Crop', id: 'LIST' }],
    }),
  }),
});

// --- Экспорт хуков ---
export const {
  useGetRegionsQuery,
  useGetLandParcelsByRegionQuery,
  useGetGreenhousesByLandParcelQuery,
  useGetPlotsByGreenhouseQuery,
  useGetCamerasByPlotQuery,
  useGetMyLeasesQuery,
  useGetAvailableCropsQuery,
  useGetPlotCropsQuery,
  useLoginMutation,
  useRegisterMutation,
  useLeasePlotMutation,
  usePlantCropMutation,
  useRemoveCropMutation,
  usePerformActionMutation,
  useCreateRegionMutation,
  useCreateLandParcelMutation,
  useCreateGreenhouseMutation,
  useCreatePlotMutation,
  useCreateCameraMutation,
  useCreateCropMutation,
} = apiSlice;