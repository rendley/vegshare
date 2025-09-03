import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

const BASE_URL = '/api/v1';

// --- Типы данных ---
// TODO: Вынести в /types
interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user_id: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'user';
  avatar_url?: string;
  created_at: string;
}

export interface Region {
  id: string;
  name: string;
}

export interface LandParcel {
  id: string;
  name: string;
  region_id: string;
}

export interface Structure {
  id: string;
  name: string;
  land_parcel_id: string;
  type?: string;
}

export interface Plot {
  id: string;
  name: string;
  structure_id: string;
  size?: string;
  status: string;
}

export interface Camera {
  id: string;
  name: string;
  plot_id: string;
  rtsp_path_name: string;
}

export interface Lease {
  id: string;
  user_id: string;
  unit_id: string;
  unit_type: string;
  status: string;
  start_date: string;
  end_date: string;
  created_at: string;
  updated_at: string;
}

export interface Crop {
  id: string;
  name: string;
  description?: string;
  planting_time?: number;
  harvest_time?: number;
}

export interface OperationLog {
    id: string;
    unit_id: string;
    unit_type: string;
    user_id: string;
    action_type: string;
    parameters: any; // Can be more specific, e.g., { crop_id: string } | { volume_liters: number }
    status: string;
    executed_at: string;
    created_at: string;
    updated_at: string;
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
  tagTypes: ['Region', 'LandParcel', 'Structure', 'Plot', 'Lease', 'OperationLog', 'Camera', 'Crop', 'User', 'StructureType'],
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
    getStructuresByLandParcel: builder.query<Structure[], string>({
      query: (landParcelId) => `farm/land-parcels/${landParcelId}/structures`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Structure' as const, id })), { type: 'Structure', id: 'LIST' }] : [{ type: 'Structure', id: 'LIST' }],
    }),
    getStructureTypes: builder.query<string[], void>({
      query: () => 'farm/structures/types',
      providesTags: ['StructureType'],
    }),
    getPlotsByStructure: builder.query<Plot[], string>({
      query: (structureId) => `plots?structure_id=${structureId}`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Plot' as const, id })), { type: 'Plot', id: 'LIST' }] : [{ type: 'Plot', id: 'LIST' }],
    }),
    getCamerasByPlot: builder.query<Camera[], string>({
      query: (plotId) => `plots/${plotId}/cameras`,
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Camera' as const, id })), { type: 'Camera', id: 'LIST' }] : [{ type: 'Camera', id: 'LIST' }],
    }),
    getMyLeases: builder.query<Lease[], void>({
      query: () => 'leasing',
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Lease' as const, id })), { type: 'Lease', id: 'LIST' }] : [{ type: 'Lease', id: 'LIST' }],
    }),
    getAvailableCrops: builder.query<Crop[], void>({
      query: () => 'catalog/crops',
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'Crop' as const, id })), { type: 'Crop', id: 'LIST' }] : [{ type: 'Crop', id: 'LIST' }],
    }),
    getActionsForUnit: builder.query<OperationLog[], string>({
        query: (unitId) => `operations/units/${unitId}/actions`,
        providesTags: (_result, _error, unitId) => [{ type: 'OperationLog', id: unitId }],
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
    leasePlot: builder.mutation<Lease, { unit_id: string; unit_type: string }>({
      query: (leaseRequest) => ({
        url: 'leasing',
        method: 'POST',
        body: leaseRequest,
      }),
      invalidatesTags: [{ type: 'Plot', id: 'LIST' }, { type: 'Lease', id: 'LIST' }],
    }),
    createAction: builder.mutation<OperationLog, { unit_id: string; unit_type: string; action_type: string; parameters: any }>(({
        query: (actionRequest) => ({
            url: 'operations/actions',
            method: 'POST',
            body: actionRequest,
        }),
        invalidatesTags: (_result, _error, arg) => [{ type: 'OperationLog', id: arg.unit_id }],
    })),
    cancelAction: builder.mutation<void, string>({
        query: (actionId) => ({
            url: `operations/actions/${actionId}`,
            method: 'DELETE',
        }),
        invalidatesTags: (_result, _error, _actionId) => [{ type: 'OperationLog', id: 'LIST' }], // This will refetch all actions for all units, which is not ideal. A more specific invalidation would be better.
    }),

    // Admin Queries
    getUsers: builder.query<User[], void>({
      query: () => 'admin/users',
      providesTags: (result) => result ? [...result.map(({ id }) => ({ type: 'User' as const, id })), { type: 'User', id: 'LIST' }] : [{ type: 'User', id: 'LIST' }],
    }),

    // Admin Mutations
    updateUserRole: builder.mutation<User, { userId: string; role: string }>({
      query: ({ userId, role }) => ({
        url: `admin/users/${userId}/role`,
        method: 'PUT',
        body: { role },
      }),
      invalidatesTags: (_, __, { userId }) => [{ type: 'User', id: userId }, { type: 'User', id: 'LIST' }],
    }),
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
    createStructure: builder.mutation<Structure, { landParcelId: string; name: string; type?: string }>({
      query: ({ landParcelId, ...body }) => ({
        url: `farm/land-parcels/${landParcelId}/structures`,
        method: 'POST',
        body,
      }),
      invalidatesTags: [{ type: 'Structure', id: 'LIST' }, 'StructureType'],
    }),
    createPlot: builder.mutation<Plot, { structure_id: string; name: string; size?: string }>({
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
    updateRegion: builder.mutation<Region, { id: string; name: string }>({
      query: ({ id, ...body }) => ({
        url: `farm/regions/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (_, __, { id }) => [{ type: 'Region', id }, { type: 'Region', id: 'LIST' }],
    }),
    deleteRegion: builder.mutation<void, string>({
      query: (id) => ({
        url: `farm/regions/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: [{ type: 'Region', id: 'LIST' }],
    }),
    updateLandParcel: builder.mutation<LandParcel, { id: string; name: string }>({
      query: ({ id, ...body }) => ({
        url: `farm/land-parcels/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (_, __, { id }) => [{ type: 'LandParcel', id }, { type: 'LandParcel', id: 'LIST' }],
    }),
    deleteLandParcel: builder.mutation<void, string>({
      query: (id) => ({
        url: `farm/land-parcels/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: [{ type: 'LandParcel', id: 'LIST' }],
    }),

    // Structures
    updateStructure: builder.mutation<Structure, { id: string; name: string; type?: string }>({
      query: ({ id, ...body }) => ({
        url: `farm/structures/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (_, __, { id }) => [{ type: 'Structure', id }, { type: 'Structure', id: 'LIST' }, 'StructureType'],
    }),
    deleteStructure: builder.mutation<void, string>({
      query: (id) => ({
        url: `farm/structures/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: [{ type: 'Structure', id: 'LIST' }, 'StructureType'],
    }),

    // Plots
    updatePlot: builder.mutation<Plot, { id: string; name: string; size?: string }>({
      query: ({ id, ...body }) => ({
        url: `plots/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (_, __, { id }) => [{ type: 'Plot', id }, { type: 'Plot', id: 'LIST' }],
    }),
    deletePlot: builder.mutation<void, string>({
      query: (id) => ({
        url: `plots/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: [{ type: 'Plot', id: 'LIST' }],
    }),

    // Cameras
    updateCamera: builder.mutation<Camera, { id: string; name: string; rtsp_path_name: string }>({
      query: ({ id, ...body }) => ({
        url: `cameras/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (_, __, { id }) => [{ type: 'Camera', id }, { type: 'Camera', id: 'LIST' }],
    }),
    deleteCamera: builder.mutation<void, string>({
      query: (id) => ({
        url: `cameras/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: [{ type: 'Camera', id: 'LIST' }],
    }),

    // Crops
    updateCrop: builder.mutation<Crop, Partial<Crop> & { id: string }>({
      query: ({ id, ...body }) => ({
        url: `catalog/crops/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: (_, __, { id }) => [{ type: 'Crop', id }, { type: 'Crop', id: 'LIST' }],
    }),
    deleteCrop: builder.mutation<void, string>({
      query: (id) => ({
        url: `catalog/crops/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: [{ type: 'Crop', id: 'LIST' }],
    }),
  }),
});

// --- Экспорт хуков ---
export const {
  useGetRegionsQuery,
  useGetLandParcelsByRegionQuery,
  useGetStructuresByLandParcelQuery,
  useGetStructureTypesQuery,
  useGetPlotsByStructureQuery,
  useGetCamerasByPlotQuery,
  useGetMyLeasesQuery,
  useGetAvailableCropsQuery,
  useGetActionsForUnitQuery,
  useLoginMutation,
  useRegisterMutation,
  useLeasePlotMutation,
  useCreateActionMutation,
  useCancelActionMutation,
  useGetUsersQuery,
  useUpdateUserRoleMutation,
  useCreateRegionMutation,
  useCreateLandParcelMutation,
  useCreateStructureMutation,
  useCreatePlotMutation,
  useCreateCameraMutation,
  useCreateCropMutation,
  useUpdateRegionMutation,
  useDeleteRegionMutation,
  useUpdateLandParcelMutation,
  useDeleteLandParcelMutation,
  useUpdateStructureMutation,
  useDeleteStructureMutation,
  useUpdatePlotMutation,
  useDeletePlotMutation,
  useUpdateCameraMutation,
  useDeleteCameraMutation,
  useUpdateCropMutation,
  useDeleteCropMutation,
} = apiSlice;