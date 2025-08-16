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
  region_id: string;
}

interface Greenhouse {
  id: string;
  name: string;
  land_parcel_id: string;
}

export interface Plot {
  id: string;
  name: string;
  status: string;
  greenhouse_id: string;
}

export interface PlotLease {
  id: string;
  plot_id: string;
  user_id: string;
  start_date: string;
  end_date: string;
  status: string;
}

// --- Слайс API ---
export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({ baseUrl: BASE_URL }),
  // Теги для автоматической инвалидации кеша
  tagTypes: ['Plot', 'Lease'],
  endpoints: (builder) => ({
    // QUERIES
    getRegions: builder.query<Region[], void>({
      query: () => 'regions',
    }),
    getLandParcelsByRegion: builder.query<LandParcel[], string>({
      query: (regionId) => `regions/${regionId}/land-parcels`,
    }),
    getGreenhousesByLandParcel: builder.query<Greenhouse[], string>({
      query: (landParcelId) => `land-parcels/${landParcelId}/greenhouses`,
    }),
    getPlotsByGreenhouse: builder.query<Plot[], string>({
      query: (greenhouseId) => `greenhouses/${greenhouseId}/plots`,
      providesTags: (result) =>
        result
          ? [
              ...result.map(({ id }) => ({ type: 'Plot' as const, id })),
              { type: 'Plot', id: 'LIST' },
            ]
          : [{ type: 'Plot', id: 'LIST' }],
    }),
    getMyLeases: builder.query<PlotLease[], void>({
      query: () => 'me/leases',
      providesTags: ['Lease'],
    }),

    // MUTATIONS
    leasePlot: builder.mutation<PlotLease, string>({
      query: (plotId) => ({
        url: `plots/${plotId}/lease`,
        method: 'POST',
      }),
      // После успешной мутации, инвалидируем кеш для списков грядок и аренд
      invalidatesTags: [{ type: 'Plot', id: 'LIST' }, 'Lease'],
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
  useLeasePlotMutation,
} = apiSlice;
