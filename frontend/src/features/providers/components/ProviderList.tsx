'use client';

import React, { useEffect, useState } from 'react';
import { 
  getProviders, 
  deleteProvider, 
  setDefaultProvider, 
  enableProvider, 
  disableProvider, 
  validateProvider 
} from '../api';
import { Provider } from '../types';
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { StatusBadge } from '@/components/ui/StatusBadge';
import { LoadingSkeleton } from '@/components/ui/LoadingSkeleton';
import { EmptyState } from '@/components/ui/EmptyState';
import { PageHeader } from '@/components/ui/PageHeader';

export const ProviderList: React.FC = () => {
  const [providers, setProviders] = useState<Provider[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadProviders = async () => {
    setLoading(true);
    try {
      const data = await getProviders();
      setProviders(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to load providers');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProviders();
  }, []);

  const handleSetDefault = async (id: string) => {
    try {
      await setDefaultProvider(id);
      await loadProviders();
    } catch (err: any) {
      alert(err.message || 'Failed to set default provider');
    }
  };

  const handleToggleEnable = async (provider: Provider) => {
    try {
      if (provider.is_enabled) {
        await disableProvider(provider.id);
      } else {
        await enableProvider(provider.id);
      }
      await loadProviders();
    } catch (err: any) {
      alert(err.message || 'Failed to toggle provider status');
    }
  };

  const handleValidate = async (id: string) => {
    try {
      await validateProvider(id, true);
      alert('Validation initiated');
      await loadProviders();
    } catch (err: any) {
      alert(err.message || 'Validation failed');
    }
  };

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this provider?')) {
      try {
        await deleteProvider(id);
        await loadProviders();
      } catch (err: any) {
        alert(err.message || 'Failed to delete provider');
      }
    }
  };

  if (loading) return <LoadingSkeleton count={3} type="card" />;

  if (error) return <div className="text-red-500 p-4">{error}</div>;

  return (
    <div className="space-y-6">
      <PageHeader 
        title="Storage Providers" 
        description="Manage your multi-cloud storage connections"
      >
        <Button>Add Provider</Button>
      </PageHeader>

      {providers.length === 0 ? (
        <EmptyState 
          title="No providers found" 
          description="You don't have any storage providers configured yet."
          action={<Button>Add Provider</Button>}
        />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {providers.map((p) => (
            <Card key={p.id} className={`border ${p.is_default ? 'border-primary' : 'border-border'}`}>
              <CardHeader className="pb-2">
                <div className="flex justify-between items-start">
                  <CardTitle className="text-xl font-bold flex items-center gap-2">
                    {p.provider_name}
                    {p.is_default && <StatusBadge status="info" label="Default" />}
                  </CardTitle>
                  <StatusBadge 
                    status={p.status === 'READY' ? 'success' : p.status === 'FAILED' ? 'error' : p.status === 'VALIDATING' ? 'warning' : 'neutral'} 
                    label={p.status} 
                  />
                </div>
                <div className="text-sm text-muted-foreground">{p.provider_type} • {p.region || 'Global'}</div>
              </CardHeader>
              <CardContent className="pb-4 space-y-4">
                
                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Health</span>
                  <span className={`font-medium ${p.health === 'HEALTHY' ? 'text-green-500' : p.health === 'UNAVAILABLE' ? 'text-red-500' : 'text-yellow-500'}`}>
                    {p.health} {p.latency_ms > 0 && `(${p.latency_ms}ms)`}
                  </span>
                </div>

                <div className="flex justify-between text-sm">
                  <span className="text-muted-foreground">Enabled</span>
                  <span className="font-medium">{p.is_enabled ? 'Yes' : 'No'}</span>
                </div>

                {p.last_error && (
                  <div className="text-xs text-red-500 bg-red-500/10 p-2 rounded truncate" title={p.last_error}>
                    {p.last_error}
                  </div>
                )}
                
                <div>
                  <div className="text-xs font-semibold mb-2 text-muted-foreground uppercase tracking-wider">Capabilities</div>
                  <div className="flex flex-wrap gap-1">
                    {Object.entries(p.capabilities).map(([key, value]) => {
                      if (value) {
                        return (
                          <span key={key} className="text-[10px] bg-secondary text-secondary-foreground px-2 py-1 rounded-full">
                            {key.replace('_', ' ')}
                          </span>
                        )
                      }
                      return null;
                    })}
                  </div>
                </div>

              </CardContent>
              <CardFooter className="flex flex-wrap gap-2 pt-0">
                <Button variant="outline" size="sm" onClick={() => handleValidate(p.id)}>
                  Validate
                </Button>
                <Button variant="outline" size="sm" onClick={() => handleToggleEnable(p)}>
                  {p.is_enabled ? 'Disable' : 'Enable'}
                </Button>
                {!p.is_default && p.is_enabled && (
                  <Button variant="outline" size="sm" onClick={() => handleSetDefault(p.id)}>
                    Make Default
                  </Button>
                )}
                <Button variant="destructive" size="sm" onClick={() => handleDelete(p.id)} disabled={p.is_default}>
                  Delete
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
};
