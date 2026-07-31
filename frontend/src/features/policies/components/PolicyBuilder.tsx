import React, { useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { RuleType } from '../types';

interface PolicyBuilderProps {
  onSubmit: (policy: any) => void;
  isSubmitting?: boolean;
}

export function PolicyBuilder({ onSubmit, isSubmitting }: PolicyBuilderProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState(10);
  const [ruleType, setRuleType] = useState<RuleType>('OBJECT_SIZE');
  const [providerId, setProviderId] = useState('');
  
  // Rule specific conditions
  const [minSize, setMinSize] = useState<number>(0);
  const [maxSize, setMaxSize] = useState<number>(10485760); // 10MB default
  const [region, setRegion] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    let conditions = {};
    if (ruleType === 'OBJECT_SIZE') {
      conditions = { min_size: minSize, max_size: maxSize };
    } else if (ruleType === 'REGION') {
      conditions = { region };
    }
    // other rules...

    const policy = {
      name,
      description,
      priority,
      rule_type: ruleType,
      conditions,
      actions: {
        provider_id: providerId
      }
    };
    
    onSubmit(policy);
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Create New Policy</CardTitle>
        <CardDescription>Define a storage routing policy for this workspace.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Policy Name</label>
              <Input 
                value={name} 
                onChange={(e) => setName(e.target.value)} 
                placeholder="e.g. Route Large Files" 
                required 
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Priority (Lower is higher)</label>
              <Input 
                type="number" 
                value={priority} 
                onChange={(e) => setPriority(parseInt(e.target.value))} 
                required 
              />
            </div>
          </div>
          
          <div className="space-y-2">
            <label className="text-sm font-medium">Description</label>
            <Input 
              value={description} 
              onChange={(e) => setDescription(e.target.value)} 
              placeholder="Optional description" 
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Rule Type</label>
            <select 
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background"
              value={ruleType}
              onChange={(e) => setRuleType(e.target.value as RuleType)}
            >
              <option value="OBJECT_SIZE">Object Size</option>
              <option value="REGION">Preferred Region</option>
              <option value="DEFAULT">Default Provider Override</option>
            </select>
          </div>

          {/* Conditional Inputs based on RuleType */}
          <div className="p-4 border rounded-md bg-muted/20 space-y-4">
            <h4 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">Conditions</h4>
            
            {ruleType === 'OBJECT_SIZE' && (
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium">Min Size (Bytes)</label>
                  <Input 
                    type="number" 
                    value={minSize} 
                    onChange={(e) => setMinSize(parseInt(e.target.value))} 
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium">Max Size (Bytes)</label>
                  <Input 
                    type="number" 
                    value={maxSize} 
                    onChange={(e) => setMaxSize(parseInt(e.target.value))} 
                  />
                </div>
              </div>
            )}

            {ruleType === 'REGION' && (
              <div className="space-y-2">
                <label className="text-sm font-medium">Region Name</label>
                <Input 
                  value={region} 
                  onChange={(e) => setRegion(e.target.value)} 
                  placeholder="e.g. us-east-1" 
                />
              </div>
            )}
            
            {ruleType === 'DEFAULT' && (
              <p className="text-sm text-muted-foreground">This rule matches all requests automatically.</p>
            )}
          </div>

          <div className="p-4 border rounded-md bg-emerald-500/5 border-emerald-500/20 space-y-4">
             <h4 className="text-sm font-semibold text-emerald-700 dark:text-emerald-400 uppercase tracking-wider">Action</h4>
             <div className="space-y-2">
                <label className="text-sm font-medium">Target Provider ID</label>
                <Input 
                  value={providerId} 
                  onChange={(e) => setProviderId(e.target.value)} 
                  placeholder="Enter the Provider UUID..." 
                  required 
                />
              </div>
          </div>

          <Button type="submit" disabled={isSubmitting} className="w-full">
            {isSubmitting ? 'Creating Policy...' : 'Create Policy'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
