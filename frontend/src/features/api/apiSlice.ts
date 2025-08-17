import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

const BASE_URL = '/api/v1';

// --- Типы данных ---
// TODO: Вынести в /types
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
  baseQuery: fetchBaseQuery({ baseUrl: BASE_URL }),
  tagTypes: ['Plot', 'Lease', 'PlotCrop'],
  endpoints: (builder) => ({
    // QUERIES
    getRegions: builder.query<Region[], void>({ query: () => 'regions' }),
    getLandParcelsByRegion: builder.query<LandParcel[], string>({ query: (regionId) => `regions/${regionId}/land-parcels` }),
    getGreenhousesByLandParcel: builder.query<Greenhouse[], string>({ query: (landParcelId) => `land-parcels/${landParcelId}/greenhouses` }),
    getPlotsByGreenhouse: builder.query<Plot[], string>({
      query: (greenhouseId) => `greenhouses/${greenhouseId}/plots`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Plot' as const, id })), { type: 'Plot', id: 'LIST' }] : [{ type: 'Plot', id: 'LIST' }],
    }),
    getMyLeases: builder.query<PlotLease[], void>({
      query: () => 'me/leases',
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Lease' as const, id })), { type: 'Lease', id: 'LIST' }] : [{ type: 'Lease', id: 'LIST' }],
    }),
    getAvailableCrops: builder.query<Crop[], void>({ query: () => 'crops' }),
    getPlotCrops: builder.query<PlotCrop[], string>({ 
      query: (plotId) => `plots/${plotId}/plantings`, // Предполагаем, что такой эндпоинт будет
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'PlotCrop' as const, id })), { type: 'PlotCrop', id: 'LIST' }] : [{ type: 'PlotCrop', id: 'LIST' }],
    }),

    // MUTATIONS
    leasePlot: builder.mutation<PlotLease, string>({
      query: (plotId) => ({ url: `plots/${plotId}/lease`, method: 'POST' }),
      invalidatesTags: [{ type: 'Plot', id: 'LIST' }, { type: 'Lease', id: 'LIST' }],
    }),
    plantCrop: builder.mutation<PlotCrop, { plotId: string; cropId: string }>({ 
      query: ({ plotId, cropId }) => ({
        url: `plots/${plotId}/plantings`,
        method: 'POST',
        body: { crop_id: cropId },
      }),
      invalidatesTags: ['PlotCrop'],
    }),
    removeCrop: builder.mutation<void, { plotId: string; plantingId: string }>({ 
      query: ({ plotId, plantingId }) => ({
        url: `plots/${plotId}/plantings/${plantingId}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['PlotCrop'],
    }),
    performAction: builder.mutation<void, { plotId: string; action: string }>({ 
      query: ({ plotId, action }) => ({
        url: `plots/${plotId}/actions`,
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
  useLeasePlotMutation,
  usePlantCropMutation,
  useRemoveCropMutation,
  usePerformActionMutation,
} = apiSlice;