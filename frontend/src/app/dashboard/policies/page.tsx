"use client";

import React, { useEffect, useState } from 'react';
import { PageHeader } from '@/components/ui/PageHeader';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { StatusBadge } from '@/components/ui/StatusBadge';
import { Shield, Activity, ArrowRight, Server, FileText } from 'lucide-react';
import { Policy, RoutingDecision, CreatePolicyRequest } from '@/features/policies/types';
import { policyApi } from '@/features/policies/api';
import { PolicyBuilder } from '@/features/policies/components/PolicyBuilder';
import { formatDistanceToNow } from 'date-fns';

export default function PoliciesPage() {
  const [policies, setPolicies] = useState<Policy[]>([]);
  const [decisions, setDecisions] = useState<RoutingDecision[]>([]);
  const [loading, setLoading] = useState(true);
  const [showBuilder, setShowBuilder] = useState(false);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [policiesRes, decRes] = await Promise.all([
        policyApi.getPolicies(),
        policyApi.getRoutingDecisions(10) // fetch latest 10 decisions
      ]);
      setPolicies(policiesRes.data);
      setDecisions(decRes.data.items || []);
    } catch (error) {
      console.error("Failed to fetch policies data", error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreatePolicy = async (policy: any) => {
    try {
      await policyApi.createPolicy(policy as CreatePolicyRequest);
      setShowBuilder(false);
      fetchData();
    } catch (error) {
      console.error("Failed to create policy", error);
    }
  };

  const handleTogglePolicy = async (id: string, currentlyEnabled: boolean) => {
    try {
      if (currentlyEnabled) {
        await policyApi.disablePolicy(id);
      } else {
        await policyApi.enablePolicy(id);
      }
      fetchData();
    } catch (error) {
      console.error("Failed to toggle policy", error);
    }
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Policy Engine"
        description="Manage intelligent routing rules and inspect routing decisions."
      >
        <Button onClick={() => setShowBuilder(!showBuilder)}>
          {showBuilder ? 'Cancel' : 'Create Policy'}
        </Button>
      </PageHeader>

      {showBuilder && (
        <div className="mb-6">
          <PolicyBuilder onSubmit={handleCreatePolicy} />
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        {/* Policies List */}
        <div className="lg:col-span-2 space-y-4">
          <h2 className="text-xl font-bold flex items-center gap-2">
            <Shield className="h-5 w-5 text-indigo-500" />
            Active Policies
          </h2>
          {loading ? (
            <div className="p-8 text-center text-muted-foreground">Loading policies...</div>
          ) : policies.length === 0 ? (
            <Card className="border-dashed border-2 bg-transparent">
               <div className="p-12 text-center text-muted-foreground">
                 <Shield className="h-10 w-10 mx-auto mb-4 opacity-50" />
                 <p>No routing policies configured.</p>
                 <p className="text-sm">Storage will fallback to the workspace default provider.</p>
               </div>
            </Card>
          ) : (
            policies.map((p) => (
              <Card key={p.id} className={p.enabled ? '' : 'opacity-60'}>
                <CardContent className="p-6">
                  <div className="flex justify-between items-start mb-4">
                    <div>
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-xs font-mono bg-indigo-500/10 text-indigo-700 dark:text-indigo-400 px-2 py-0.5 rounded border border-indigo-500/20">
                          Priority: {p.priority}
                        </span>
                        <StatusBadge 
                          status={p.enabled ? 'success' : 'neutral'} 
                          label={p.enabled ? 'Enabled' : 'Disabled'} 
                        />
                      </div>
                      <h3 className="text-lg font-semibold">{p.name}</h3>
                      <p className="text-sm text-muted-foreground">{p.description}</p>
                    </div>
                    <Button variant="outline" size="sm" onClick={() => handleTogglePolicy(p.id, p.enabled)}>
                      {p.enabled ? 'Disable' : 'Enable'}
                    </Button>
                  </div>
                  
                  <div className="grid grid-cols-2 gap-4 mt-4 p-4 bg-muted/30 rounded-md">
                     <div>
                       <span className="text-xs text-muted-foreground uppercase tracking-wider block mb-1">Condition</span>
                       <div className="font-medium text-sm flex items-center gap-2">
                         <FileText className="h-4 w-4 text-muted-foreground" />
                         {p.rule_type}: {JSON.stringify(p.conditions)}
                       </div>
                     </div>
                     <div>
                       <span className="text-xs text-muted-foreground uppercase tracking-wider block mb-1">Action</span>
                       <div className="font-medium text-sm flex items-center gap-2">
                         <ArrowRight className="h-4 w-4 text-emerald-500" />
                         Route to: <span className="font-mono text-xs">{p.actions.provider_id.substring(0, 8)}...</span>
                       </div>
                     </div>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>

        {/* Audit Trail / Routing Decisions */}
        <div className="space-y-4">
          <h2 className="text-xl font-bold flex items-center gap-2">
            <Activity className="h-5 w-5 text-emerald-500" />
            Routing Audit Trail
          </h2>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium">Recent Decisions</CardTitle>
              <CardDescription>Live log of provider selections</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {decisions.length === 0 ? (
                  <p className="text-sm text-muted-foreground text-center py-4">No routing decisions yet.</p>
                ) : (
                  decisions.map((d) => (
                    <div key={d.id} className="text-sm border-l-2 border-indigo-500/50 pl-4 py-1 relative">
                      <div className="absolute w-2 h-2 rounded-full bg-indigo-500 -left-[5px] top-2"></div>
                      <div className="flex justify-between items-start mb-1">
                        <span className="font-semibold truncate max-w-[150px]">{d.object_key || 'Bucket Op'}</span>
                        <span className="text-xs text-muted-foreground">
                          {formatDistanceToNow(new Date(d.timestamp))} ago
                        </span>
                      </div>
                      <div className="text-xs text-muted-foreground mb-1 line-clamp-2">
                        {d.decision_reason}
                      </div>
                      <div className="flex items-center gap-2 mt-2">
                        <StatusBadge 
                          status={d.is_fallback ? 'warning' : 'info'} 
                          label={d.is_fallback ? 'Fallback' : d.matched_rule} 
                          dot={false}
                        />
                        <span className="flex items-center gap-1 text-xs font-mono text-emerald-600 dark:text-emerald-400 bg-emerald-500/10 px-1.5 py-0.5 rounded">
                          <Server className="h-3 w-3" />
                          {d.provider_id.substring(0, 8)}
                        </span>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </div>

      </div>
    </div>
  );
}
