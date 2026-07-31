import React from 'react';
import { ProviderList } from '@/features/providers/components/ProviderList';

export const metadata = {
  title: 'Storage Providers - CloudStoreX',
  description: 'Manage multi-cloud storage connections',
};

export default function ProvidersPage() {
  return (
    <div className="container mx-auto py-8 max-w-7xl">
      <ProviderList />
    </div>
  );
}
