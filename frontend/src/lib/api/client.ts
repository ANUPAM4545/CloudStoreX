import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { APIResponse, APIError } from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000,
});

apiClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => {
    // Backend APIResponse structure: { success: true, data: ... }
    return response.data;
  },
  (error: AxiosError<APIResponse>) => {
    if (error.response?.status === 401) {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token');
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
      }
    }

    const apiError: APIError = error.response?.data?.error || {
      code: 'unknown_error',
      message: error.message || 'An unexpected error occurred.',
    };

    return Promise.reject(apiError);
  }
);
