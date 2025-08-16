import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

const BASE_URL = '/api/v1';

// --- Типы данных ---
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

interface Plot {
  id: string;
  name: string;
  status: string;
  greenhouse_id: string;
}

// --- Слайс API ---
export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({ baseUrl: BASE_URL }),
  endpoints: (builder) => ({
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
    }),
  }),
});

// --- Экспорт хуков ---
export const {
  useGetRegionsQuery,
  useGetLandParcelsByRegionQuery,
  useGetGreenhousesByLandParcelQuery,
  useGetPlotsByGreenhouseQuery,
} = apiSlice;