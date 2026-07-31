import { api as apiClient } from "@/lib/api";
import { Policy, CreatePolicyRequest, UpdatePolicyRequest, RoutingDecision } from "../types";

export const policyApi = {
  getPolicies: () => 
    apiClient.get<Policy[]>('/policies'),
    
  getPolicy: (id: string) => 
    apiClient.get<Policy>(`/policies/${id}`),
    
  createPolicy: (data: CreatePolicyRequest) => 
    apiClient.post<Policy>('/policies', data),
    
  updatePolicy: (id: string, data: UpdatePolicyRequest) => 
    apiClient.put<Policy>(`/policies/${id}`, data),
    
  deletePolicy: (id: string) => 
    apiClient.delete(`/policies/${id}`),
    
  enablePolicy: (id: string) => 
    apiClient.post(`/policies/${id}/enable`, {}),
    
  disablePolicy: (id: string) => 
    apiClient.post(`/policies/${id}/disable`, {}),
    
  getRoutingDecisions: (limit: number = 100, offset: number = 0) => 
    apiClient.get<{items: RoutingDecision[], total: number}>(`/routing-decisions?limit=${limit}&offset=${offset}`),
};
