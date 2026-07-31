import { api } from '@/lib/api';
import { Provider, CreateProviderRequest, UpdateProviderRequest } from '../types';

interface ProviderResponse {
  data: Provider;
}

interface ProvidersResponse {
  data: Provider[];
}

export const getProviders = async (): Promise<Provider[]> => {
  const response = await api.get<any, ProvidersResponse>('/providers');
  return response.data;
};

export const getProvider = async (id: string): Promise<Provider> => {
  const response = await api.get<any, ProviderResponse>(`/providers/${id}`);
  return response.data;
};

export const createProvider = async (data: CreateProviderRequest): Promise<Provider> => {
  const response = await api.post<any, ProviderResponse>('/providers', data);
  return response.data;
};

export const updateProvider = async (id: string, data: UpdateProviderRequest): Promise<Provider> => {
  const response = await api.put<any, ProviderResponse>(`/providers/${id}`, data);
  return response.data;
};

export const deleteProvider = async (id: string): Promise<void> => {
  await api.delete(`/providers/${id}`);
};

export const validateProvider = async (id: string, extended: boolean = false): Promise<void> => {
  await api.post(`/providers/${id}/validate?extended=${extended}`);
};

export const setDefaultProvider = async (id: string): Promise<void> => {
  await api.post(`/providers/${id}/default`);
};

export const enableProvider = async (id: string): Promise<void> => {
  await api.post(`/providers/${id}/enable`);
};

export const disableProvider = async (id: string): Promise<void> => {
  await api.post(`/providers/${id}/disable`);
};
